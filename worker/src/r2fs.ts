import { createTarDecoder } from 'modern-tar';

const MAX_RETRY_DELAY = 1000; // Cap retry delay
const RETRY_JITTER_MAX = 11; // Maximum jitter for retry backoff
const BOUNDARY = "entry-boundary";
export const METHODS = "GET, HEAD, PUT, PATCH, DELETE, MOVE, COPY, OPTIONS";

export async function handle(req, env, basepath="") {
    const url = new URL(req.url);
    let wantsIndex = false;
    if (url.pathname.endsWith("/...")) {
        url.pathname = url.pathname.slice(0, -4);
        wantsIndex = true;
    }
    if (basepath === "/") {
        basepath = "";
    }
    // path is normalized path without basepath
    const path = (url.pathname.slice(basepath.length) === "") ? "/" : url.pathname.slice(basepath.length);
    // key is path without trailing slash
    const key = (path !== "/" && path.endsWith("/")) ? path.slice(0, -1) : path;
    
    switch (req.method) {
        case "PATCH":
            return handlePatch(req, env, key);
        case "PUT":
            return handlePut(req, env, key, path);
        case "HEAD":
            return handleHead(req, env, key, basepath);
        case "GET":
            return handleGet(req, env, key, basepath, wantsIndex);
        case "DELETE":
            return handleDelete(req, env, key);
        case "MOVE":
        case "COPY":
            return handleMoveCopy(req, env, key, basepath);
        case "OPTIONS":
            return handleOptions();
        default:
            return new Response("Method Not Allowed\n", {
                status: 405,
                headers: { Allow: METHODS },
            });
    }
}

export async function handlePatch(req, env, key) {
    if (isAttr(key)) {
        //return new Response("Attributes cannot be patched\n", { status: 400 });
        return new Response("NoOp\n", {status: 200});
    }

    let customMetadata = {};
    
    if (req.headers.get("Content-Type")?.includes("application/x-tar")) {
        // tar patch
        const decoder = createTarDecoder();
        const entriesStream = req.body.pipeThrough(decoder);
        return new Response(new ReadableStream({
            async start(controller) {

                const puts: Promise<void>[] = [];
                const deleteBatch: string[] = [];
                const deletes: Promise<void>[] = [];

                for await (const entry of entriesStream) {
                    let mode: number = entry.header.mode || 0;
                    switch (entry.header.type) {
                        case "directory":
                            // S_IFDIR is 0o040000 (octal), which is 16384 (decimal)
                            mode = mode | 0o040000;
                            customMetadata["Content-Type"] = "application/x-directory";
                            break;
                        case "file":
                            customMetadata["Content-Type"] = "application/octet-stream";
                            break;
                        case "symlink":
                            // S_IFLNK is 0o120000 (octal), which is 65536 (decimal)
                            mode = mode | 0o120000;
                            customMetadata["Content-Type"] = "application/x-symlink";
                            break;
                    }
                    customMetadata["Content-Mode"] = mode.toString();
                    customMetadata["Content-Length"] = (entry.header.type === "symlink")
                        ? entry.header.linkname.length.toString() 
                        : entry.header.size.toString();
                    customMetadata["Content-Ownership"] = "0:0";
                    if (entry.header.mtime !== undefined) {
                        customMetadata["Content-Modified"] = Math.floor(entry.header.mtime.getTime() / 1000).toString();
                    }
                    
                    let entryName = entry.header.name.replace(/^\/|\/$/g, '');
                    if (entryName === "./") {
                        entryName = ".";
                    } else if (entryName.startsWith("./")) {
                        entryName = entryName.slice(2);
                    }
                    const entryKey = (entryName === ".") 
                        ? key 
                        : (key !== "/" ? [key, entryName].join("/") : "/" + entryName);
                    
                    if (entry.header.pax?.["delete"] === "") {
                        deleteBatch.push(entryKey);
                        controller.enqueue(encode("- "+entryKey+"\n"));
                        if (deleteBatch.length === 1000) {
                            deletes.push(env.bucket.delete(deleteBatch));
                            deleteBatch.length = 0;
                        }
                    } else {
                        let buf = await new Response(entry.body).arrayBuffer();
                        if (entry.header.type === "symlink") {
                            buf = encode(entry.header.linkname);
                        }
                        puts.push(env.bucket.put(entryKey, buf, {
                            httpMetadata: {"Content-Type": customMetadata["Content-Type"]},
                            customMetadata,
                        }));
                        controller.enqueue(encode("+ "+entryKey+"\n"));
                    }
                }

                if (deleteBatch.length > 0) {
                    deletes.push(env.bucket.delete(deleteBatch));
                    deleteBatch.length = 0;
                }
                await Promise.all(puts);
                await Promise.all(deletes);
                controller.close();
            }
        }), { headers: { "Content-Type": `application/x-tar-apply` } });
    }
    
    // regular metadata update
    headersToMetadata(req.headers, customMetadata);
    // logMeta(customMetadata);
    if (!await compareAndSwap(env.bucket, key, async (object) => {
        customMetadata = Object.assign(object.customMetadata || {}, customMetadata);
        return {
            httpMetadata: req.headers,
            customMetadata,
        };
    })) {
        return new Response("Object Not Found\n", { status: 404 });
    }

    return new Response("OK\n");
}

export async function handlePut(req, env, key, path) {
    let customMetadata = {};
    let attrKey = undefined;
    if (isAttr(key)) {
        [key, attrKey] = splitAttr(key);
        const object = await env.bucket.get(key);
        if (object === null) {
            return new Response("Object Not Found\n", { status: 404 });
        }
        customMetadata = object.customMetadata || {};
        // Use microseconds for Change-Timestamp
        customMetadata["Change-Timestamp"] = (Date.now() * 1000).toString();
        const attrValue = new TextDecoder().decode(await req.arrayBuffer());
        customMetadata["Attribute-" + attrKey] = attrValue.replace(/\n+$/, "");
        await env.bucket.put(key, object.body, {
            httpMetadata: object.httpMetadata || {},
            customMetadata,
        });
        return new Response("OK\n");
    }

    headersToMetadata(req.headers, customMetadata);

    if (!await compareAndSwap(env.bucket, key, async (object) => {
        const existingTimestamp = parseMicroseconds(object.customMetadata["Change-Timestamp"]);
        const newTimestamp = parseMicroseconds(customMetadata["Change-Timestamp"]);
        if (newTimestamp > existingTimestamp) {
            customMetadata = Object.assign(object.customMetadata || {}, customMetadata);
        } else {
            // keep existing metadata if timestamp is not newer
            customMetadata = object.customMetadata;
        }
        return {
            value: req.body,
            httpMetadata: req.headers,
            customMetadata,
        };
    })) {
        // new object
        const defaultModeDir = "16868"; // default dir mode is 0744 plus dir flag 0040000
        const defaultModeFile = "33188"; // default file mode is 0644 plus file flag 0100000

        let defaultMode = defaultModeFile;
        if (path.endsWith("/")) {
            customMetadata["Content-Type"] = "application/x-directory";
            defaultMode = defaultModeDir;
        }
        customMetadata["Content-Mode"] ||= defaultMode;
        customMetadata["Content-Ownership"] ||= "0:0";
        customMetadata["Content-Modified"] ||= Math.floor(Date.now() / 1000).toString();
        
        await env.bucket.put(key, req.body, {
            httpMetadata: req.headers,
            customMetadata,
        });
    }

    return new Response("OK\n");
}

export async function handleHead(req, env, key, basepath) {
    let attrKey = undefined;
    if (isAttr(key)) {
        if (isAttrs(key)) {
            return new Response(null, {
                status: 200,
                headers: {"Content-Type": "application/x-directory"},
            });
        }
        [key, attrKey] = splitAttr(key);
    }

    const object = await env.bucket.get(key, {
        onlyIf: req.headers,
        range: req.headers,
    });

    if (object === null) {
        return new Response("Object Not Found\n", { status: 404 });
    }

    // ATTRIBUTES
    if (attrKey) {
        if (object.customMetadata["Attribute-" + attrKey] === undefined) {
            return new Response("Attribute Not Found\n", { status: 404 });
        }
        return new Response(null, {
            status: 200,
            headers: {"Content-Type": "plain/text"},
        });
    }

    return new Response(null, {
        status: 200,
        headers: headersFromObject(object, basepath),
    });
}

export async function handleGet(req, env, key, basepath, wantsIndex) {
    let attrKey = undefined;
    let wantsAttrs = false;
    if (isAttr(key)) {
        if (isAttrs(key)) {
            wantsAttrs = true;
            key = key.slice(0, -6); // remove /:attr
        } else {
            [key, attrKey] = splitAttr(key);
        }
    }

    const object = await env.bucket.get(key, {
        onlyIf: req.headers,
        range: req.headers,
    });

    if (object === null) {
        return new Response("Object Not Found\n", { status: 404 });
    }

    // ATTRIBUTES
    if (attrKey) {
        const attr = object.customMetadata["Attribute-" + attrKey];
        if (attr === undefined) {
            return new Response("Attribute Not Found\n", { status: 404 });
        }
        return new Response(attr+"\n", {
            status: 200,
            headers: {
                "Content-Type": "plain/text",
            },
        });
    }
    if (wantsAttrs) {
        const attrs = await attributeEntries(object);
        return new Response(attrs, {
            status: 200,
            headers: {
                "Content-Type": "application/x-directory",
            },
        });
    }

    // FILES
    if (object.customMetadata["Content-Type"] !== "application/x-directory") {
        // When no body is present, preconditions have failed
        return new Response("body" in object ? object.body : undefined, {
            status: "body" in object ? 200 : 412,
            headers: headersFromObject(object, basepath),
        });
    }

    // DIRECTORIES
    if (wantsIndex) {
        const objects = await objectsInDir(env.bucket, key, -1);
        return new Response(new ReadableStream({
            async start(controller) {
                controller.enqueue(encode(
                    `--${BOUNDARY}\r\n` +
                    `${formatHeaders(headersFromObject(object, basepath))}\r\n` +
                    `${formatEntries(entriesFromObjects(objects, key))}\r\n`
                ));
                for (const obj of objects) {
                    controller.enqueue(encode(
                        `--${BOUNDARY}\r\n` +
                        `${formatHeaders(headersFromObject(obj, basepath, !isObjectDir(obj)))}\r\n` +
                        ((isObjectDir(obj)) ? `${formatEntries(entriesFromObjects(objects, obj.key))}\r\n` : "")
                    ));
                }
                controller.enqueue(encode(`--${BOUNDARY}--\r\n`));
                controller.close();
            }
        }), { headers: { "Content-Type": `multipart/mixed; boundary=${BOUNDARY}` } });
    }

    // multipart/mixed listing gives metadata for all files in the directory
    // important for caching
    if (req.headers.get("Accept")?.includes("multipart/mixed")) {
        const objects = await objectsInDir(env.bucket, key, 2);
        return new Response(new ReadableStream({
            async start(controller) {
                controller.enqueue(encode(
                    `--${BOUNDARY}\r\n` +
                    `${formatHeaders(headersFromObject(object, basepath))}\r\n` +
                    `${formatEntries(entriesFromObjects(objects, key))}\r\n`
                ));
                for (const obj of objects) {
                    if (obj.key.slice(key.length).split("/").length > (key === "/" ? 1 : 2)) {
                        continue;
                    }
                    controller.enqueue(encode(
                        `--${BOUNDARY}\r\n` +
                        `${formatHeaders(headersFromObject(obj, basepath, !isObjectDir(obj)))}\r\n` +
                        ((isObjectDir(obj)) ? `${formatEntries(entriesFromObjects(objects, obj.key))}\r\n` : "")
                    ));
                }
                controller.enqueue(encode(`--${BOUNDARY}--\r\n`));
                controller.close();
            }
        }), { headers: { "Content-Type": `multipart/mixed; boundary=${BOUNDARY}` } });
    }
    
    // normal listing
    const listing = await directoryEntries(env.bucket, key);
    const headers = headersFromObject(object, basepath);
    headers.delete("Content-Length");
    headers.delete("ETag");
    return new Response(listing, {headers});
}

export async function handleDelete(req, env, key) {
    let attrKey = undefined;
    if (isAttr(key)) {
        [key, attrKey] = splitAttr(key);
    }

    const object = await env.bucket.get(key);
    if (object === null) {
        return new Response("Object Not Found\n", { status: 404 });
    }

    if (attrKey) {
        const customMetadata = object.customMetadata || {};
        delete customMetadata["Attribute-" + attrKey];
        customMetadata["Change-Timestamp"] = Math.floor(Date.now() / 1000).toString();
        await env.bucket.put(key, object.body, {
            httpMetadata: object.httpMetadata || {},
            customMetadata,
        });
        return new Response("OK\n");
    }

    // If the object is a directory, recursively delete its contents
    if (object.customMetadata["Content-Type"] === "application/x-directory") {
        const objects = await objectsInDir(env.bucket, key, -1);
        const batchSize = 1000;
        for (let i = 0; i < objects.length; i += batchSize) {
            const batch = objects.slice(i, i + batchSize);
            await env.bucket.delete(batch.map(obj => obj.key));
        }
    }

    await env.bucket.delete(key);

    return new Response("OK\n");
}

export async function handleMoveCopy(req, env, key, basepath) {
    if (isAttr(key)) {
        return new Response("Attributes cannot be moved/copied\n", { status: 400 });
    }

    // Parse source and destination from headers
    const dest = req.headers.get("Destination").slice(basepath.length);
    if (!dest || !dest.startsWith("/")) {
        return new Response("Missing or invalid Destination header\n", { status: 400 });
    }
    const destKey = (dest !== "/" && dest.endsWith("/")) ? dest.slice(0, -1) : dest;

    if (destKey === key) {
        return new Response("Cannot move/copy to same path\n", { status: 400 });
    }

    // Prevent moving/copying to root
    if (dest === "/") {
        return new Response("Cannot move/copy to root\n", { status: 400 });
    }

    // Get source object
    const srcObject = await env.bucket.get(key, { type: "stream" });
    if (!srcObject) {
        return new Response("Source Not Found\n", { status: 404 });
    }

    // If destination exists and overwrite is not allowed, fail
    const overwrite = req.headers.get("Overwrite")?.toLowerCase() !== "f";
    const destExists = await env.bucket.get(destKey);
    if (destExists && !overwrite) {
        return new Response("Destination Exists\n", { status: 412 });
    }

    // Write to destination
    await env.bucket.put(destKey, srcObject.body, {
        httpMetadata: srcObject.httpMetadata || {},
        customMetadata: srcObject.customMetadata || {},
    });

    // If the source is a directory, recursively copy its contents
    if (srcObject.customMetadata["Content-Type"] === "application/x-directory") {
        const objects = await objectsInDir(env.bucket, key, -1);
        const batchSize = 256;
        const deleteSubKeys: string[] = [];
        for (let i = 0; i < objects.length; i += batchSize) {
            const batch = objects.slice(i, i + batchSize);
            const ops: Promise<void>[] = [];
            for (const obj of batch) {
                ops.push(new Promise(async (resolve) => {
                    // Determine the new key under the destination directory
                    const newKey = destKey + obj.key.slice(key.length);
                    // Read the source object (as stream to avoid memory bloat)
                    const srcObj = await env.bucket.get(obj.key, { type: "stream" });
                    if (srcObj) {
                        await env.bucket.put(newKey, srcObj.body, {
                            httpMetadata: srcObj.httpMetadata || {},
                            customMetadata: srcObj.customMetadata || {},
                        });
                    }
                    deleteSubKeys.push(obj.key);
                    resolve();
                }));
            }
            await Promise.all(ops);
        }
        // If MOVE, delete the sources
        if (req.method === "MOVE") {
            const batchSize = 1000;
            for (let i = 0; i < deleteSubKeys.length; i += batchSize) {
                const batch = deleteSubKeys.slice(i, i + batchSize);
                await env.bucket.delete(batch);
            }
            await env.bucket.delete(key);
        }
    }

    return new Response("OK\n");
}

export function handleOptions() {
    return new Response("OK\n", {
        headers: { Allow: METHODS },
    });
}

function encode(str) {
    return new TextEncoder().encode(str);
}

function isAttr(path: string): boolean {
    return path.includes("/:attr/") || path.endsWith("/:attr");
}

function isAttrs(path: string): boolean {
    return path.endsWith("/:attr");
}

function splitAttr(path: string): string[] {
    return path.split("/:attr/");
}

function isObjectDir(object: R2Object): boolean {
    return (
        object.customMetadata?.["Content-Type"] === "application/x-directory" ||
        object.httpMetadata?.["Content-Type"]?.[0] === "application/x-directory" ||
        isDirectoryMode(object.customMetadata?.["Content-Mode"] || "33188")
    );
}

function entriesFromObjects(objects: R2Object[], dir: string): Map<string, string> {
    if (dir === "/") {
        dir = "";
    }
    const entries = new Map<string, string>();
    for (const obj of objects.filter(obj => obj.key.startsWith(dir+"/") && !obj.key.slice(dir.length+1).includes("/"))) {
        const defaultMode = (obj.httpMetadata?.["Content-Type"]?.[0] === "application/x-directory") ? "16877" : "33188";
        const entryName = obj.key.slice(dir.length+1);
        entries.set(entryName, obj.customMetadata?.["Content-Mode"] || defaultMode);
    }
    return entries;
}

function formatHeaders(headers: Headers): string {
    let result = "";
    headers.forEach((value, key) => {
        result += `${key}: ${value}\r\n`;
    });
    return result;
}


function headersFromObject(object, basepath, useEmptyRange = false): Headers {
    const headers = new Headers();
    object.writeHttpMetadata(headers);
    headers.set("ETag", object.httpEtag);
    headers.set("Content-Type", object.customMetadata["Content-Type"] || "application/octet-stream");
    headers.set("Content-Location", basepath + object.key);
    if (useEmptyRange) {
        headers.set("Content-Range", `bytes 0-0/${object.size}`);
    } else {
        headers.set("Content-Length", object.size.toString());
    }
    if (object.customMetadata) {
        for (const [k, v] of Object.entries(object.customMetadata)) {
            headers.set(k, v as string);
        }
    }
    return headers;
}

function headersToMetadata(headers, metadata) {
    [
        "Change-Timestamp",
        "Content-Type",
        "Content-Mode",
        "Content-Modified",
        "Content-Ownership",
    ].forEach(header => {
        if (headers.has(header.toLowerCase())) {
            metadata[header] = headers.get(header.toLowerCase());
        }
    });
    // add any attributes
    headers.forEach((value, key) => {
        if (key.toLowerCase().startsWith("attribute-")) {
            metadata["Attribute-" + key.slice(10)] = value;
        }
    });
}

async function objectsInDir(bucket: any, dir: string, depth: number = 1): Promise<R2Object[]> {
    const objects: R2Object[] = [];
    const prefix = dir === "/" ? "/" : dir + "/";
    let cursor: string | undefined = undefined;
    do {
        const page = await bucket.list({
            prefix,
            include: ["customMetadata", "httpMetadata"],
            cursor,
            limit: 1000,
        });

        for (const obj of page.objects || []) {
            const name = obj.key.slice(prefix.length);
            if (!name) {
                continue;
            }
            const partCount = name.split("/").length;
            if (depth !== -1 && partCount > depth) {
                continue;
            }
            objects.push(obj as R2Object);
        }

        cursor = page.truncated ? page.cursor : undefined;
    } while (cursor);

    return objects;
}

// Build a directory listing string from bucket.list() results for a given prefix
async function directoryEntries(bucket: any, key: string): Promise<string> {
    const prefix = key === "/" ? "/" : key + "/";
    const entries = new Map<string, string>();
    let cursor: string | undefined = undefined;
    do {
        const page = await bucket.list({
            prefix,
            delimiter: "/",
            include: ["customMetadata"],
            cursor,
            limit: 1000,
        });

        // Files directly under prefix
        for (const obj of page.objects || []) {
            const name = obj.key.slice(prefix.length);
            if (!name || name.includes("/")) {
                continue;
            }
            const mode = obj.customMetadata?.["Content-Mode"] || "33188"; // default file mode 0644|0100000
            entries.set(name, String(mode));
        }

        cursor = page.truncated ? page.cursor : undefined;
    } while (cursor);

    return formatEntries(new Map([...entries.entries()].sort()));
}

export async function getAttrs(bucket: any, key: string): Promise<Record<string, string>|null> {
    const object = await bucket.get(key);
    if (object === null) {
        console.log("no object", key)
        return null;
    }
    const attrs = {};
    for (const [k, v] of Object.entries(object.customMetadata || {})) {
        if (k.startsWith("Attribute-")) {
            attrs[k.slice(10)] = v;
        }
    }
    return attrs;
}

async function attributeEntries(object: R2Object): Promise<string> {
    const entries = new Map<string, string>();
    for (const [k, v] of Object.entries(object.customMetadata)) {
        if (k.startsWith("Attribute-")) {
            entries.set(k.slice(10), "33188"); // default file mode is 0644 plus file flag 0100000
        }
    }
    return formatEntries(new Map([...entries.entries()].sort()));
}

// Helper function to check if a mode string represents a directory
function isDirectoryMode(modeStr: string): boolean {
    const mode = parseInt(modeStr, 10);
    // Check for directory flag (0040000 in octal, 16384 in decimal)
    return (mode & 0o040000) !== 0;
}

// updateFn should return putOptions with optional value to set the key value
async function compareAndSwap(bucket, key, updateFn, maxRetries = 3) {
    for (let i = 0; i < maxRetries; i++) {
      // Get current state with its ETag
      const object = await bucket.get(key);
      if (object === null) {
        return false;
      }
      const currentETag = object?.etag;
      
      // Apply update
      const putOptions = await updateFn(object);
      
      // Try to write ONLY if ETag matches (object unchanged)
      try {
        let newValue;
        if (putOptions["value"]) {
            newValue = putOptions["value"];
            delete putOptions["value"];
        }
        if (!putOptions["httpMetadata"]) {
            putOptions["httpMetadata"] = {};
        }
        if (currentETag) {
          putOptions["httpMetadata"]["ifMatch"] = currentETag;
        } else {
          // For new objects, use if-none-match: * to ensure object doesn't exist
          putOptions["httpMetadata"]["ifNoneMatch"] = '*';
        }
        
        await bucket.put(key, newValue || object?.body, putOptions);
        return true;
        
      } catch (error) {
        if (error.status === 412 && i < maxRetries - 1) {
          // Object was modified, retry with optimized backoff
          const jitter = Math.random() * RETRY_JITTER_MAX; // Add jitter to prevent thundering herd
          const baseDelay = Math.pow(2, i) * 20;
          const delay = Math.min(MAX_RETRY_DELAY, baseDelay + jitter);
          await new Promise(r => setTimeout(r, delay));
          continue;
        }
        throw error;
      }
    }
}

function formatEntries(entries: Map<string, string>): string {
    return Array.from(entries.entries()).map(([name, mode]) => `${name} ${mode}`).join("\n")+"\n";
}


function parseMicroseconds(str: string): number {
    // Parses an integer string representing microseconds
    const micros = parseInt(str, 10);
    return isNaN(micros) ? 0 : micros;
}

function logMeta(object: Object) {
    console.log("Content-Type:", object["Content-Type"]);
    console.log("Content-Length:", object["Content-Length"]);
    console.log("Content-Mode:", object["Content-Mode"]);
    console.log("Content-Modified:", object["Content-Modified"]);
    console.log("Content-Ownership:", object["Content-Ownership"]);
}