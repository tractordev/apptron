import { WanixRuntime } from "/wanix.min.js";
import { register } from "/hanko/elements.js";

// querySelector conveniences, but dont import these in html components
// because they have their own that work differently for shadowRoot
export function $(selector) { return document.querySelector(selector); }
export function $$(selector) { return document.querySelectorAll(selector); }

export async function setupWanix() {
    const params = new URLSearchParams(window.location.search);
    if (params.get("cache") === "clear" || (isLocalhost() && !params.get("cache"))) {
        await clearAllCache("assets");
        await clearAllCache("bundles");
    }
    const w = new WanixRuntime({
        helpers: true,
        debug9p: params.get('debug9p') === "true",
        wasm: null,
        network: params.get('network') || `${isLocalhost() ? "ws" : "wss"}://${appHost()}/x/net`
    });
    // getting the bundle ourself, and the function to get other bundles
    w._bundle = getBundle("/bundles/sys.tar.gz");
    w._getBundle = getBundle;
    // getting then loading the wasm ourselves
    getCachedOrFetch("/wanix.wasm").then(wasm => w._loadWasm(wasm));
    return w;
}

let auth = null;
export async function getAuth() {
    if (auth) {
        return auth;
    }
    if (!getMeta("auth-url")) {
        throw new Error("auth-url meta tag not found");
    }
    const { hanko } = await register(getMeta("auth-url"), isLocalhost() ? undefined : {
        cookieDomain: "." + appHost()
    });
    auth = hanko;
    auth.validatedSession = auth.validateSession();
    auth.validatedSession.then(session => {
        if (session.is_valid) {
            console.log("valid session for user", session.claims.username);
        }
    });
    return auth;
}

export function getMeta(name) {
    const meta = document.querySelector('meta[name="' + name + '"]');
    if (!meta) {
        return null;
    }
    return meta.content;
}

export function isLocalhost() {
    const hostname = window.location.hostname;
    return hostname === "localhost" || hostname === "127.0.0.1" || hostname === "::1";
}

export function isUserDomain() {
    const params = new URLSearchParams(window.location.search);
    if (params.get("user")) {
        return true;
    }
    const subdomain = window.location.hostname.split(".").slice(0, -2).join(".");
    if (!subdomain) {
        return false;
    }
    if (subdomain.length >= 32) {
        // env domain
        return false;
    }
    return true;
}

export function isEnvDomain() {
    const params = new URLSearchParams(window.location.search);
    if (params.get("env")) {
        return true;
    }
    const subdomain = window.location.hostname.split(".").slice(0, -2).join(".");
    if (!subdomain) {
        return false;
    }
    if (subdomain.length < 32) {
        // user domain
        return false;
    }
    return true;
}

export function envUUID() {
    if (!isEnvDomain()) {
        return null;
    }
    const params = new URLSearchParams(window.location.search);
    if (params.get("env")) {
        return params.get("env");
    }
    const subdomain = window.location.hostname.split(".").slice(0, -2).join(".");
    if (subdomain.length < 32) {
        return null;
    }
    return subdomain;
}

export async function envUsername() {
    const params = new URLSearchParams(window.location.search);
    if (params.get("user")) {
        return params.get("user");
    }
    const hostname = new URL(await currentURL()).hostname;
    return hostname.split(".").slice(0, -2).join(".");
}

export function appHost() {
    const hostname = window.location.origin.replace("https://", "").replace("http://", "");
    if (isLocalhost()) {
        return hostname;
    }
    const parts = hostname.split(".");
    if (parts.length >= 2) {
        return parts.slice(-2).join(".");
    }
    return hostname;
}

export function urlFor(path, params = {}, user = null) {
    let host = appHost();
    if (user && !isLocalhost()) {
        host = user + "." + host;
    }
    const currentURL = new URL(window.location.href);
    const url = new URL(currentURL.protocol + "//" + host + path);
    if (params && Object.keys(params).length > 0) {
        for (const [key, value] of Object.entries(params)) {
            url.searchParams.set(key, value);
        }
    }
    if (user && isLocalhost()) {
        url.searchParams.set("user", user);
    }
    return url.toString()
}

export function currentURL() {
    if (isEnvDomain()) {
        const reply = new MessageChannel();
        top.postMessage({ self: true, reply: reply.port2 }, getOrigin(), [reply.port2]);
        return new Promise((resolve, reject) => {
            reply.port1.onmessage = (e) => resolve(e.data);
        });
    }
    return Promise.resolve(window.location.href);
}

export function getOrigin() {
    let origin = window.location.protocol + "//" + appHost();
    if (window.apptron) {
        origin = window.location.protocol + "//" + window.apptron.user.username + "." + appHost();
    }
    if (isLocalhost()) {
        origin = "*";
    }
    return origin;
}

export function redirectTo(url) {
    if (isEnvDomain()) {
        top.postMessage({ redirect: url }, getOrigin());
        return;
    }
    window.location.href = url;
}

export async function authRedirect(defaultTarget = "/", user = null) {
    const currentParams = new URLSearchParams(window.location.search);
    const redirect = currentParams.get("redirect") || defaultTarget;
    redirectTo(urlFor(redirect, {}, user));
}

export function secondsSince(timestamp) {
    const then = new Date(timestamp);
    const now = new Date();
    const diffInMs = now - then;
    return Math.floor(diffInMs / 1000);
}

export async function getCachedOrFetch(url, gzipped = false, cacheName = "assets") {
    try {
        // Open the cache
        const cache = await caches.open(cacheName);

        // Check if the asset is already cached
        const cachedResponse = await cache.match(url);

        if (cachedResponse) {
            console.log('Found in cache:', url);
            if (gzipped) {
                if (!("DecompressionStream" in window)) {
                    throw new Error("DecompressionStream not supported in this browser.");
                }
                // Decompress stream and return as ArrayBuffer
                const decompressed = cachedResponse.body
                    .pipeThrough(new DecompressionStream("gzip"));
                const decompressedBuffer = await new Response(decompressed).arrayBuffer();
                return decompressedBuffer;
            } else {
                // Return ArrayBuffer from cached response
                return await cachedResponse.arrayBuffer();
            }
        }

        // Not in cache, fetch from network
        console.log('Not in cache, fetching:', url);
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Clone the response since we need to use it twice
        // (once for cache, once for returning the ArrayBuffer)
        const responseToCache = response.clone();

        // Store in cache for future use
        await cache.put(url, responseToCache);
        console.log('Stored in cache:', url);

        if (gzipped) {
            if (!("DecompressionStream" in window)) {
                throw new Error("DecompressionStream not supported in this browser.");
            }
            // Decompress stream and return as ArrayBuffer
            const decompressed = response.body
                .pipeThrough(new DecompressionStream("gzip"));
            const decompressedBuffer = await new Response(decompressed).arrayBuffer();
            return decompressedBuffer;
        } else {
            // Return ArrayBuffer from the original response
            return await response.arrayBuffer();
        }

    } catch (error) {
        console.error('Error in getCachedOrFetch:', error);
        throw error;
    }
}

export async function clearAllCache(cacheName = "assets") {
    const deleted = await caches.delete(cacheName);
    console.log('Deleted entire cache:', cacheName, deleted);
    return deleted;
}

export function modalDialog(el) {
    el.querySelectorAll('[data-action="close"]').forEach(closer => {
        closer.addEventListener("click", () => el.close());
    });
    el.addEventListener("click", (e) => {
        const r = el.getBoundingClientRect();
        const inBounds = (
            e.clientX >= r.left &&
            e.clientX <= r.right &&
            e.clientY >= r.top &&
            e.clientY <= r.bottom
        );
        if (!inBounds) el.close();
    });
    return el;
}

let cacheFrame = null;
export async function getBundle(name) {
    if (!document) {
        return null;
    }
    
    if (!cacheFrame) {
        cacheFrame = new Promise(resolve => {
            const el = document.createElement("iframe");
            // todo: use apptron.dev/bundles
            el.src = "/bundles";
            el.style.display = "none";
            el.onload = () => {
                resolve(el);
            };
            document.body.appendChild(el);
        });
    }

    const el = await cacheFrame;
    const channel = new MessageChannel();
    el.contentWindow.postMessage({ type: "bundle", name: name, port: channel.port2 }, "*", [channel.port2]);
    return await new Promise((resolve, reject) => {
        channel.port1.onmessage = (e) => resolve(e.data.bundle);
    });
}

export async function copyText(text) {
    try {
        await navigator.clipboard.writeText(text);
        return true;
    } catch {
        try {
            const ta = document.createElement("textarea");
            ta.value = text;
            ta.setAttribute("readonly", "");
            ta.style.position = "absolute";
            ta.style.left = "-9999px";
            document.body.appendChild(ta);
            ta.select();
            const ok = document.execCommand("copy");
            document.body.removeChild(ta);
            return ok;
        } catch {
            return false;
        }
    }
}