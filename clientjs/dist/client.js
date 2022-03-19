// deno-fmt-ignore-file
// deno-lint-ignore-file
// This code was bundled using `deno bundle` and it's not recommended to edit it manually

var payloadSizes, debug4, frames, options1;
function copy(a, c, b = 0) {
    b = Math.max(0, Math.min(b, c.byteLength));
    const d = c.byteLength - b;
    return a.byteLength > d && (a = a.subarray(0, d)), c.set(a, b), a.byteLength;
}
const MIN_READ = 32 * 1024, MAX_SIZE = 2 ** 32 - 2;
class Buffer1 {
    constructor(a){
        this._buf = a === void 0 ? new Uint8Array(0) : new Uint8Array(a), this._off = 0;
    }
    bytes(a = {
        copy: !0
    }) {
        return a.copy === !1 ? this._buf.subarray(this._off) : this._buf.slice(this._off);
    }
    empty() {
        return this._buf.byteLength <= this._off;
    }
    get length() {
        return this._buf.byteLength - this._off;
    }
    get capacity() {
        return this._buf.buffer.byteLength;
    }
    truncate(a) {
        if (a === 0) {
            this.reset();
            return;
        }
        if (a < 0 || a > this.length) throw Error("bytes.Buffer: truncation out of range");
        this._reslice(this._off + a);
    }
    reset() {
        this._reslice(0), this._off = 0;
    }
    _tryGrowByReslice(b) {
        const a = this._buf.byteLength;
        return b <= this.capacity - a ? (this._reslice(a + b), a) : -1;
    }
    _reslice(a) {
        this._buf = new Uint8Array(this._buf.buffer, 0, a);
    }
    readSync(a) {
        if (this.empty()) return this.reset(), a.byteLength === 0 ? 0 : null;
        const b = copy(this._buf.subarray(this._off), a);
        return this._off += b, b;
    }
    read(a) {
        const b = this.readSync(a);
        return Promise.resolve(b);
    }
    writeSync(a) {
        const b = this._grow(a.byteLength);
        return copy(a, this._buf, b);
    }
    write(a) {
        const b = this.writeSync(a);
        return Promise.resolve(b);
    }
    _grow(a) {
        const b = this.length;
        b === 0 && this._off !== 0 && this.reset();
        const d = this._tryGrowByReslice(a);
        if (d >= 0) return d;
        const c = this.capacity;
        if (a <= Math.floor(c / 2) - b) copy(this._buf.subarray(this._off), this._buf);
        else if (c + a > MAX_SIZE) throw new Error("The buffer cannot be grown beyond the maximum size.");
        else {
            const b = new Uint8Array(Math.min(2 * c + a, MAX_SIZE));
            copy(this._buf.subarray(this._off), b), this._buf = b;
        }
        return this._off = 0, this._reslice(Math.min(b + a, MAX_SIZE)), b;
    }
    grow(a) {
        if (a < 0) throw Error("Buffer.grow: negative count");
        const b = this._grow(a);
        this._reslice(b);
    }
    async readFrom(b) {
        let a = 0;
        const c = new Uint8Array(MIN_READ);
        while(!0){
            const e = this.capacity - this.length < MIN_READ, f = e ? c : new Uint8Array(this._buf.buffer, this.length), d = await b.read(f);
            if (d === null) return a;
            e ? this.writeSync(f.subarray(0, d)) : this._reslice(this.length + d), a += d;
        }
    }
    readFromSync(b) {
        let a = 0;
        const c = new Uint8Array(MIN_READ);
        while(!0){
            const e = this.capacity - this.length < MIN_READ, f = e ? c : new Uint8Array(this._buf.buffer, this.length), d = b.readSync(f);
            if (d === null) return a;
            e ? this.writeSync(f.subarray(0, d)) : this._reslice(this.length + d), a += d;
        }
    }
}
class JSONCodec1 {
    constructor(a = !1){
        this.debug = a;
    }
    encoder(a) {
        return new JSONEncoder1(a, this.debug);
    }
    decoder(a) {
        return new JSONDecoder1(a, this.debug);
    }
}
class JSONEncoder1 {
    constructor(a, b = !1){
        this.w = a, this.enc = new TextEncoder, this.debug = b;
    }
    async encode(b) {
        this.debug && console.log("<<", b);
        let c = this.enc.encode(JSON.stringify(b)), a = 0;
        while(a < c.length)a += await this.w.write(c.subarray(a));
    }
}
class JSONDecoder1 {
    constructor(a, b = !1){
        this.r = a, this.dec = new TextDecoder, this.debug = b;
    }
    async decode(c) {
        const a = new Uint8Array(c), d = await this.r.read(a);
        if (d === null) return Promise.resolve(null);
        let b = JSON.parse(this.dec.decode(a));
        return this.debug && console.log(">>", b), Promise.resolve(b);
    }
}
class FrameCodec1 {
    constructor(a){
        this.codec = a;
    }
    encoder(a) {
        return new FrameEncoder1(a, this.codec);
    }
    decoder(a) {
        return new FrameDecoder1(a, this.codec.decoder(a));
    }
}
class FrameEncoder1 {
    constructor(a, b){
        this.w = a, this.codec = b;
    }
    async encode(e) {
        const a = new Buffer1, f = this.codec.encoder(a);
        await f.encode(e);
        const d = new DataView(new ArrayBuffer(4));
        d.setUint32(0, a.length);
        const b = new Uint8Array(a.length + 4);
        b.set(new Uint8Array(d.buffer), 0), b.set(a.bytes(), 4);
        let c = 0;
        while(c < b.length)c += await this.w.write(b.subarray(c));
    }
}
class FrameDecoder1 {
    constructor(a, b){
        this.r = a, this.dec = b;
    }
    async decode(e) {
        const a = new Uint8Array(4), b = await this.r.read(a);
        if (b === null) return null;
        const c = new DataView(a.buffer), d = c.getUint32(0);
        return await this.dec.decode(d);
    }
}
function HandlerFunc1(a) {
    return {
        respondRPC: a
    };
}
function NotFoundHandler1() {
    return HandlerFunc1((a, b)=>{
        a.return(new Error(`not found: ${b.selector}`));
    });
}
function cleanSelector(a) {
    return a === "" ? "/" : (a[0] != "/" && (a = "/" + a), a = a.replace(".", "/"), a);
}
class RespondMux1 {
    constructor(){
        this.handlers = {};
    }
    async respondRPC(b, a) {
        const c = this.handler(a);
        await c.respondRPC(b, a);
    }
    handler(b) {
        const a = this.match(b.selector);
        return a || NotFoundHandler1();
    }
    remove(a) {
        a = cleanSelector(a);
        const b = this.match(a);
        return delete this.handlers[a], b || null;
    }
    match(a) {
        return a = cleanSelector(a), this.handlers.hasOwnProperty(a) ? this.handlers[a] : null;
    }
    handle(a, b) {
        if (a === "") throw "invalid selector";
        if (a = cleanSelector(a), !b) throw "invalid handler";
        if (this.match(a)) throw "selector already registered";
        this.handlers[a] = b;
    }
}
class Call1 {
    constructor(a, b){
        this.selector = a, this.decoder = b;
    }
    receive() {
        return this.decoder.decode();
    }
}
class ResponseHeader1 {
    constructor(){
        this.Error = void 0, this.Continue = !1;
    }
}
class Response {
    constructor(a, b){
        this.channel = a, this.codec = b, this.error = void 0, this.continue = !1;
    }
    send(a) {
        this.codec.encoder(this.channel).encode(a);
    }
    receive() {
        return this.codec.decoder(this.channel).decode();
    }
}
class Client1 {
    constructor(a, b){
        this.session = a, this.codec = b;
    }
    async call(b, c) {
        const a = await this.session.open();
        try {
            const e = new FrameCodec1(this.codec), f = e.encoder(a), g = e.decoder(a);
            await f.encode({
                Selector: b
            }), await f.encode(c);
            const h = await g.decode(), d = new Response(a, e);
            if (d.error = h.Error, d.error !== void 0 && d.error !== null) throw d.error;
            return d.reply = await g.decode(), d.continue = h.Continue, d.continue || await a.close(), d;
        } catch (d) {
            return await a.close(), console.error(d, b, c), Promise.reject(d);
        }
    }
}
async function Respond1(a, d, b) {
    const e = new FrameCodec1(d), f = e.decoder(a), h = await f.decode(), g = new Call1(h.Selector, f);
    g.caller = new Client1(a.session, d);
    const i = new ResponseHeader1, c = new responder1(a, e, i);
    return b || (b = new RespondMux1), await b.respondRPC(c, g), c.responded || await c.return(null), Promise.resolve();
}
function VirtualCaller1(b1) {
    function a1(b2, c1) {
        return new Proxy(Object.assign(()=>{}, {
            path: b2,
            callable: c1
        }), {
            get (b, c, d) {
                return c.startsWith("__") ? Reflect.get(b, c, d) : a1(b.path ? `${b.path}.${c}` : c, b.callable);
            },
            apply ({ path: a , callable: b  }, d, c = []) {
                return b(a, c);
            }
        });
    }
    return a1("", (a2, c)=>b1.call(a2, c).then((a)=>a.reply
        )
    );
}
class responder1 {
    constructor(a, b, c){
        this.ch = a, this.codec = b, this.header = c, this.responded = !1;
    }
    send(a) {
        return this.codec.encoder(this.ch).encode(a);
    }
    return(a) {
        return this.respond(a, !1);
    }
    async continue(a) {
        return await this.respond(a, !0), this.ch;
    }
    async respond(a, b) {
        return this.responded = !0, this.header.Continue = b, a instanceof Error && (this.header.Error = a.message, a = null), await this.send(this.header), await this.send(a), b || await this.ch.close(), Promise.resolve();
    }
}
class Peer1 {
    constructor(a, b){
        this.session = a, this.codec = b, this.caller = new Client1(a, b), this.responder = new RespondMux1;
    }
    async respond() {
        while(!0){
            const a = await this.session.accept();
            if (a === null) break;
            Respond1(a, this.codec, this.responder);
        }
    }
    async call(a, b) {
        return this.caller.call(a, b);
    }
    handle(a, b) {
        this.responder.handle(a, b);
    }
    respondRPC(a, b) {
        this.responder.respondRPC(a, b);
    }
    virtualize() {
        return VirtualCaller1(this.caller);
    }
}
function concat(c2, d) {
    const a = new Uint8Array(d);
    let b = 0;
    return c2.forEach((c)=>{
        a.set(c, b), b += c.length;
    }), a;
}
class queue {
    constructor(){
        this.q = [], this.waiters = [], this.closed = !1;
    }
    push(a) {
        if (this.closed) throw "closed queue";
        if (this.waiters.length > 0) {
            const b = this.waiters.shift();
            b && b(a);
            return;
        }
        this.q.push(a);
    }
    shift() {
        return this.closed ? Promise.resolve(null) : new Promise((a)=>{
            if (this.q.length > 0) {
                a(this.q.shift() || null);
                return;
            }
            this.waiters.push(a);
        });
    }
    close() {
        if (this.closed) return;
        this.closed = !0, this.waiters.forEach((a)=>{
            a(null);
        });
    }
}
class ReadBuffer {
    constructor(){
        this.readBuf = new Uint8Array(0), this.gotEOF = !1, this.readers = [];
    }
    read(a) {
        return new Promise((b)=>{
            let c = ()=>{
                if (this.readBuf === void 0) {
                    b(null);
                    return;
                }
                if (this.readBuf.length == 0) {
                    if (this.gotEOF) {
                        this.readBuf = void 0, b(null);
                        return;
                    }
                    this.readers.push(c);
                    return;
                }
                const d = this.readBuf.slice(0, a.length);
                this.readBuf = this.readBuf.slice(d.length), this.readBuf.length == 0 && this.gotEOF && (this.readBuf = void 0), a.set(d), b(d.length);
            };
            c();
        });
    }
    write(a) {
        for(this.readBuf && (this.readBuf = concat([
            this.readBuf,
            a
        ], this.readBuf.length + a.length)); !this.readBuf || this.readBuf.length > 0;){
            let a = this.readers.shift();
            if (!a) break;
            a();
        }
        return Promise.resolve(a.length);
    }
    eof() {
        this.gotEOF = !0, this.flushReaders();
    }
    close() {
        this.readBuf = void 0, this.flushReaders();
    }
    flushReaders() {
        while(!0){
            const a = this.readers.shift();
            if (!a) return;
            a();
        }
    }
}
const CloseID = 106;
payloadSizes = new Map([
    [
        100,
        12
    ],
    [
        101,
        16
    ],
    [
        102,
        4
    ],
    [
        103,
        8
    ],
    [
        104,
        8
    ],
    [
        105,
        4
    ],
    [
        106,
        4
    ]
]), debug4 = {
    messages: !1,
    bytes: !1
};
class Encoder {
    constructor(a){
        this.w = a;
    }
    async encode(c) {
        debug4.messages && console.log("<<ENC", c);
        const b = Marshal(c);
        debug4.bytes && console.log("<<ENC", b);
        let a = 0;
        while(a < b.length)a += await this.w.write(b.subarray(a));
        return a;
    }
}
class Decoder {
    constructor(a){
        this.r = a;
    }
    async decode() {
        const a = await readPacket(this.r);
        if (a === null) return Promise.resolve(null);
        debug4.bytes && console.log(">>DEC", a);
        const b = Unmarshal(a);
        return debug4.messages && console.log(">>DEC", b), b;
    }
}
function Marshal(a) {
    if (a.ID === 106) {
        const c = a, b = new DataView(new ArrayBuffer(5));
        return b.setUint8(0, c.ID), b.setUint32(1, c.channelID), new Uint8Array(b.buffer);
    }
    if (a.ID === 104) {
        const b = a, c = new DataView(new ArrayBuffer(9));
        c.setUint8(0, b.ID), c.setUint32(1, b.channelID), c.setUint32(5, b.length);
        const d = new Uint8Array(9 + b.length);
        return d.set(new Uint8Array(c.buffer), 0), d.set(b.data, 9), d;
    }
    if (a.ID === 105) {
        const c = a, b = new DataView(new ArrayBuffer(5));
        return b.setUint8(0, c.ID), b.setUint32(1, c.channelID), new Uint8Array(b.buffer);
    }
    if (a.ID === 100) {
        const c = a, b = new DataView(new ArrayBuffer(13));
        return b.setUint8(0, c.ID), b.setUint32(1, c.senderID), b.setUint32(5, c.windowSize), b.setUint32(9, c.maxPacketSize), new Uint8Array(b.buffer);
    }
    if (a.ID === 101) {
        const c = a, b = new DataView(new ArrayBuffer(17));
        return b.setUint8(0, c.ID), b.setUint32(1, c.channelID), b.setUint32(5, c.senderID), b.setUint32(9, c.windowSize), b.setUint32(13, c.maxPacketSize), new Uint8Array(b.buffer);
    }
    if (a.ID === 102) {
        const c = a, b = new DataView(new ArrayBuffer(5));
        return b.setUint8(0, c.ID), b.setUint32(1, c.channelID), new Uint8Array(b.buffer);
    }
    if (a.ID === 103) {
        const c = a, b = new DataView(new ArrayBuffer(9));
        return b.setUint8(0, c.ID), b.setUint32(1, c.channelID), b.setUint32(5, c.additionalBytes), new Uint8Array(b.buffer);
    }
    throw `marshal of unknown type: ${a}`;
}
async function readPacket(d) {
    const c = new Uint8Array(1), f = await d.read(c);
    if (f === null) return Promise.resolve(null);
    const b = c[0], e = payloadSizes.get(b);
    if (e === void 0 || b < 100 || b > 106) return Promise.reject(`bad packet: ${b}`);
    const a = new Uint8Array(e), g = await d.read(a);
    if (g === null) return Promise.reject("unexpected EOF");
    if (b === 104) {
        const f = new DataView(a.buffer), b = f.getUint32(4), e = new Uint8Array(b), g = await d.read(e);
        return g === null ? Promise.reject("unexpected EOF") : concat([
            c,
            a,
            e
        ], b + a.length + 1);
    }
    return concat([
        c,
        a
    ], a.length + 1);
}
function Unmarshal(b) {
    const a = new DataView(b.buffer);
    switch(b[0]){
        case 106:
            return {
                ID: b[0],
                channelID: a.getUint32(1)
            };
        case 104:
            let c = a.getUint32(5), d = new Uint8Array(b.buffer.slice(9));
            return {
                ID: b[0],
                channelID: a.getUint32(1),
                length: c,
                data: d
            };
            return {
                ID: b[0],
                channelID: a.getUint32(1),
                length: c,
                data: d
            };
        case 105:
            return {
                ID: b[0],
                channelID: a.getUint32(1)
            };
        case 100:
            return {
                ID: b[0],
                senderID: a.getUint32(1),
                windowSize: a.getUint32(5),
                maxPacketSize: a.getUint32(9)
            };
        case 101:
            return {
                ID: b[0],
                channelID: a.getUint32(1),
                senderID: a.getUint32(5),
                windowSize: a.getUint32(9),
                maxPacketSize: a.getUint32(13)
            };
        case 102:
            return {
                ID: b[0],
                channelID: a.getUint32(1)
            };
        case 103:
            return {
                ID: b[0],
                channelID: a.getUint32(1),
                additionalBytes: a.getUint32(5)
            };
        default:
            throw `unmarshal of unknown type: ${b[0]}`;
    }
}
const channelMaxPacket1 = 1 << 15, maxPacketLength1 = Number.MAX_VALUE, channelWindowSize1 = 64 * channelMaxPacket1;
class Channel1 {
    constructor(a){
        this.localId = 0, this.remoteId = 0, this.maxIncomingPayload = 0, this.maxRemotePayload = 0, this.sentEOF = !1, this.sentClose = !1, this.remoteWin = 0, this.myWindow = 0, this.ready = new queue, this.session = a, this.writers = [], this.readBuf = new ReadBuffer;
    }
    ident() {
        return this.localId;
    }
    async read(b) {
        let a = await this.readBuf.read(b);
        if (a !== null) try {
            await this.adjustWindow(a);
        } catch (a3) {
            if (a3 !== "EOF") throw a3;
        }
        return a;
    }
    write(a) {
        return this.sentEOF ? Promise.reject("EOF") : new Promise((d, e)=>{
            let b = 0;
            const c = ()=>{
                if (this.sentEOF || this.sentClose) {
                    e("EOF");
                    return;
                }
                if (a.byteLength == 0) {
                    d(b);
                    return;
                }
                const h = Math.min(this.maxRemotePayload, a.length), g = this.reserveWindow(h);
                if (g == 0) {
                    this.writers.push(c);
                    return;
                }
                const f = a.slice(0, g);
                this.send({
                    ID: 104,
                    channelID: this.remoteId,
                    length: f.length,
                    data: f
                }).then(()=>{
                    if (b += f.length, a = a.slice(f.length), a.length == 0) {
                        d(b);
                        return;
                    }
                    this.writers.push(c);
                });
            };
            c();
        });
    }
    reserveWindow(a) {
        return this.remoteWin < a && (a = this.remoteWin), this.remoteWin -= a, a;
    }
    addWindow(a) {
        for(this.remoteWin += a; this.remoteWin > 0;){
            const a = this.writers.shift();
            if (!a) break;
            a();
        }
    }
    async closeWrite() {
        this.sentEOF = !0, await this.send({
            ID: 105,
            channelID: this.remoteId
        }), this.writers.forEach((a)=>a()
        ), this.writers = [];
    }
    async close() {
        if (!this.sentClose) {
            for(await this.send({
                ID: 106,
                channelID: this.remoteId
            }), this.sentClose = !0; await this.ready.shift() !== null;);
            return;
        }
        this.shutdown();
    }
    shutdown() {
        this.readBuf.close(), this.writers.forEach((a)=>a()
        ), this.ready.close(), this.session.rmCh(this.localId);
    }
    async adjustWindow(a) {
        this.myWindow += a, await this.send({
            ID: 103,
            channelID: this.remoteId,
            additionalBytes: a
        });
    }
    send(a) {
        if (this.sentClose) throw "EOF";
        return this.sentClose = a.ID === CloseID, this.session.enc.encode(a);
    }
    handle(a) {
        if (a.ID === 104) {
            this.handleData(a);
            return;
        }
        if (a.ID === 106) {
            this.close();
            return;
        }
        if (a.ID === 105 && this.readBuf.eof(), a.ID === 102) {
            this.session.rmCh(a.channelID), this.ready.push(!1);
            return;
        }
        if (a.ID === 101) {
            if (a.maxPacketSize < 9 || a.maxPacketSize > maxPacketLength1) throw "invalid max packet size";
            this.remoteId = a.senderID, this.maxRemotePayload = a.maxPacketSize, this.addWindow(a.windowSize), this.ready.push(!0);
            return;
        }
        a.ID === 103 && this.addWindow(a.additionalBytes);
    }
    handleData(a) {
        if (a.length > this.maxIncomingPayload) throw "incoming packet exceeds maximum payload size";
        if (this.myWindow < a.length) throw "remote side wrote too much";
        this.myWindow -= a.length, this.readBuf.write(a.data);
    }
}
class Session1 {
    constructor(a){
        this.conn = a, this.enc = new Encoder(a), this.dec = new Decoder(a), this.channels = [], this.incoming = new queue, this.done = this.loop();
    }
    async open() {
        const a = this.newChannel();
        if (a.maxIncomingPayload = channelMaxPacket1, await this.enc.encode({
            ID: 100,
            windowSize: a.myWindow,
            maxPacketSize: a.maxIncomingPayload,
            senderID: a.localId
        }), await a.ready.shift()) return a;
        throw "failed to open";
    }
    accept() {
        return this.incoming.shift();
    }
    async close() {
        for (const b of Object.keys(this.channels)){
            const a = parseInt(b);
            this.channels[a] !== void 0 && this.channels[a].shutdown();
        }
        this.conn.close(), await this.done;
    }
    async loop() {
        try {
            while(!0){
                const a = await this.dec.decode();
                if (a === null) {
                    this.close();
                    return;
                }
                if (a.ID === 100) {
                    await this.handleOpen(a);
                    continue;
                }
                const b = a, c = this.getCh(b.channelID);
                if (c === void 0) throw `invalid channel (${b.channelID}) on op ${b.ID}`;
                await c.handle(b);
            }
        } catch (a) {
            throw new Error(`session loop: ${a}`);
        }
    }
    async handleOpen(b) {
        if (b.maxPacketSize < 9 || b.maxPacketSize > maxPacketLength1) {
            await this.enc.encode({
                ID: 102,
                channelID: b.senderID
            });
            return;
        }
        const a = this.newChannel();
        a.remoteId = b.senderID, a.maxRemotePayload = b.maxPacketSize, a.remoteWin = b.windowSize, a.maxIncomingPayload = channelMaxPacket1, this.incoming.push(a), await this.enc.encode({
            ID: 101,
            channelID: a.remoteId,
            senderID: a.localId,
            windowSize: a.myWindow,
            maxPacketSize: a.maxIncomingPayload
        });
    }
    newChannel() {
        const a = new Channel1(this);
        return a.remoteWin = 0, a.myWindow = channelWindowSize1, a.localId = this.addCh(a), a;
    }
    getCh(b) {
        const a = this.channels[b];
        return a && a.localId !== b && console.log("bad ids:", b, a.localId, a.remoteId), a;
    }
    addCh(a) {
        return this.channels.forEach((c, b)=>{
            if (c === void 0) return this.channels[b] = a, b;
        }), this.channels.push(a), this.channels.length - 1;
    }
    rmCh(a) {
        delete this.channels[a];
    }
}
function connect2(b, a) {
    return new Promise((d)=>{
        const c = new WebSocket(b);
        c.onopen = ()=>d(new Conn(c))
        , a && (c.onclose = a);
    });
}
class Conn {
    constructor(b3){
        this.isClosed = !1, this.waiters = [], this.chunks = [], this.ws = b3, this.ws.binaryType = "arraybuffer", this.ws.onmessage = (a)=>{
            const b = new Uint8Array(a.data);
            if (this.chunks.push(b), this.waiters.length > 0) {
                const a = this.waiters.shift();
                a && a();
            }
        };
        const a4 = this.ws.onclose;
        this.ws.onclose = (b)=>{
            a4 && a4.bind(this.ws)(b), this.close();
        };
    }
    read(a5) {
        return new Promise((b)=>{
            var c3 = ()=>{
                if (this.isClosed) {
                    b(null);
                    return;
                }
                if (this.chunks.length === 0) {
                    this.waiters.push(c3);
                    return;
                }
                let d = 0;
                while(d < a5.length){
                    const c = this.chunks.shift();
                    if (c === null || c === void 0) {
                        b(null);
                        return;
                    }
                    const e = c.slice(0, a5.length - d);
                    if (a5.set(e, d), d += e.length, c.length > e.length) {
                        const a = c.slice(e.length);
                        this.chunks.unshift(a);
                    }
                }
                b(d);
            };
            c3();
        });
    }
    write(a) {
        return this.ws.send(a), Promise.resolve(a.byteLength);
    }
    close() {
        if (this.isClosed) return;
        this.isClosed = !0, this.waiters.forEach((a)=>a()
        ), this.ws.close();
    }
}
const mod = {
    connect: connect2,
    Conn
};
frames = {};
function frameElementID(a) {
    return a.frameElement ? a.frameElement.id : "";
}
window.addEventListener("message", (b)=>{
    if (!b.source) return;
    const a = frameElementID(b.source);
    if (!frames[a]) {
        const b = new CustomEvent("connection", {
            detail: a
        });
        if (!window.dispatchEvent(b)) return;
        if (!frames[a]) {
            console.warn("incoming message with no connection for frame ID in window:", a, window.location);
            return;
        }
    }
    const c = frames[a], d = new Uint8Array(b.data);
    if (c.chunks.push(d), c.waiters.length > 0) {
        const a = c.waiters.shift();
        a && a();
    }
});
class Conn1 {
    constructor(a){
        this.isClosed = !1, this.waiters = [], this.chunks = [], a && a.contentWindow ? (this.frame = a.contentWindow, frames[a.id] = this) : (this.frame = window.parent, frames[frameElementID(window.parent)] = this);
    }
    read(a6) {
        return new Promise((b)=>{
            var c4 = ()=>{
                if (this.isClosed) {
                    b(null);
                    return;
                }
                if (this.chunks.length === 0) {
                    this.waiters.push(c4);
                    return;
                }
                let d = 0;
                while(d < a6.length){
                    const c = this.chunks.shift();
                    if (c === null || c === void 0) {
                        b(null);
                        return;
                    }
                    const e = c.slice(0, a6.length - d);
                    if (a6.set(e, d), d += e.length, c.length > e.length) {
                        const a = c.slice(e.length);
                        this.chunks.unshift(a);
                    }
                }
                b(d);
            };
            c4();
        });
    }
    write(a) {
        return this.frame.postMessage(a.buffer, "*"), Promise.resolve(a.byteLength);
    }
    close() {
        if (this.isClosed) return;
        this.isClosed = !0, this.waiters.forEach((a)=>a()
        );
    }
}
options1 = {
    transport: mod
};
async function connect1(a, b) {
    const c = await options1.transport.connect(a);
    return open1(c, b);
}
function open1(a, d, b) {
    a === window.parent && (a = new Conn1), typeof a == "string" && (a = new Conn1(document.querySelector(`iframe#${a}`)));
    const e = new Session1(a), c = new Peer1(e, d);
    if (b) {
        for(const a in b)c.handle(a, HandlerFunc1(b[a]));
        c.respond();
    }
    return c;
}
(()=>{
    if (window) {
        window.requestAnimationFrame(async ()=>{
            window["$host"] = await connect(`ws://${window.location.host}/`);
        });
    }
})();
async function connect(url) {
    return new Client(await connect1(url, new JSONCodec1()));
}
class Client {
    rpc;
    constructor(peer){
        this.rpc = peer.virtualize();
    }
    get screen() {
        return new Screen(this.rpc);
    }
    get shell() {
        return new Shell(this.rpc);
    }
    get window() {
        return {
            New: async (options)=>{
                const w = await this.rpc.window.New(options);
                return new Window(this.rpc, w.ID);
            }
        };
    }
}
class Screen {
    rpc;
    constructor(rpc){
        this.rpc = rpc;
    }
    Displays() {
        return this.rpc.screen.Displays();
    }
}
class Shell {
    rpc;
    constructor(rpc){
        this.rpc = rpc;
    }
    ShowNotification(n) {
        this.rpc.shell.ShowNotification(n);
    }
    ShowMessage(msg) {
        this.rpc.shell.ShowMessage(msg);
    }
    ShowFilePicker(fd) {
        return this.rpc.shell.ShowFilePicker(fd);
    }
    ReadClipboard() {
        return this.rpc.shell.ReadClipboard();
    }
    WriteClipboard(text) {
        return this.rpc.shell.WriteClipboard(text);
    }
    RegisterShortcut(accelerator) {
        return this.rpc.shell.RegisterShortcut(accelerator);
    }
    IsShortcutRegistered(accelerator) {
        return this.rpc.shell.IsShortcutRegistered(accelerator);
    }
    UnregisterShortcut(accelerator) {
        return this.rpc.shell.UnregisterShortcut(accelerator);
    }
    UnregisterAllShortcuts() {
        return this.rpc.shell.UnregisterAllShortcuts();
    }
}
class Window {
    ID;
    rpc;
    constructor(rpc, id){
        this.rpc = rpc;
        this.ID = id;
    }
    async destroy() {
        await this.rpc.window.Destroy(this.ID);
    }
    async focus() {
        await this.rpc.window.Focus(this.ID);
    }
    async getOuterPosition() {
        return await this.rpc.window.GetOuterPosition(this.ID);
    }
    async getOuterSize() {
        return await this.rpc.window.GetOuterSize(this.ID);
    }
    async isDestroyed() {
        return await this.rpc.window.IsDestroyed(this.ID);
    }
    async isVisible() {
        return await this.rpc.window.IsVisible(this.ID);
    }
    async setVisible(visible) {
        return await this.rpc.window.SetVisible(this.ID, visible);
    }
    async setMaximized(maximized) {
        return await this.rpc.window.SetMaximized(this.ID, maximized);
    }
    async setMinimized(minimized) {
        return await this.rpc.window.SetMinimized(this.ID, minimized);
    }
    async setFullscreen(fullscreen) {
        return await this.rpc.window.SetFullscreen(this.ID, fullscreen);
    }
    async setMinSize(size) {
        return await this.rpc.window.SetMinSize(this.ID, size);
    }
    async setMaxSize(size) {
        return await this.rpc.window.SetMaxSize(this.ID, size);
    }
    async setResizable(resizable) {
        return await this.rpc.window.SetResizable(this.ID, resizable);
    }
    async setAlwaysOnTop(always) {
        return await this.rpc.window.SetAlwaysOnTop(this.ID, always);
    }
    async setSize(size) {
        return await this.rpc.window.SetSize(this.ID, size);
    }
    async setPosition(position) {
        return await this.rpc.window.SetPosition(this.ID, position);
    }
    async setTitle(title) {
        return await this.rpc.window.SetTitle(this.ID, title);
    }
}
export { connect as connect };
export { Client as Client };
export { Window as Window };

