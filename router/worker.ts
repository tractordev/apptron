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
        if (url.pathname.startsWith("/x/net")) {
            return getContainer(env.session).fetch(req);
        }
        if (url.pathname === "/workbench.json") {
            const resp = await env.assets.fetch(req);
            const workbench = await resp.json();
            workbench.additionalBuiltinExtensions = [{
                scheme: url.protocol.replace(":", ""),
                authority: url.host,
                path: "/system"
            }];
            workbench.profile = {
                name: "Default",
                contents: JSON.stringify({
                    "globalState": JSON.stringify({
                        "storage": {
                            "workbench.activity.pinnedViewlets2": JSON.stringify([
                                { "id": "workbench.view.explorer", "pinned": true, "visible": true, "order": 0 },
                                { "id": "workbench.view.search", "pinned": false, "visible": false, "order": 1 },
                                { "id": "workbench.view.scm", "pinned": false, "visible": false, "order": 2 },
                                { "id": "workbench.view.debug", "pinned": true, "visible": true, "order": 3 },
                                { "id": "workbench.view.extensions", "pinned": false, "visible": false, "order": 4 }
                            ])
                        }
                    })
                })
            };
            return new Response(JSON.stringify(workbench), {
                headers: {
                    "content-type": "application/json",
                },
            });
        }
        return env.assets.fetch(req);
    },
};