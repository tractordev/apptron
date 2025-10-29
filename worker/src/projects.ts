import { Context } from "./context";
import { validateToken } from "./auth";
import { uuidv4, putdir, deletepath } from "./util";
import { getAttrs } from "./r2fs";

export async function handle(req: Request, env: any, ctx: Context) {

    if (!await validateToken(env.AUTH_URL, ctx.tokenRaw)) {
        return new Response("Forbidden", { status: 403 });
    }

    switch (req.method) {
    case "GET":
        return handleGet(req, env, ctx);
    case "POST":
        return handlePost(req, env, ctx);
    case "PUT":
        return handlePut(req, env, ctx);
    case "DELETE":
        return handleDelete(req, env, ctx);
    default:
        return new Response("Method Not Allowed", { status: 405 });
    }
}

export async function handleGet(req: Request, env: any, ctx: Context) {
    const projects = await list(env, ctx.username);
    return new Response(JSON.stringify(projects), { status: 200 });
}

export async function handlePost(req: Request, env: any, ctx: Context) {
    const project = await req.json();
    if (!project["name"]) {
        return new Response("Bad Request", { status: 400 });
    }
    let name = project["name"].trim();
    // Remove all characters except alphanumeric, dash, and underscore, replace spaces with dashes
    name = name.replace(/\s+/g, "-").replace(/[^A-Za-z0-9\-_]/g, "");
    project["name"] = name;
    project["uuid"] = uuidv4();

    let resp;
    resp = await putdir(req, env, `/env/${project["uuid"]}`, {
        "name": project["name"],
        "owner": ctx.userUUID,
    });
    if (!resp.ok) {
        return resp;
    }

    resp = await putdir(req, env, `/env/${project["uuid"]}/project`);
    if (!resp.ok) {
        await deletepath(req, env, `/env/${project["uuid"]}`);
        return resp;
    }

    resp = await putdir(req, env, `/etc/index/${ctx.username}/${project["name"]}`, {
        "uuid": project["uuid"],
        "name": project["name"],
        "description": project["description"] || "",
        "owner": ctx.userUUID,
        "visibility": project["visibility"] || "private",
    });
    if (!resp.ok) {
        await deletepath(req, env, `/env/${project["uuid"]}`);
        return resp;
    }

    const projectURL = new URL(req.url);
    projectURL.pathname = `/edit/${project["name"]}`;
    return new Response(null, { status: 201, headers: { "Location": projectURL.toString() } });
}

export async function handlePut(req: Request, env: any, ctx: Context) {
    const url = new URL(req.url);

    if (!url.pathname.startsWith("/projects/")) {
        return new Response("Not Found", { status: 404 });
    }

    const projectName = url.pathname.split("/").pop() || "";
    if (!projectName) {
        return new Response("Bad Request", { status: 400 });
    }

    const update = await req.json();

    // Look up existing project metadata
    const attrs = await getAttrs(env.bucket, `/etc/index/${ctx.username}/${projectName}`);
    if (!attrs) {
        return new Response("Not Found", { status: 404 });
    }

    // Update description (and other metadata if you like)
    const newAttrs = {
        "uuid": attrs["uuid"],
        "owner": ctx.userUUID,
        "name": projectName,
        "description": update["description"] || attrs["description"] || "",
        "visibility": update["visibility"] || attrs["visibility"] || "private",
    };

    // Write updated attributes back using mkdir (idempotent PUT)
    const updateResp = await putdir(req, env, `/etc/index/${ctx.username}/${projectName}`, newAttrs);

    if (!updateResp.ok) {
        return updateResp;
    }

    return new Response(JSON.stringify({ name: projectName, description: newAttrs.description }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
    });
}

export async function handleDelete(req: Request, env: any, ctx: Context) {
    const url = new URL(req.url);
    let resp;

    if (!url.pathname.startsWith("/projects/")) {
        return new Response("Not Found", { status: 404 });
    }
    const projectName = url.pathname.split("/").pop() || "";
    if (!projectName) {
        return new Response("Not Found", { status: 404 });
    }

    const attrs = await getAttrs(env.bucket, `/etc/index/${ctx.username}/${projectName}`);
    if (!attrs) {
        return new Response("Not Found", { status: 404 });
    }

    resp = await deletepath(req, env, `/etc/index/${ctx.username}/${projectName}`);
    if (!resp.ok) {
        return resp;
    }

    resp = await deletepath(req, env, `/env/${attrs["uuid"]}`);
    if (!resp.ok) {
        return resp;
    }

    return new Response(null, { status: 204 });
}

export async function list(env: any, username: string): Promise<Record<string, string>[]> {
    const projects: Record<string, string>[] = [];
    let cursor: string | undefined = undefined;
    do {
        const prefix = `/etc/index/${username}/`;
        const page = await env.bucket.list({
            prefix,
            include: ["customMetadata"],
            cursor,
            limit: 1000,
        });
        for (const obj of page.objects || []) {
            const project = {
                name: obj.key.slice(prefix.length),
                visibility: "private", // default
            };
            for (const [key, value] of Object.entries(obj.customMetadata)) {
                if (key.startsWith("Attribute-")) {
                    project[key.slice(10)] = value;
                }
            }
            projects.push(project);
        }
        cursor = page.truncated ? page.cursor : undefined;
    } while (cursor);
    return projects;
}

export async function getByName(env: any, username: string, projectName: string): Promise<Record<string, string> | null> {
    const attrs = await getAttrs(env.bucket, `/etc/index/${username}/${projectName}`);
    if (!attrs) {
        return null;
    }
    return attrs;
}

export async function getByUUID(env: any, uuid: string): Promise<Record<string, string> | null> {
    const envAttrs = await getAttrs(env.bucket, `/env/${uuid}`);
    if (!envAttrs) {
        return null;
    }
    const userAttrs = await getAttrs(env.bucket, `/usr/${envAttrs["owner"]}`);
    if (!userAttrs) {
        return null;
    }
    const project = await getByName(env, userAttrs["username"], envAttrs["name"]);
    if (!project) {
        return null;
    }
    if (!project["name"]) {
        project["name"] = envAttrs["name"];
    }
    return project;
}