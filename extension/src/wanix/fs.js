//@ts-ignore
import * as duplex from "../duplex/duplex.min.js";

export class WanixFS {
    constructor(port) {
        const sess = new duplex.Session(new duplex.PortConn(port));
        this.peer = new duplex.Peer(sess, new duplex.CBORCodec());
    }

    async readDir(name) {
        return (await this.peer.call("ReadDir", [name])).value;
    }

    async makeDir(name) {
        await this.peer.call("Mkdir", [name]);
    }

    async makeDirAll(name) {
        await this.peer.call("MkdirAll", [name]);
    }

    async bind(name, newname) {
        await this.peer.call("Bind", [name, newname]);
    }

    async unbind(name, newname) {
        await this.peer.call("Unbind", [name, newname]);
    }
    
    async readFile(name) {
        return (await this.peer.call("ReadFile", [name])).value;
    }

    async readText(name) {
        return (new TextDecoder()).decode(await this.readFile(name));
    }

    async waitFor(name, timeoutMs=1000) {
        await this.peer.call("WaitFor", [name, timeoutMs]);
    }

    async stat(name) {
        return (await this.peer.call("Stat", [name])).value;
    }

    async writeFile(name, contents) {
        if (typeof contents === "string") {
            contents = (new TextEncoder()).encode(contents);
        }
        return (await this.peer.call("WriteFile", [name, contents])).value;
    }

    async appendFile(name, contents) {
        if (typeof contents === "string") {
            contents = (new TextEncoder()).encode(contents);
        }
        return (await this.peer.call("AppendFile", [name, contents])).value;
    }

    async rename(oldname, newname) {
        await this.peer.call("Rename", [oldname, newname]);
    }

    async copy(oldname, newname) {
        await this.peer.call("Copy", [oldname, newname]);
    }

    async remove(name) {
        await this.peer.call("Remove", [name]);
    }

    async removeAll(name) {
        await this.peer.call("RemoveAll", [name]);
    }

    async truncate(name, size) {
        await this.peer.call("Truncate", [name, size]);
    }

    async open(name) {
        return (await this.peer.call("Open", [name])).value;
    }

    async read(fd, count) {
        return (await this.peer.call("Read", [fd, count])).value;
    }

    async write(fd, data) {
        return (await this.peer.call("Write", [fd, data])).value;
    }

    async close(fd) {
        return (await this.peer.call("Close", [fd])).value;
    }

    async sync(fd) {
        return (await this.peer.call("Sync", [fd])).value;
    }

    async openReadable(name) {
        const fd = await this.open(name);
        return this.readable(fd);
    }

    async openWritable(name) {
        const fd = await this.open(name);
        return this.writable(fd);
    }

    writable(fd) {
        const self = this;
        return new WritableStream({
            write(chunk) {
                return self.write(fd, chunk);
            },
        });
    }

    readable(fd) {
        const self = this;
        return new ReadableStream({
            async pull(controller) {
                const data = await self.read(fd, 1024);
                if (data === null) {
                    controller.close();
                }
                controller.enqueue(data);
            },
        });
    }
}