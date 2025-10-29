import { Container, getContainer } from "@cloudflare/containers";
import { validateToken } from "./auth";
import { handle as handleR2FS, getAttrs } from "./r2fs";
import { isLocal, redirectToSignin, insertMeta, insertHTML, uuidv4, mkdir } from "./util";
import { ADMIN_USERS, HOST_DOMAIN } from "./config";
import { Context, parseContext } from "./context";
import * as projects from "./projects";

export class Session extends Container {
    defaultPort = 8080;
    sleepAfter = "1h";
}

export default {
    async fetch(req: Request, env: any) {
        const url = new URL(req.url);
        const ctx = parseContext(req, env);

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
                "auth-url": env.AUTH_URL,
                "env-name": envName,
            });
        }

        if (["/dashboard", "/shell"].includes(url.pathname)) {
            if (!await validateToken(env.AUTH_URL, ctx.tokenRaw)) {
                return redirectToSignin(env, url);
            }
        }

        if (url.pathname.startsWith("/data")) {
            if (!await validateToken(env.AUTH_URL, ctx.tokenRaw)) {
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
            if (!await validateToken(env.AUTH_URL, ctx.tokenRaw)) {
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
            return projects.handle(req, env, ctx);
        }

        if (url.pathname.startsWith("/x/local")) {
            return new Response("OK", { status: 200 });
        }
        
        if (url.pathname.startsWith("/x/net") || 
            url.host.startsWith("_") ||
            url.pathname === "/bundle.tgz") {
            return getContainer(env.session).fetch(req);
        }

        if (["/signin", "/signout", "/shell", "/dashboard", "/debug"].includes(url.pathname)) {
            const resp = await env.assets.fetch(req);

            const contentType = resp.headers.get('content-type');
            if (!contentType || !contentType.includes('text/html')) {
                return resp;
            }

            return insertMeta(resp, {
                "auth-url": env.AUTH_URL
            });
        }

        return env.assets.fetch(req);
    },
};

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
