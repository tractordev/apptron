import { HOST_DOMAIN } from "./config";
import { handle as handleR2FS } from "./r2fs";

export function isLocal(env: any) {
    return !!(env && env.LOCALHOST);
}

export function redirectToSignin(env: any, url: URL) {
    if (isLocal(env)) {
        url.host = env.LOCALHOST;
    } else {
        url.host = HOST_DOMAIN;
    }
    url.pathname = "/signin";
    return Response.redirect(url.toString(), 307);
}

export function insertMeta(resp: Response, meta: Record<string, string>) {
    return new HTMLRewriter().on('head', {
        element(element) {
            for (const [name, content] of Object.entries(meta)) {
                element.append(`<meta name="${name}" content="${content}">`, { html: true });
            }
        }
    }).transform(resp);
}

export function insertHTML(resp: Response, element: string, content: string) {
    return new HTMLRewriter().on(element, {
        element(element) {
            element.append(content, { html: true });
        }
    }).transform(resp);
}

export function uuidv4() {
    // Generate a RFC4122 version 4 UUID string.
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = crypto.getRandomValues(new Uint8Array(1))[0] & 15;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}



export async function rm(req: Request, env: any, path: string): Promise<Response> {
    // Ensure path starts with a "/" and does not end with one (unless path is just "/")
    if (!path.startsWith("/")) {
        path = "/" + path;
    }
    if (path.length > 1 && path.endsWith("/")) {
        path = path.slice(0, -1);
    }
    const url = new URL(req.url);
    url.host = (isLocal(env) ? env.LOCALHOST : HOST_DOMAIN);
    url.pathname = `/data${path}/`;
    const delReq = new Request(url.toString(), {method: "DELETE"});
    return handleR2FS(delReq, env, "/data");
}

export async function mkdir(req: Request, env: any, path: string, attrs?: Record<string, string>): Promise<Response> {
    // Ensure path starts with a "/" and does not end with one (unless path is just "/")
    if (!path.startsWith("/")) {
        path = "/" + path;
    }
    if (path.length > 1 && path.endsWith("/")) {
        path = path.slice(0, -1);
    }
    const url = new URL(req.url);
    url.host = (isLocal(env) ? env.LOCALHOST : HOST_DOMAIN);
    url.pathname = `/data${path}/`;
    const headers = {
        "Content-Type": "application/x-directory",
        "Change-Timestamp": (Date.now() * 1000).toString(),
    }
    if (attrs) {
        for (const [key, value] of Object.entries(attrs)) {
            headers[`Attribute-${key}`] = value;
        }
    }
    const putReq = new Request(url.toString(), {method: "PUT", headers});
    return handleR2FS(putReq, env, "/data");
}