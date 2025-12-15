
import * as vscode from 'vscode';
import { WanixBridge } from './bridge.js';


declare const navigator: unknown;

export async function activate(context: vscode.ExtensionContext) {
	if (typeof navigator !== 'object') {	// do not run under node.js
		console.error("not running in browser");
		return;
	}
	
	const channel = new MessageChannel();
	const bridge = new WanixBridge(channel.port2, "vm/1/fsys");
	context.subscriptions.push(bridge);

	const port = (context as any).messagePassingProtocol;
	port.postMessage({type: "_port", port: channel.port1}, [channel.port1]);

	bridge.ready.then((wfsys) => {
		console.log("bridge ready");
		const terminal = createTerminal(wfsys);
		context.subscriptions.push(terminal);
		terminal.show();

		(async () => {
			const dec = new TextDecoder();
			const stream = await wfsys.openReadable("#commands/data1");
			for await (const chunk of stream) {
				const args = dec.decode(chunk).trim().split(" ");
				const cmd = args.shift();
				vscode.commands.executeCommand(`apptron.${cmd}`, ...args);
			}
		})();		
	});


	context.subscriptions.push(vscode.commands.registerCommand('apptron.open-preview', (filepath?: string) => {
		if (!filepath) {
			return;
		}
		vscode.commands.executeCommand('markdown.showPreview', vscode.Uri.parse(`wanix://${filepath}`));
	}));

	context.subscriptions.push(vscode.commands.registerCommand('apptron.open-file', (filepath?: string) => {
		if (!filepath) {
			return;
		}
		vscode.commands.executeCommand('vscode.open', vscode.Uri.parse(`wanix://${filepath}`));
	}));

	context.subscriptions.push(vscode.commands.registerCommand('apptron.open-folder', (filepath?: string) => {
		if (!filepath) {
			return;
		}
		const folders = vscode.workspace.workspaceFolders;
		const insertIndex = folders ? folders.length : 0;
		const uri = vscode.Uri.parse(`wanix://${filepath}`);
		vscode.workspace.updateWorkspaceFolders(
			insertIndex, // insert at the end
			0, // number of folders to remove
			{ uri }
		);
	}));

	
	console.log('Apptron system extension activated');
}

function createTerminal(wx: any) {
	const writeEmitter = new vscode.EventEmitter<string>();
	let channel: any = undefined;
	const dec = new TextDecoder();
	const enc = new TextEncoder();
	const pty = {
		onDidWrite: writeEmitter.event,
		open: () => {
			(async () => {
				const stream = await wx.openReadable("#console/data");
				for await (const chunk of stream) {
					writeEmitter.fire(dec.decode(chunk));
				}
			})();
		},
		close: () => {
			// if (channel) {
			// 	channel.close();
			// }
		},
		handleInput: (data: string) => {
			wx.appendFile("#console/data", data);
		}
	};
	return vscode.window.createTerminal({ name: `Shell`, pty });
}


// @ts-ignore
// polyfill for ReadableStream.prototype[Symbol.asyncIterator] on safari
if (!ReadableStream.prototype[Symbol.asyncIterator]) {
	// @ts-ignore
    ReadableStream.prototype[Symbol.asyncIterator] = async function* () {
        const reader = this.getReader();
        try {
            while (true) {
                const { done, value } = await reader.read();
                if (done) return;
                yield value;
            }
        } finally {
            reader.releaseLock();
        }
    };
}