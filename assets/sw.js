/**
 * This file provides a generic message-based service worker infrastructure for custom HTTP-like request/response
 * handling between the main thread (web clients) and a service worker.
 * 
 * STILL IN DEVELOPMENT, DO NOT USE IN PRODUCTION
 * 
 * Key Features:
 * --------------
 * - Instantly claims clients and skips waiting (`install` and `activate` logic).
 * - Maintains a per-client registration so the SW can coordinate with responder endpoints in each client.
 * - For requests with paths beginning with `/:/`, the SW:
 *   - Finds the currently "registered" responder for the client.
 *   - Forwards the fetch metadata (method, URL, headers) to the registered responder via a MessageChannel.
 *   - Waits for a reply or timeout, then returns the reply (with headers/body/status) as a `Response` object.
 *   - On error or timeout returns a network error response.
 * - The `register(handler)` function is designed to be called on the client side:
 *   - It sets up a local responder that can handle SW requests by running a user-defined `handler(Request)`.
 *   - It registers the responder with the SW via `postMessage`, and handles communication using MessageChannel ports.
 *   - Ensures pages can drive fetch logic from the "application" side with a handler for `/ : /` requests.
 * 
 * Intended Usage:
 * ---------------
 * 1. On the main thread, call `register(handler)` to provide a custom request handler for SW requests.
 * 2. Any fetch from the page with a path like `/:/something` will go through the SW, which relays the request
 *    via messages to the registered handler, receives the result, and responds to the fetch.
 * 3. Supports advanced scenarios like mocking server APIs, on-the-fly content, offline logic, etc.
 * 
 * Notes:
 * ------
 * - This system is not for general HTTP request interception, only requests matching specific path prefixes.
 * - All communication uses MessageChannel ports for structured, race-condition-safe messaging.
 * - Timeout, error, and registration logic are handled to ensure robust messaging regardless of page state.
 * - The code is ES module-compatible.
 * 
 * Example (client page):
 * ----------------------
 *   import { register } from '/sw.js';
 *   await register(async (req) => {
 *     // handle Request, return a Response
 *     return new Response("Hello!");
 *   });
 * 
 *   const resp = await fetch('/:/' + ...); // routed by service worker to your handler!
 */

if (globalThis["ServiceWorkerGlobalScope"] && self instanceof ServiceWorkerGlobalScope) {
    const registered = new Map();

    async function cleanupDeadClients() {
        const clientsToDelete = [];
        for (const clientId of registered.keys()) {
            const client = await clients.get(clientId);
            if (!client) clientsToDelete.push(clientId);
        }
        for (const clientId of clientsToDelete) {
            registered.delete(clientId);
        }
    }

    self.addEventListener("install", () => self.skipWaiting()); // Activate immediately, don't wait
    self.addEventListener("activate", event => event.waitUntil(clients.claim())); // Take control of all pages immediately

    self.addEventListener("message", (event) => {
        if (event.data.responder) {
            registered.set(event.source.id, {clientId: event.source.id, ...event.data});
            event.data.ready.postMessage(true);
        }
    });

    self.addEventListener("fetch", async (event) => {
        // find the registration for the fetching client
        let registration = registered.get(event.clientId);
        if (!registration) {
            // no registration found, find the most recent one
            let last = null;
            for (const reg of registered.values()) {
                last = reg;
            }
            if (!last) {
                return;
            }
            registration = last;
        }

        const { timeout = 1000, prefix = "/" } = registration.options;

        const req = event.request;
        const url = new URL(req.url);
        if (!url.pathname.startsWith(prefix)) return;

        const headers = {}
        for (var p of req.headers) {
            headers[p[0]] = p[1]
        }

        event.respondWith(new Promise(async (resolve) => {
            await cleanupDeadClients(); // no awaits before respondWith

            const ch = new MessageChannel();
            const response = new Promise(r => ch.port1.onmessage = e => r(e.data));
            registration.responder.postMessage({
                request: {
                    method: req.method, 
                    url: req.url, 
                    headers: headers,
                },
                responder: ch.port2
            }, [ch.port2]);
            try {
                const reply = await Promise.race([response, new Promise((_, reject) => {
                    setTimeout(() => reject(new Error('Timeout')), timeout);
                })]);
                if (reply.error) {
                    console.warn(reply.error);
                    resolve(Response.error());
                    return;
                }
                resolve(new Response(reply.body, reply));
            } catch (error) {
                console.error(error);
                resolve(Response.error());
            }
        }))
    });

} 

export async function register(handler, options = {}) {
    const responder = new MessageChannel();
    const ready = new MessageChannel();
    
    if (!handler) {
        handler = () => new Response("No handler yet", { status: 503 });
    }

    responder.port1.onmessage = async (event) => {
        const req = new Request(event.data.request.url, {
            method: event.data.request.method,
            headers: event.data.request.headers,
            body: event.data.request.body,
        });
        const resp = await handler(req);
        event.data.responder.postMessage({
            body: await resp.bytes(),
            headers: Object.fromEntries(resp.headers.entries()),
            status: resp.status,
            statusText: resp.statusText,
        });
    };

    await navigator.serviceWorker.register(import.meta.url, {type: "module", scope: options.scope || "/"});
    const registration = await navigator.serviceWorker.ready;
    registration.active.postMessage({
        responder: responder.port2, 
        ready: ready.port2,
        options: options
    }, [responder.port2, ready.port2]);
    
    await new Promise(resolve => ready.port1.onmessage = resolve);

    return (h) => handler = h;
}
