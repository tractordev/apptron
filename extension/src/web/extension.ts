
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
		const terminal = createTerminal(wfsys);
		context.subscriptions.push(terminal);
		terminal.show();
	});

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


