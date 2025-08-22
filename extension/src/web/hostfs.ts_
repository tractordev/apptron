import {
	CancellationToken,
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

interface RemoteEntry {
    IsDir: boolean;
    Name: string;
    Ctime: number;
    Mtime: number;
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
		this.mtime = entry.Mtime;
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
		this.mtime = entry.Mtime;
		this.size = entry.Size;
		this.name = entry.Name;
	}
}

export type Entry = File | Directory;

export class HostFS implements FileSystemProvider, /*FileSearchProvider, TextSearchProvider,*/ Disposable {
	static scheme = 'hostfs';

	private readonly disposable: Disposable;
    private readonly peer: any;

	constructor(peer: any) {
        this.peer = peer;
		this.disposable = Disposable.from(
			workspace.registerFileSystemProvider(HostFS.scheme, this, { isCaseSensitive: true }),
			// workspace.registerFileSearchProvider(MemFS.scheme, this),
			// workspace.registerTextSearchProvider(MemFS.scheme, this)
		);
	}

	dispose() {
		this.disposable?.dispose();
	}

	// root = new Directory(Uri.parse('memfs:/'), '');

	// --- manage file metadata

	stat(uri: Uri): Thenable<FileStat> {
        // console.log("stat", uri);
        return new Promise<FileStat>((resolve: (value: FileStat) => void, reject: (reason?: any) => void) => {
            this._lookup(uri, false).then(resolve, reject);
        });
	}

	readDirectory(uri: Uri): Thenable<[string, FileType][]> {
        // console.log("readDir", uri);
        return new Promise<[string, FileType][]>(async (resolve: (value: [string, FileType][]) => void, reject: (reason?: any) => void) => {
            try {
                const resp = await this.peer.call("vscode.ReadDir", [uri.path]);
                let result: [string, FileType][] = [];
                for (const entry of resp.value as RemoteEntry[]) {
                    result.push([entry.Name, (entry.IsDir) ? FileType.Directory : FileType.File]);
                }
                resolve(result);
            } catch (e) {
                reject(e);
            }
        });
	}

	// --- manage file contents

	readFile(uri: Uri): Thenable<Uint8Array> {
        // console.log("readFile", uri);
        return new Promise<Uint8Array>(async (resolve: (value: Uint8Array) => void, reject: (reason?: any) => void) => {
            try {
                const resp = await this.peer.call("vscode.ReadFile", [uri.path]);
                resolve(resp.value);
            } catch(e) {
                reject(FileSystemError.FileNotFound());
            }
        });
	}

	writeFile(uri: Uri, content: Uint8Array, options: { create: boolean, overwrite: boolean }): Thenable<void> {
        return new Promise<void>(async (resolve: (value: void) => void, reject: (reason?: any) => void) => {
            let entry = await this._lookup(uri, true);
            if (entry instanceof Directory) {
                reject(FileSystemError.FileIsADirectory(uri));
            }
            if (!entry && !options.create) {
                reject(FileSystemError.FileNotFound(uri));
            }
            if (entry && options.create && !options.overwrite) {
                reject(FileSystemError.FileExists(uri));
            }

            try {
                await this.peer.call("vscode.WriteFile", [uri.path, content]);
            } catch (e) {
                reject(e);
            }
            

            if (!entry) {
                //entry = new File(uri, basename);
                // parent.entries.set(basename, entry);
                this._fireSoon({ type: FileChangeType.Created, uri });
            }
            // entry.mtime = Date.now();
            // entry.size = content.byteLength;
            // entry.data = content;

            this._fireSoon({ type: FileChangeType.Changed, uri });

            resolve();
        });
	}

	// --- manage files/folders

    copy(source: Uri, destination: Uri, options: {overwrite: boolean}): Thenable<void> {
        return new Promise<void>(async (resolve: (value: void) => void, reject: (reason?: any) => void) => {
            reject("not implemented");
        });
    }

	rename(oldUri: Uri, newUri: Uri, options: { overwrite: boolean }): Thenable<void> {
        return new Promise<void>(async (resolve: (value: void) => void, reject: (reason?: any) => void) => {
            if (!options.overwrite && await this._lookup(newUri, true)) {
                reject(FileSystemError.FileExists(newUri));
            }
    
            let entry = await this._lookup(oldUri, false);
            let oldParent = await this._lookupParentDirectory(oldUri);
    
            let newParent = await this._lookupParentDirectory(newUri);
            let newName = this._basename(newUri.path);
    
            // TODO
    
            // oldParent.entries.delete(entry.name);
            // entry.name = newName;
            // newParent.entries.set(newName, entry);
    
            this._fireSoon(
                { type: FileChangeType.Deleted, uri: oldUri },
                { type: FileChangeType.Created, uri: newUri }
            );

            resolve();
        });
	}

	delete(uri: Uri, options: {recursive: boolean}): Thenable<void> {
        return new Promise<void>(async (resolve: (value: void) => void, reject: (reason?: any) => void) => {
            let dirname = uri.with({ path: this._dirname(uri.path) });
            // let basename = this._basename(uri.path);
            // let parent = this._lookupAsDirectory(dirname, false);
            // if (!parent.entries.has(basename)) {
            // 	throw FileSystemError.FileNotFound(uri);
            // }

            // TODO

            // parent.entries.delete(basename);
            // parent.mtime = Date.now();
            // parent.size -= 1;
            this._fireSoon({ type: FileChangeType.Changed, uri: dirname }, { uri, type: FileChangeType.Deleted });

            resolve();
        });
	}

	createDirectory(uri: Uri): Promise<void> {
        return new Promise<void>(async (resolve: (value: void) => void, reject: (reason?: any) => void) => {
            // let basename = this._basename(uri.path);
            let dirname = uri.with({ path: this._dirname(uri.path) });
            // let parent = await this._lookupAsDirectory(dirname, false);

            try {
                await this.peer.call("vscode.MakeDir", [uri.path]);
            } catch(e) {
                reject(e);
            }

            // let entry = new Directory(uri, basename);
            // parent.entries.set(entry.name, entry);
            // parent.mtime = Date.now();
            // parent.size += 1;
            this._fireSoon({ type: FileChangeType.Changed, uri: dirname }, { type: FileChangeType.Created, uri });

            resolve();
        });
	}

	// --- lookup

	private async _lookup(uri: Uri, silent: false): Promise<Entry>;
	private async _lookup(uri: Uri, silent: boolean): Promise<Entry | undefined>;
	private async _lookup(uri: Uri, silent: boolean): Promise<Entry | undefined> {
        try {
            const resp = await this.peer.call("vscode.Stat", [uri.path]);
            const entry: RemoteEntry = resp.value as RemoteEntry;
            if (entry.IsDir) {
                return new Directory(uri, entry);
            } else {
                return new File(uri, entry);
            }
        } catch (e) {
            if (!silent) {
                //console.error(e);
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
