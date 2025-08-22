
import * as vscode from 'vscode';
// import { HostFS } from './hostfs.js';

//@ts-ignore
import * as duplex from "../duplex/duplex.min.js";

declare const navigator: unknown;

export async function activate(context: vscode.ExtensionContext) {
	console.log('Apptron System extension activated');
	// if (typeof navigator !== 'object') {	// do not run under node.js
	// 	console.error("not running in browser");
	// 	return;
	// }

	// const channel = new MessageChannel();
	// self.postMessage({type: "_port", port: channel.port2}, [channel.port2]);

	// const sess = new duplex.Session(new duplex.PortConn(channel.port1));
	// const peer = new duplex.Peer(sess, new duplex.CBORCodec());
	// peer.respond();

	// const fs = new HostFS(peer);
	// context.subscriptions.push(fs);

	// const terminal = createTerminal(peer);
	// terminal.show();

}

function createTerminal(peer: any) {
	const writeEmitter = new vscode.EventEmitter<string>();
	let channel: any = undefined;
	const dec = new TextDecoder();
	const enc = new TextEncoder();
	const pty = {
		onDidWrite: writeEmitter.event,
		open: () => {
			(async () => {
				const resp = await peer.call("vscode.Terminal");
				channel = resp.channel;
				const b = new Uint8Array(1024);
				let gotEOF = false;
				while (gotEOF === false) {
				  const n = await channel.read(b);
				  if (n === null) {
					gotEOF = true;
				  } else {
					writeEmitter.fire(dec.decode(b.subarray(0, n)));
				  }
				}
			})();
		},
		close: () => {
			if (channel) {
				channel.close();
			}
		},
		handleInput: (data: string) => {
			if (channel) {
				channel.write(enc.encode(data));
			}
		}
	};
	return vscode.window.createTerminal({ name: `Shell`, pty });
}


