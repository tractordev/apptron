import { Container, getContainer } from "@cloudflare/containers";
import { validateToken } from "./auth";
import { handle as handleR2FS, getAttrs } from "./r2fs";
import { isLocal, redirectToSignin, insertMeta, insertHTML, uuidv4, putdir } from "./util";
import { ADMIN_USERS, HOST_DOMAIN, PUBLISH_DOMAINS } from "./config";
import { Context, parseContext } from "./context";
import * as projects from "./projects";
import * as publicsite from "./public";

export class Session extends Container {
    defaultPort = 8080;
    sleepAfter = "2h";
}

const CORS_HEADERS = {
    "Access-Control-Allow-Origin": "*",
    "Access-Control-Allow-Methods": "GET, HEAD, PUT, POST, PATCH, DELETE, MOVE, COPY, OPTIONS",
    "Access-Control-Allow-Headers": "Content-Type, Authorization",
    "Vary": "Origin",
};

function applyCORS(resp: Response) {
    const newresp = new Response(resp.body, {
        status: resp.status,
        statusText: resp.statusText,
        headers: new Headers(resp.headers)
    });
    for (const [key, value] of Object.entries(CORS_HEADERS)) {
        newresp.headers.set(key, value);
    }
    return newresp;
}

export default {
    async fetch(req: Request, env: any) {
        const url = new URL(req.url);
        const ctx = parseContext(req, env);

        if (
            ctx.portDomain ||
            url.pathname.startsWith("/x/net") || 
            url.pathname.startsWith("/bundles/")
        ) {
            return getContainer(env.session).fetch(req);
        }

        if (req.method === "OPTIONS") {
            return new Response("", { 
                status: 200,
                headers: {...CORS_HEADERS},
            });
        }

        for (const domain of PUBLISH_DOMAINS) {
            if ((isLocal(env) && (url.pathname.split("/")[1]||"").endsWith("."+domain)) 
                || url.host.endsWith("."+domain)) {
                return publicsite.handle(req, env, ctx);
            }
        }

        if (url.pathname.endsWith(".map")) {
            return new Response("", { status: 200 });
        }

        if (ctx.envDomain && url.pathname.startsWith("/edit/")) {
            const project = await projects.getByUUID(env, ctx.subdomain);
            if (project === null) {
                return new Response("Not Found", { status: 404 });
            }
            if (project["visibility"] !== "public" && project["owner"] !== ctx.userUUID) {
                return new Response("Forbidden", { status: 403 });
            }

            const envReq = new Request(new URL("/_env", req.url).toString(), req);
            const envResp = await env.assets.fetch(envReq);
            const resp = new Response(envResp.body, {
                status: envResp.status,
                statusText: envResp.statusText,
                headers: new Headers(envResp.headers)
            });
            // resp.headers.set("Cross-Origin-Opener-Policy", "same-origin");
            // resp.headers.set("Cross-Origin-Embedder-Policy", "require-corp");
            // resp.headers.set("Cross-Origin-Resource-Policy", "cross-origin");
            
            const contentType = resp.headers.get('content-type');
            if (!contentType || !contentType.includes('text/html')) {
                return resp;
            }

            return insertMeta(resp, {
                "auth-url": env.AUTH_URL,
                "project": escapeJSON(JSON.stringify(project)),
            });
        }

        if (["/dashboard", "/shell"].includes(url.pathname)) {
            if (!await validateToken(env.AUTH_URL, ctx.tokenRaw)) {
                return redirectToSignin(env, url);
            }
        }

        if (url.pathname.startsWith("/data")) {
            const authenticated = await validateToken(env.AUTH_URL, ctx.tokenRaw);

            let dataPath = url.pathname.slice(5);
            if (dataPath.endsWith("/...")) {
                dataPath = dataPath.slice(0, -4);
            }
            // admin data urls
            if (dataPath.includes("/:attr/") || 
                dataPath.startsWith("/etc/") ||
                ["","/etc","/env","/usr"].indexOf(dataPath) !== -1) {
                if (!authenticated) {
                    return new Response("Forbidden", { status: 403 });
                }
                if (!ADMIN_USERS.includes(ctx.tokenJWT?.username)) {
                    return new Response("Forbidden", { status: 403 });
                }
            }
            // user data urls
            if (dataPath.startsWith("/usr/")) {
                if (!authenticated) {
                    return new Response("Forbidden", { status: 403 });
                }
                const parts = dataPath.split("/");
                if (!parts[2] || parts[2] !== ctx.userUUID) {
                    return new Response("Forbidden", { status: 403 });
                }
            }
            // env data urls
            if (dataPath.startsWith("/env/")) {
                const envUUID = dataPath.split("/")[2];
                const project = await projects.getByUUID(env, envUUID);
                if (project === null) {
                    return new Response("Not Found", { status: 404 });
                }
                // not public and not owner
                if (project["visibility"] !== "public" && project["owner"] !== ctx.userUUID) {
                    return new Response("Forbidden", { status: 403 });
                }
                // public, not owner, and not GET or HEAD request
                if (project["visibility"] === "public" && project["owner"] !== ctx.userUUID && ["GET", "HEAD"].indexOf(req.method) === -1) {
                    return new Response("Forbidden", { status: 403 });
                }
            }
            return handleR2FS(req, env, "/data");
        }

        // <username>.apptron.dev/edit/<env-name>
        if (url.pathname.startsWith("/edit/")) {
            const parts = url.pathname.split("/");
            const envName = parts[2];
            const project = await projects.getByName(env, ctx.subdomain, envName);
            if (project === null) {
                return new Response("Not Found", { status: 404 });
            }
            if (project["visibility"] !== "public" && project["owner"] !== ctx.userUUID) {
                return new Response("Not Found", { status: 404 });
                // return new Response("Forbidden", { status: 403 });
            }
            return await envPage(req, env, project, "/edit/"+envName);
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
            
            const usrResp = await putdir(req, env, `/usr/${user["user_id"]}`, {
                "username": user["username"],
            });
            if (!usrResp.ok) {
                return usrResp;
            }
            
            const idxResp = await putdir(req, env, `/etc/index/${user["username"]}`, {
                "uuid": user["user_id"],
            });
            if (!idxResp.ok) {
                return idxResp;
            }

            return new Response(null, { status: 204 });
        }
        
        if (ctx.userDomain && url.pathname.startsWith("/projects")) {
            return applyCORS(await projects.handle(req, env, ctx));
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

        // final fallback: serve 404 page without changing URL
        if (req.method === "GET" && url.pathname !== "/404") {
            const notFoundReq = new Request(new URL("/404", req.url).toString(), req);
            const notFoundResp = await env.assets.fetch(notFoundReq);

            return new Response(notFoundResp.body, {
                status: 404,
                headers: notFoundResp.headers,
            });
        }

        return env.assets.fetch(req);
    },
};

function escapeJSON(json: string) {
    return json.replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
}

function ensureSystemDirs(req: Request, env: any) {
    console.log("Ensuring system directories exist...");
    return Promise.all([
        putdir(req, env, "/"),
        putdir(req, env, "/etc"),
        putdir(req, env, "/etc/index"),
        putdir(req, env, "/usr"),
        putdir(req, env, "/env"),
    ]);
}


async function envPage(req: Request, env: any, project: any, path: string) {
    const url = new URL(req.url);
    if (isLocal(env)) {
        url.searchParams.set("env", project["uuid"]);
        url.host = env.LOCALHOST;
    } else {
        url.host = project["uuid"] + "." + HOST_DOMAIN;
    }
    url.pathname = path;
    const envReq = new Request(new URL("/_frame", req.url).toString(), req);
    const envResp = await env.assets.fetch(envReq);
    const resp = new Response(envResp.body, {
        status: envResp.status,
        statusText: envResp.statusText,
        headers: new Headers(envResp.headers)
    });
    // resp.headers.set("Cross-Origin-Opener-Policy", "same-origin");
    // resp.headers.set("Cross-Origin-Embedder-Policy", "require-corp");
    // resp.headers.set("Cross-Origin-Resource-Policy", "cross-origin");
    return insertHTML(resp, "body", `<iframe src="${url.toString()}" allow="usb; serial; hid; clipboard-read; clipboard-write;"></iframe>`);
}
//sandbox="allow-scripts allow-same-origin allow-forms allow-modals allow-popups allow-popups-to-escape-sandbox"

