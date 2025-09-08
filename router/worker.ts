import { Container, getContainer } from "@cloudflare/containers";

export class Session extends Container {
  defaultPort = 8080;
  sleepAfter = "1h";
}

export default {
    async fetch(req, env) {  
        const url = new URL(req.url);
        if (url.pathname === "/") {
            return Response.redirect("https://github.com/tractordev/apptron", 307);
        }
        if (url.pathname.startsWith("/x/local")) {
            return new Response("OK", { status: 200 });
        }
        if (url.pathname.startsWith("/x/net") || 
            url.host.startsWith("_") ||
            url.pathname === "/bundle.tgz") {
            return getContainer(env.session).fetch(req);
        }
        return env.assets.fetch(req);
    },
};