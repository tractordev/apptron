/**
 * Web component that loads and mounts the VS Code workbench.
 * Call mount() when the host is ready (e.g. after any other scripts that must run first).
 *
 * Cascade: workbench CSS is loaded before app styles. App styles are in @layer app so
 * unlayered workbench rules (e.g. .monaco-tl-twistie padding) are not overridden.
 */
const DEFAULT_PROFILE = {
  name: "Default",
  contents: JSON.stringify({
    globalState: JSON.stringify({
      storage: {
        "workbench.explorer.views.state.hidden":
          '[{"id":"outline","isHidden":true},{"id":"timeline","isHidden":true},{"id":"workbench.explorer.openEditorsView","isHidden":true},{"id":"workbench.explorer.emptyView","isHidden":false},{"id":"npm","isHidden":true}]',
        "workbench.panel.pinnedPanels":
          '[{"id":"workbench.panel.markers","pinned":false,"visible":false,"order":0},{"id":"workbench.panel.output","pinned":false,"visible":false,"order":1},{"id":"workbench.panel.repl","pinned":true,"visible":false,"order":2},{"id":"terminal","pinned":true,"visible":false,"order":3},{"id":"workbench.panel.testResults","pinned":true,"visible":false,"order":3},{"id":"refactorPreview","pinned":true,"visible":false}]',
        "workbench.activity.pinnedViewlets2": JSON.stringify([
          { id: "workbench.view.explorer", pinned: true, visible: true, order: 0 },
          { id: "workbench.view.search", pinned: true, visible: true, order: 1 },
          { id: "workbench.view.scm", pinned: false, visible: false, order: 2 },
          { id: "workbench.view.debug", pinned: false, visible: false, order: 3 },
          { id: "workbench.view.extensions", pinned: false, visible: false, order: 4 },
        ]),
      },
    }),
  }),
};

export class VSCodeWorkbenchComponent extends HTMLElement {
  constructor() {
    super();
    this._mounted = false;
  }

  connectedCallback() {
    this.style.flex = "1";
    this.style.minHeight = "0";
    this.style.display = "flex";
    this.style.flexDirection = "column";
  }

  /** Load and mount the workbench. Idempotent. */
  mount() {
    if (this._mounted) return;
    this._mounted = true;

    const nls = document.createElement("script");
    nls.src = "/vscode/out/nls.messages.js";
    const loader = document.createElement("script");
    loader.src = "/vscode/out/vs/loader.js";

    const runBootstrap = () => {
      const go = () => {
        if (typeof require === "undefined") {
          setTimeout(go, 0);
          return;
        }
        require.config({ baseUrl: window.location.origin + "/vscode/out" });
        this._createWorkbench();
      };
      go();
    };

    loader.onload = runBootstrap;
    loader.onerror = () => this.dispatchEvent(new CustomEvent("error", { detail: new Error("Failed to load VS Code loader") }));
    document.head.appendChild(loader);
    document.head.appendChild(nls);
  }

  _createWorkbench() {
    const ch = new MessageChannel();
    ch.port2.onmessage = (event) => {
      if (event.data.port) window.postMessage(event.data, "*", [event.data.port]);
    };

    const url = new URL(window.location.href);
    const scheme = url.protocol.replace(":", "");
    const hostParts = url.host.split(".");
    if (hostParts.length > 2) hostParts.shift();
    const hostJoin = hostParts.join(".");

    const config = {
      messagePorts: new Map([["progrium.apptron-system", ch.port1]]),
      productConfiguration: {
        extensionEnabledApiProposals: { "progrium.apptron-system": ["ipc"] },
        webviewContentExternalBaseUrlTemplate:
          `${scheme}://{{uuid}}.${hostJoin}/vscode/out/vs/workbench/contrib/webview/browser/pre/`,
      },
      configurationDefaults: {
        "workbench.colorTheme": "Tractor Dark",
        "workbench.secondarySideBar.defaultVisibility": "hidden",
        "workbench.statusBar.visible": false,
        "workbench.layoutControl.enabled": false,
        "window.commandCenter": false,
        "workbench.startupEditor": "none",
        "workbench.activityBar.location": "hidden",
        "workbench.tips.enabled": false,
        "workbench.welcomePage.walkthroughs.openOnInstall": false,
        "problems.visibility": false,
        "editor.minimap.enabled": false,
        "terminal.integrated.tabs.showActions": false,
      },
      developmentOptions: { logLevel: 0 },
      additionalBuiltinExtensions: [
        { scheme, authority: url.host, path: "/system" },
        { scheme, authority: url.host, path: "/preview" },
      ],
      profile: DEFAULT_PROFILE,
      folderUri: { scheme: "wanix", path: "/project" },
    };

    require(["vs/workbench/workbench.web.main"], (wb) => {
      wb.create(this, {
        ...config,
        workspaceProvider: {
          trusted: true,
          workspace: { folderUri: wb.URI.revive(config.folderUri) },
          open(workspace, options) {
            console.log("openFolder requested", workspace);
            return Promise.resolve(true);
          },
        },
      });
      this.dispatchEvent(new CustomEvent("vscode-ready"));
    });
  }
}

customElements.define("vscode-workbench", VSCodeWorkbenchComponent);
