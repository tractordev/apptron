import { PUBLISH_DOMAINS } from "./config";
import { Context } from "./context";
import { isLocal } from "./util";
import * as projects from "./projects";

export async function handle(req: Request, env: any, ctx: Context) {
    const url = new URL(req.url);
    let domain = undefined;
    let username = undefined;
    if (isLocal(env)) {
        domain = url.pathname.split("/")[1];
        username = domain.split(".")[0];
        url.pathname = url.pathname.slice(domain.length + 1);
    } else {
        domain = url.host;
        username = domain.split(".")[0];
    }
    let envName = url.pathname.split("/")[1];
    url.pathname = url.pathname.slice(envName.length + 1);

    const project = await projects.getByName(env, username, envName);
    if (!project) {
        return new Response('Not found', { status: 404 });
    }
      
    let objectKey = `/env/${project["uuid"]}/public${url.pathname}`;
    let object = await env.bucket.get(objectKey);
    if (!object || object.customMetadata["Content-Type"] === "application/x-directory") {
        objectKey = `/env/${project["uuid"]}/public${url.pathname}/index.html`.replace(/\/{2,}/g, "/");
        object = await env.bucket.get(objectKey);
        if (!object) {
            object = await env.bucket.get(`/env/${project["uuid"]}/public/404.html`);
            if (object) {
                return new Response(object.body, {
                    headers: {
                        'Content-Type': object.httpMetadata.contentType || 'text/html',
                    },
                    status: 404,
                });
            } else {
                return new Response('Not found', { status: 404 });
            }
        }
    }

    return new Response(object.body, {
        headers: {
            'Content-Type': object.httpMetadata.contentType || 'text/html',
        },
    });
}