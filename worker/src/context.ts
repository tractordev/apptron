import { HOST_DOMAIN } from "./config";
import { parseJWT } from "./auth";

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

export function parseContext(req: Request, env: any): Context {
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
