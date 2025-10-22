import { Container, getContainer } from "@cloudflare/containers";
import { validateToken } from "./auth";
import { handle as handleR2FS, getAttrs } from "./r2fs";

const HOST_DOMAIN = "apptron.dev";
const ADMIN_USERS = ["progrium"];

export default {
    async fetch(req: Request, env: any) {
        const url = new URL(req.url);
        const ctx = parseContext(req, env);
        const authURL = env.AUTH_URL;

        if (url.pathname.endsWith(".map")) {
            return new Response("", { status: 200 });
        }

        if (ctx.envDomain && url.pathname.startsWith("/edit/")) {
            const parts = url.pathname.split("/");
            const envName = parts[2];

            const envReq = new Request(new URL("/_env", req.url).toString(), req);
            const resp = await env.assets.fetch(envReq);
            
            const contentType = resp.headers.get('content-type');
            if (!contentType || !contentType.includes('text/html')) {
                return resp;
            }

            return insertMeta(resp, {
                "auth-url": authURL,
                "env-name": envName,
            });
        }

        if (["/dashboard", "/shell"].includes(url.pathname)) {
            if (!await validateToken(authURL, ctx.tokenRaw)) {
                return redirectToSignin(env, url);
            }
        }

        if (url.pathname.startsWith("/data")) {
            if (!await validateToken(authURL, ctx.tokenRaw)) {
                return new Response("Forbidden", { status: 403 });
            }
            if (url.pathname.includes("/:attr/")) {
                if (!ADMIN_USERS.includes(ctx.tokenJWT?.username)) {
                    return new Response("Forbidden", { status: 403 });
                }
            }
            // todo: validate user access to data path!
            return handleR2FS(req, env, "/data");
        }

        if (url.pathname.startsWith("/edit/")) {
            const parts = url.pathname.split("/");
            const envName = parts[2];
            const envUUID = await lookupEnvUUID(env, ctx, envName);
            if (envUUID === null) {
                return new Response("Not Found", { status: 404 });
            }
            return await envPage(req, env, envUUID, "/edit/"+envName);
        }

        if (url.pathname === "/" && req.method === "GET") {
            await ensureSystemDirs(req, env);
            return redirectToSignin(env, url);
        }

        if (ctx.userDomain && url.pathname === "/" && req.method === "PUT") {
            // ensure user is set up
            if (!await validateToken(authURL, ctx.tokenRaw)) {
                return new Response("Forbidden", { status: 403 });
            }
            const user = await req.json();
            
            const usrURL = new URL(req.url);
            usrURL.pathname = `/data/usr/${user["user_id"]}/`;
            usrURL.host = (isLocal(env) ? env.LOCALHOST : HOST_DOMAIN);
            const usrReq = new Request(usrURL.toString(), {method: "PUT"});
            const usrResp = await handleR2FS(usrReq, env, "/data");
            if (!usrResp.ok) {
                return usrResp;
            }

            usrURL.pathname = `/data/etc/index/${user["username"]}/`;
            const idxReq = new Request(usrURL.toString(), {
                method: "PUT", 
                headers: {
                    "Content-Type": "application/x-directory",
                    "Attribute-UUID": user["user_id"],
                },
            });
            const idxResp = await handleR2FS(idxReq, env, "/data");
            if (!idxResp.ok) {
                return idxResp;
            }

            return new Response(null, { status: 204 });
        }
        
        if (ctx.userDomain && url.pathname.startsWith("/projects")) {
            if (!await validateToken(authURL, ctx.tokenRaw)) {
                return new Response("Forbidden", { status: 403 });
            }
            let resp;
            switch (req.method) {
            case "GET":
                const projects: Record<string, string>[] = [];
                let cursor: string | undefined = undefined;
                do {
                    const prefix = `/etc/index/${ctx.userName}/`;
                    const page = await env.bucket.list({
                        prefix,
                        include: ["customMetadata"],
                        cursor,
                        limit: 1000,
                    });
                    for (const obj of page.objects || []) {
                        const project = {
                            name: obj.key.slice(prefix.length)
                        };
                        for (const [key, value] of Object.entries(obj.customMetadata)) {
                            if (key.startsWith("Attribute-")) {
                                project[key.slice(10)] = value;
                            }
                        }
                        projects.push(project);
                    }
                    cursor = page.truncated ? page.cursor : undefined;
                } while (cursor);

                return new Response(JSON.stringify(projects), { status: 200 });
            case "POST":
                const project = await req.json();
                if (!project["name"]) {
                    return new Response("Bad Request", { status: 400 });
                }
                let name = project["name"].trim();
                // Remove all characters except alphanumeric, dash, and underscore, replace spaces with dashes
                name = name.replace(/\s+/g, "-").replace(/[^A-Za-z0-9\-_]/g, "");
                project["name"] = name;
                project["uuid"] = uuidv4();

                
                resp = await mkdir(req, env, `/env/${project["uuid"]}`, {
                    "name": project["name"],
                    "owner": ctx.userUUID,
                });
                if (!resp.ok) {
                    return resp;
                }

                resp = await mkdir(req, env, `/env/${project["uuid"]}/project`);
                if (!resp.ok) {
                    await rm(req, env, `/env/${project["uuid"]}`);
                    return resp;
                }

                resp = await mkdir(req, env, `/etc/index/${ctx.userName}/${project["name"]}`, {
                    "uuid": project["uuid"],
                    "description": project["description"] || "",
                    "owner": ctx.userUUID,
                });
                if (!resp.ok) {
                    await rm(req, env, `/env/${project["uuid"]}`);
                    return resp;
                }

                const projectURL = new URL(req.url);
                projectURL.pathname = `/edit/${project["name"]}`;
                return new Response(null, { status: 201, headers: { "Location": projectURL.toString() } });
            case "DELETE":
                if (!url.pathname.startsWith("/projects/")) {
                    return new Response("Not Found", { status: 404 });
                }
                const projectName = url.pathname.split("/").pop() || "";
                if (!projectName) {
                    return new Response("Not Found", { status: 404 });
                }

                const attrs = await getAttrs(env.bucket, `/etc/index/${ctx.userName}/${projectName}`);
                if (!attrs) {
                    return new Response("Not Found", { status: 404 });
                }

                resp = await rm(req, env, `/etc/index/${ctx.userName}/${projectName}`);
                if (!resp.ok) {
                    return resp;
                }

                resp = await rm(req, env, `/env/${attrs["uuid"]}`);
                if (!resp.ok) {
                    return resp;
                }

                return new Response(null, { status: 204 });
            }
        }

        if (url.pathname.startsWith("/x/local")) {
            return new Response("OK", { status: 200 });
        }
        
        if (url.pathname.startsWith("/x/net") || 
            url.host.startsWith("_") ||
            url.pathname === "/bundle.tgz") {
            return getContainer(env.session).fetch(req);
        }

        if (["/signin", "/signout", "/shell", "/dashboard"].includes(url.pathname)) {
            const resp = await env.assets.fetch(req);

            const contentType = resp.headers.get('content-type');
            if (!contentType || !contentType.includes('text/html')) {
                return resp;
            }

            return insertMeta(resp, {
                "auth-url": authURL
            });
        }

        return env.assets.fetch(req);
    },
};

export class Session extends Container {
    defaultPort = 8080;
    sleepAfter = "1h";
}

export interface Context {
    tokenRaw?: string;
    tokenJWT?: Record<string, any>;
    userUUID?: string;
    userName?: string;
    userDomain: boolean;
    envUUID?: string;
    envName?: string;
    envDomain: boolean;
}

function ensureSystemDirs(req: Request, env: any) {
    console.log("Ensuring system directories exist...");
    return Promise.all([
        mkdir(req, env, "/"),
        mkdir(req, env, "/etc"),
        mkdir(req, env, "/etc/index"),
        mkdir(req, env, "/usr"),
        mkdir(req, env, "/env"),
    ]);
}

async function rm(req: Request, env: any, path: string): Promise<Response> {
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

async function mkdir(req: Request, env: any, path: string, attrs?: Record<string, string>): Promise<Response> {
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
    }
    if (attrs) {
        for (const [key, value] of Object.entries(attrs)) {
            headers[`Attribute-${key}`] = value;
        }
    }
    const putReq = new Request(url.toString(), {method: "PUT", headers});
    return handleR2FS(putReq, env, "/data");
}

async function lookupEnvUUID(env: any, ctx: Context, envName: string) {
    if (!envName) {
        return null;
    }
    const envObject = await env.bucket.get(`/etc/index/${ctx.userName}/${envName}`);
    if (envObject === null) {
        return null;
    }
    if (envObject.customMetadata["Attribute-uuid"] === undefined) {
        return null;
    }
    return envObject.customMetadata["Attribute-uuid"];
}

function parseJWT(token: string): Record<string, any> {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    return JSON.parse(atob(base64));
  }

function parseContext(req: Request, env: any): Context {
    const url = new URL(req.url);
    const ctx: Context = {
        userDomain: false,
        envDomain: false,
    };

    ctx.tokenRaw = url.searchParams.get("token") || undefined;
    if (!ctx.tokenRaw) {
        const cookie = req.headers.get("cookie") || "";
        const match = cookie.match(/hanko=([^;]+)/);
        if (match) {
            ctx.tokenRaw = match[1] || undefined;
        }
    }
    if (!ctx.tokenRaw) {
        ctx.tokenRaw = req.headers.get("Authorization")?.split(" ")[1] || undefined;
    }

    if (ctx.tokenRaw) {
        ctx.tokenJWT = parseJWT(ctx.tokenRaw);
        ctx.userUUID = ctx.tokenJWT["sub"]; // should be user_id
    }

    if (url.host.endsWith("." + HOST_DOMAIN)) {
        const subdomain = url.host.slice(0, -("." + HOST_DOMAIN).length);
        if (subdomain.length >= 32) {
            ctx.envDomain = true;
            ctx.envUUID = subdomain;
        } else {
            ctx.userDomain = true;
            ctx.userName = subdomain;
        }
    }

    if (url.searchParams.get("env")) {
        ctx.envUUID = url.searchParams.get("env") || undefined;
        ctx.envDomain = true;
    } else if (url.searchParams.get("user")) {
        ctx.userName = url.searchParams.get("user") || undefined;
        ctx.userDomain = true;
    }

    return ctx;
}

function isLocal(env: any) {
    return !!(env && env.LOCALHOST);
}

function redirectToSignin(env: any, url: URL) {
    if (isLocal(env)) {
        url.host = env.LOCALHOST;
    } else {
        url.host = HOST_DOMAIN;
    }
    url.pathname = "/signin";
    return Response.redirect(url.toString(), 307);
}

function insertMeta(resp: Response, meta: Record<string, string>) {
    return new HTMLRewriter().on('head', {
        element(element) {
            for (const [name, content] of Object.entries(meta)) {
                element.append(`<meta name="${name}" content="${content}">`, { html: true });
            }
        }
    }).transform(resp);
}

function insertHTML(resp: Response, element: string, content: string) {
    return new HTMLRewriter().on(element, {
        element(element) {
            element.append(content, { html: true });
        }
    }).transform(resp);
}

async function envPage(req: Request, env: any, envUUID: string, path: string) {
    const url = new URL(req.url);
    if (isLocal(env)) {
        url.searchParams.set("env", envUUID);
        url.host = env.LOCALHOST;
    } else {
        url.host = envUUID + "." + HOST_DOMAIN;
    }
    url.pathname = path;
    const envReq = new Request(new URL("/_frame", req.url).toString(), req);
    return insertHTML(await env.assets.fetch(envReq), "body", `<iframe src="${url.toString()}" allow="usb; serial; hid; clipboard-read; clipboard-write; cross-origin-isolated"
        sandbox="allow-scripts allow-same-origin allow-forms allow-modals allow-popups allow-popups-to-escape-sandbox"></iframe>`);
}


function uuidv4() {
    // Generate a RFC4122 version 4 UUID string.
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = crypto.getRandomValues(new Uint8Array(1))[0] & 15;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}
