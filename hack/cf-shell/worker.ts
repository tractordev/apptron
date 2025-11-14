import { Container, getContainer } from "@cloudflare/containers";

export class Shell extends Container {
    defaultPort = 8080;
    sleepAfter = "1h";
}

export default {
    async fetch(req: Request, env: any) {
        return getContainer(env.shell).fetch(req);
    }
}