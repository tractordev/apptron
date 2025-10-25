import { WanixRuntime } from "/wanix.js";
import { register } from "/hanko/elements.js";

export function setupWanix() {
    const params = new URLSearchParams(window.location.search);
    const w = new WanixRuntime({ 
        helpers: true, 
        bundle: params.get('bundle') || urlFor("/bundle.tgz"),
        wasm: params.get('wasm') || urlFor("/wanix.wasm"),
        network: params.get('network') || "wss://apptron.dev/x/net"
    });
    return w;
}

let auth = null;
export async function getAuth() {
    if (auth) {
        return auth;
    }
    if (isLocalhost()) {
        const { hanko } = await register(getMeta("auth-url"));
        auth = hanko;
        return auth;
    }
    const { hanko } = await register(getMeta("auth-url"), {
        cookieDomain: "."+appHost()
    });
    auth = hanko;
    return auth;
}

export function getMeta(name) {
    const meta = document.querySelector('meta[name="'+name+'"]');
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

export function redirectTo(url) {
    if (isEnvDomain()) {
        let origin = window.location.protocol + "//" + appHost();
        if (window.apptron) {
            origin = window.location.protocol + "//" + window.apptron.user.username + "." + appHost();
        }
        if (isLocalhost()) {
            origin = "*";
        }
        top.postMessage({ redirect: url }, origin);
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