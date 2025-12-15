import {
	Disposable,
	Event,
	EventEmitter,
	FileChangeEvent,
	FileChangeType,
	FileStat,
	FileSystemError,
	FileSystemProvider,
	FileType,
	Uri,
	workspace,
} from 'vscode';

//@ts-ignore
import { WanixFS } from "../wanix/fs.js";

interface RemoteEntry {
    IsDir: boolean;
    Name: string;
    // Ctime: number;
    ModTime: number;
    Size: number;
}

export class File implements FileStat {

	type: FileType;
	ctime: number;
	mtime: number;
	size: number;

	name: string;

	constructor(public uri: Uri, entry: RemoteEntry) {
		this.type = FileType.File;
		this.ctime = 0;
		this.mtime = entry.ModTime || 0;
		this.size = entry.Size;
		this.name = entry.Name;
	}
}

export class Directory implements FileStat {

	type: FileType;
	ctime: number;
	mtime: number;
	size: number;

	name: string;

	constructor(public uri: Uri, entry: RemoteEntry) {
		this.type = FileType.Directory;
		this.ctime = 0;
		this.mtime = entry.ModTime || 0;
		this.size = entry.Size;
		this.name = entry.Name;
	}
}

export type Entry = File | Directory;

export class WanixBridge implements FileSystemProvider, /*FileSearchProvider, TextSearchProvider,*/ Disposable {
	static scheme = 'wanix';

	public wfsys: any;
	public readonly ready: Promise<any>;
	private readonly disposable: Disposable;
    private root: string;	

	constructor(wanixReceiver: MessagePort, root: string) {
		this.ready = new Promise((resolve) => {
			wanixReceiver.onmessage = async (event) => {
				if (event.data.wanix) {
					console.log("wanix port received");
					const wfsys = new WanixFS(event.data.wanix);
					wfsys.waitFor("vm/1/fsys").then(() => {
						this.wfsys = wfsys;
						resolve(wfsys);
					});
				}
			}
		});
		this.root = root;
		this.disposable = Disposable.from(
			workspace.registerFileSystemProvider(WanixBridge.scheme, this, { isCaseSensitive: true }),
			// workspace.registerFileSearchProvider(MemFS.scheme, this),
			// workspace.registerTextSearchProvider(MemFS.scheme, this)
		);
	}

	join(path: string): string {
		if (path === "/") {
			return this.root;
		}
		return this.root + path;
	}

	dispose() {
		this.disposable?.dispose();
	}

	// --- manage file metadata

	stat(uri: Uri): Thenable<FileStat> {
		return this._stat(uri);
	}

	async _stat(uri: Uri): Promise<FileStat> {
		if (!this.wfsys) {
			if (uri.path !== "/project") {
				if (uri.path.includes(".vscode")) {
					throw FileSystemError.FileNotFound(uri);
				}
				return new File(uri, {
					IsDir: false,
					Name: this._basename(uri.path),
					ModTime: 0,
					Size: 0,
				});
			}
			// todo: watch root to force reload?
			return new Directory(uri, {
				IsDir: true,
				Name: uri.path,
				ModTime: 0,
				Size: 0,
			});
		}
		await this.ready;
		return await this._lookup(uri, false);
	}

	readDirectory(uri: Uri): Thenable<[string, FileType][]> {
		return this._readDirectory(uri);
	}

	async _readDirectory(uri: Uri): Promise<[string, FileType][]> {
        await this.ready;
		const entries = await this.wfsys.readDir(this.join(uri.path));
		let result: [string, FileType][] = [];
		for (const entry of entries) {
			result.push([entry.replace(/\/$/, ''), (entry.endsWith('/')) ? FileType.Directory : FileType.File]);
		}
		return result;
	}

	// --- manage file contents

	readFile(uri: Uri): Thenable<Uint8Array> {
		return this._readFile(uri);
	}

	async _readFile(uri: Uri): Promise<Uint8Array> {
		await this.ready;
		return await this.wfsys.readFile(this.join(uri.path));
	}

	writeFile(uri: Uri, content: Uint8Array, options: { create: boolean, overwrite: boolean }): Thenable<void> {
		return this._writeFile(uri, content, options);
	}

	async _writeFile(uri: Uri, content: Uint8Array, options: { create: boolean, overwrite: boolean }): Promise<void> {
		await this.ready;
		let entry = await this._lookup(uri, true);
		if (entry instanceof Directory) {
			throw FileSystemError.FileIsADirectory(uri);
		}
		if (!entry && !options.create) {
			throw FileSystemError.FileNotFound(uri);
		}
		if (entry && options.create && !options.overwrite) {
			throw FileSystemError.FileExists(uri);
		}

		await this.wfsys.writeFile(this.join(uri.path), content);
		
		if (!entry) {
			this._fireSoon({ type: FileChangeType.Created, uri });
		} else {
			this._fireSoon({ type: FileChangeType.Changed, uri });
		}
		this._fireSoon(
			{ type: FileChangeType.Changed, uri: uri.with({ path: this._dirname(uri.path) }) }
		);
	}

	// --- manage files/folders

    copy(source: Uri, destination: Uri, options: {overwrite: boolean}): Thenable<void> {
		return this._copy(source, destination, options);
	}

	async _copy(source: Uri, destination: Uri, options: {overwrite: boolean}): Promise<void> {
		await this.ready;
		if (!options.overwrite && await this._lookup(destination, true)) {
			throw FileSystemError.FileExists(destination);
		}

		await this.wfsys.copy(this.join(source.path), this.join(destination.path));

		this._fireSoon(
			{ type: FileChangeType.Changed, uri: destination.with({ path: this._dirname(destination.path) }) },
			{ type: FileChangeType.Created, uri: destination }
		);
	}

	rename(oldUri: Uri, newUri: Uri, options: { overwrite: boolean }): Thenable<void> {
		return this._rename(oldUri, newUri, options);
	}

	async _rename(oldUri: Uri, newUri: Uri, options: { overwrite: boolean }): Promise<void> {
		await this.ready;
		if (!options.overwrite && await this._lookup(newUri, true)) {
			throw FileSystemError.FileExists(newUri);
		}

		await this.wfsys.rename(this.join(oldUri.path), this.join(newUri.path));

		this._fireSoon(
			{ type: FileChangeType.Changed, uri: oldUri.with({ path: this._dirname(oldUri.path) }) },
			{ type: FileChangeType.Deleted, uri: oldUri },
			{ type: FileChangeType.Changed, uri: newUri.with({ path: this._dirname(newUri.path) }) },
			{ type: FileChangeType.Created, uri: newUri }
		);
	}

	delete(uri: Uri, options: {recursive: boolean}): Thenable<void> {
		return this._delete(uri, options);
	}

	async _delete(uri: Uri, options: {recursive: boolean}): Promise<void> {
		await this.ready;
		if (options.recursive) {
			await this.wfsys.removeAll(this.join(uri.path));
		} else {
			await this.wfsys.remove(this.join(uri.path));
		}

		this._fireSoon(
			{ type: FileChangeType.Changed, uri: uri.with({ path: this._dirname(uri.path) }) }, 
			{ uri, type: FileChangeType.Deleted }
		);
	}

	createDirectory(uri: Uri): Promise<void> {
		return this._createDirectory(uri);
	}

	async _createDirectory(uri: Uri): Promise<void> {
		await this.ready;
		await this.wfsys.makeDir(this.join(uri.path));
		this._fireSoon(
			{ type: FileChangeType.Changed, uri: uri.with({ path: this._dirname(uri.path) }) }, 
			{ type: FileChangeType.Created, uri }
		);
	}

	// --- lookup

	private async _lookup(uri: Uri, silent: false): Promise<Entry>;
	private async _lookup(uri: Uri, silent: boolean): Promise<Entry | undefined>;
	private async _lookup(uri: Uri, silent: boolean): Promise<Entry | undefined> {
        try {
            const entry = await this.wfsys.stat(this.join(uri.path));
            if (entry.IsDir) {
                return new Directory(uri, entry);
            } else {
                return new File(uri, entry);
            }
        } catch (e) {
            if (!silent) {
                // console.error(e);
                throw FileSystemError.FileNotFound(uri);
            } else {
                return undefined;
            }
        }
	}

	private async _lookupAsDirectory(uri: Uri, silent: boolean): Promise<Directory> {
		let entry = await this._lookup(uri, silent);
		if (entry instanceof Directory) {
			return entry;
		}
		throw FileSystemError.FileNotADirectory(uri);
	}

	private async _lookupAsFile(uri: Uri, silent: boolean): Promise<File> {
		let entry = await this._lookup(uri, silent);
		if (entry instanceof File) {
			return entry;
		}
		throw FileSystemError.FileIsADirectory(uri);
	}

	private async _lookupParentDirectory(uri: Uri): Promise<Directory> {
		const dirname = uri.with({ path: this._dirname(uri.path) });
		return await this._lookupAsDirectory(dirname, false);
	}

	// --- manage file events

	private _emitter = new EventEmitter<FileChangeEvent[]>();
	private _bufferedEvents: FileChangeEvent[] = [];
	private _fireSoonHandle?: any;

	readonly onDidChangeFile: Event<FileChangeEvent[]> = this._emitter.event;

	watch(_resource: Uri): Disposable {
		// ignore, fires for all changes...
		return new Disposable(() => { });
	}

	private _fireSoon(...events: FileChangeEvent[]): void {
		this._bufferedEvents.push(...events);

		if (this._fireSoonHandle) {
			clearTimeout(this._fireSoonHandle);
		}

		this._fireSoonHandle = setTimeout(() => {
			this._emitter.fire(this._bufferedEvents);
			this._bufferedEvents.length = 0;
		}, 5);
	}

	// --- path utils

	private _basename(path: string): string {
		path = this._rtrim(path, '/');
		if (!path) {
			return '';
		}

		return path.substr(path.lastIndexOf('/') + 1);
	}

	private _dirname(path: string): string {
		path = this._rtrim(path, '/');
		if (!path) {
			return '/';
		}

		return path.substr(0, path.lastIndexOf('/'));
	}

	private _rtrim(haystack: string, needle: string): string {
		if (!haystack || !needle) {
			return haystack;
		}

		const needleLen = needle.length,
			haystackLen = haystack.length;

		if (needleLen === 0 || haystackLen === 0) {
			return haystack;
		}

		let offset = haystackLen,
			idx = -1;

		while (true) {
			idx = haystack.lastIndexOf(needle, offset - 1);
			if (idx === -1 || idx + needleLen !== offset) {
				break;
			}
			if (idx === 0) {
				return '';
			}
			offset = idx;
		}

		return haystack.substring(0, offset);
	}

}
