import { HOST_DOMAIN } from "./config";
import { parseJWT } from "./auth";

export interface Context {
    tokenRaw?: string;
    tokenJWT?: Record<string, any>;
    userUUID?: string;
    username?: string;

    userDomain: boolean;
    envDomain: boolean;
    portDomain: boolean;
    subdomain?: string; // username or env UUID
}

export function parseContext(req: Request, env: any): Context {
    const url = new URL(req.url);
    const ctx: Context = {
        userDomain: false,
        envDomain: false,
        portDomain: false,
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
        ctx.username = ctx.tokenJWT["username"];
    }

    if (url.host.endsWith("." + HOST_DOMAIN)) {
        const subdomain = url.host.slice(0, -("." + HOST_DOMAIN).length);
        if (subdomain.startsWith("tcp-") && subdomain.split("-").length === 4) {
            ctx.portDomain = true;
        } else if (subdomain.length >= 32) {
            ctx.envDomain = true;
        } else {
            ctx.userDomain = true;
        }
        ctx.subdomain = subdomain;
    }

    if (url.searchParams.get("env")) {
        ctx.subdomain = url.searchParams.get("env") || undefined;
        ctx.envDomain = true;
    } else if (url.searchParams.get("user")) {
        ctx.subdomain = url.searchParams.get("user") || undefined;
        ctx.userDomain = true;
    } else if (url.searchParams.get("port")) {
        ctx.portDomain = true;
        ctx.subdomain = url.searchParams.get("port") || undefined;
    }

    return ctx;
}
