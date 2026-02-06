import { PUBLISH_DOMAINS } from "./config";
import { Context } from "./context";
import { isLocal } from "./util";
import * as projects from "./projects";
import mime from 'mime';

export async function handle(req: Request, env: any, ctx: Context) {
    const url = new URL(req.url);
    
    // Check if url.pathname does not end with "/" and does not include a file extension
    if (!url.pathname.endsWith("/") && !/\.[^\/]+$/.test(url.pathname)) {
        url.pathname = url.pathname + "/";
        return Response.redirect(url.toString(), 302);
    }

    let domain = undefined;
    let username = undefined;
    let basePath = "";
    if (isLocal(env)) {
        domain = url.pathname.split("/")[1];
        username = domain.split(".")[0];
        url.pathname = url.pathname.slice(domain.length + 1);
        basePath = `/${domain}`;
    } else {
        domain = url.host;
        username = domain.split(".")[0];
    }
    // this will have to change with custom domains (no basepaths)
    let envName = url.pathname.split("/")[1];
    basePath += `/${envName}`;
    url.pathname = url.pathname.slice(envName.length + 1);

    const project = await projects.getByName(env, username, envName);
    if (!project) {
        return new Response('Not found', { status: 404 });
    }
    
    let publicPath = `/env/${project["uuid"]}/public`;
    let objectKey = `${publicPath}${url.pathname}`;
    let object = await env.bucket.get(objectKey);
    if (!object && url.pathname === "/sw.js") {
        return new Response(await generateSW(env, publicPath), { status: 200, headers: { 'Content-Type': 'application/javascript' } });
    }
    if (!object || object.customMetadata["Content-Type"] === "application/x-directory") {
        objectKey = `${publicPath}${url.pathname}/index.html`.replace(/\/{2,}/g, "/");
        object = await env.bucket.get(objectKey);
        if (!object) {
            object = await env.bucket.get(`${publicPath}/404.html`);
            if (object) {
                return new Response(object.body, {
                    headers: {
                        'Content-Type': object.httpMetadata.contentType || mime.getType(url.pathname) || 'text/html',
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
            'Content-Type': object.httpMetadata.contentType || mime.getType(url.pathname) || 'text/html',
        },
    });
}

async function generateSW(env, publicPath) {
    const listAll = async (bucket, prefix) => {
        const files = [];
        let result;
        let cursor;
        do {
            result = await bucket.list({ prefix, cursor });
            files.push(...result.objects
                .map(obj => obj.key.slice(publicPath.length)));
            let indexes = result.objects
                .filter(o => o.key.endsWith("/index.html"))
                .map(o => o.key.slice(publicPath.length))
                .map(k => k.replace(/\/index\.html$/, "/"));
            files.push(...indexes);
            cursor = result.cursor;
        } while (cursor);
        return files.sort().filter(f => f); // remove empty
    };
    const assets = await listAll(env.bucket, publicPath);
    return `
if (globalThis["ServiceWorkerGlobalScope"] && self instanceof ServiceWorkerGlobalScope) {
    const CACHE_NAME = 'my-site-v1';
    const MAX_AGE = 1000 * 60 * 60; // 1 hour in milliseconds
    const ASSETS = [
${assets.map(a => `dirname(import.meta.url)+'${a}'`).join(",\n")}
    ];

    // Install: cache assets
    self.addEventListener('install', (e) => {
        e.waitUntil(
            caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
        );
    });

    // Fetch: network-first with stale-while-revalidate if fresh enough
    self.addEventListener('fetch', (e) => {
        e.respondWith(
            caches.open(CACHE_NAME).then(async (cache) => {
                const cached = await cache.match(e.request);

                if (cached) {
                    const cachedTime = cached.headers.get('sw-cached-at');
                    const age = Date.now() - parseInt(cachedTime || 0);

                    if (age < MAX_AGE) {
                        // Fresh enough: stale-while-revalidate
                        fetchAndCache(cache, e.request);
                        return cached;
                    } else {
                        // Too old: network-first
                        return fetchAndCache(cache, e.request).catch(() => cached);
                    }
                }

                // No cache: fetch from network
                return fetchAndCache(cache, e.request);
            })
        );
    });

    async function fetchAndCache(cache, request) {
        const response = await fetch(request);

        // Clone and add timestamp header
        const headers = new Headers(response.headers);
        headers.set('sw-cached-at', Date.now().toString());

        const timestampedResponse = new Response(await response.clone().blob(), {
            status: response.status,
            statusText: response.statusText,
            headers: headers
        });

        cache.put(request, timestampedResponse);
        return response;
    }
} else {
    navigator.serviceWorker.register(import.meta.url, {type: "module"});
}

function dirname(path) {
  const i = path.lastIndexOf('/');
  return i === -1 ? '.' : i === 0 ? '/' : path.slice(0, i);
}
`;
}

