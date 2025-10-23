
import * as vscode from 'vscode';

declare const navigator: unknown;

class BrowserPanel {
	public static currentPanel: BrowserPanel | undefined;
	public static readonly viewType = 'preview.browser';

	private readonly _panel: vscode.WebviewPanel;
	private readonly _extensionUri: vscode.Uri;
	private readonly _extensionContext: vscode.ExtensionContext;
	private _disposables: vscode.Disposable[] = [];
	private _currentUrl: string = 'https://www.google.com';
	private _history: string[] = ['https://www.google.com'];
	private _historyIndex: number = 0;
	private _zoomLevel: number = 1.0;

	public static createOrShow(extensionUri: vscode.Uri, extensionContext: vscode.ExtensionContext, url?: string) {
		const column = vscode.window.activeTextEditor
			? vscode.window.activeTextEditor.viewColumn
			: undefined;

		// If we already have a panel, show it.
		if (BrowserPanel.currentPanel) {
			BrowserPanel.currentPanel._panel.reveal(column);
			if (url) {
				BrowserPanel.currentPanel.navigate(url);
			}
			return;
		}

		// Otherwise, create a new panel.
		const panel = vscode.window.createWebviewPanel(
			BrowserPanel.viewType,
			'Browser Preview',
			column || vscode.ViewColumn.One,
			{
				enableScripts: true,
				retainContextWhenHidden: true,
				localResourceRoots: [extensionUri]
			}
		);

		BrowserPanel.currentPanel = new BrowserPanel(panel, extensionUri, extensionContext, url);
	}


	private constructor(panel: vscode.WebviewPanel, extensionUri: vscode.Uri, extensionContext: vscode.ExtensionContext, url?: string) {
		this._panel = panel;
		this._extensionUri = extensionUri;
		this._extensionContext = extensionContext;

		if (url) {
			this._currentUrl = url;
			this._history = [url];
		}

		// Set webview options for restored panels
		this._panel.webview.options = {
			enableScripts: true,
			localResourceRoots: [this._extensionUri]
		};

		// Set the webview's initial html content
		this._update();

		// Listen for when the panel is disposed
		// This happens when the user closes the panel or when the panel is closed programmatically
		this._panel.onDidDispose(() => this.dispose(), null, this._disposables);


		// Handle messages from the webview
		this._panel.webview.onDidReceiveMessage(
			(message: any) => {
				switch (message.type) {
					case 'navigate':
						this.navigate(message.url);
						break;
					case 'back':
						this.goBack();
						break;
					case 'forward':
						this.goForward();
						break;
					case 'reload':
						this.reload();
						break;
					case 'zoom':
						this.setZoom(message.level);
						break;
					case 'iframe-navigate':
						this.handleIframeNavigation(message.url);
						break;
				}
			},
			null,
			this._disposables
		);
	}

	private navigate(url: string) {
		try {
			// Clean and validate URL
			url = url.trim();
			if (!url) {
				vscode.window.showErrorMessage('Please enter a valid URL');
				return;
			}

			// Handle search queries (treat as Google search if no domain detected)
			if (!url.includes('.') && !url.startsWith('http://') && !url.startsWith('https://')) {
				url = `https://www.google.com/search?q=${encodeURIComponent(url)}`;
			}
			// Ensure URL has protocol
			else if (!url.startsWith('http://') && !url.startsWith('https://')) {
				url = 'https://' + url;
			}

			// Basic URL validation
			try {
				new URL(url);
			} catch {
				vscode.window.showErrorMessage(`Invalid URL: ${url}`);
				return;
			}
			
			this._currentUrl = url;
			
			// Add to history if it's a new navigation (not back/forward)
			if (this._historyIndex < this._history.length - 1) {
				this._history = this._history.slice(0, this._historyIndex + 1);
			}
			this._history.push(url);
			this._historyIndex = this._history.length - 1;
			
			this._update();
			this._panel.title = `Browser Preview - ${this.getDomainFromUrl(url)}`;
		} catch (error) {
			vscode.window.showErrorMessage(`Failed to navigate to: ${url}`);
		}
	}

	private handleIframeNavigation(url: string) {
		// Update history when iframe navigates via clicked links
		this._currentUrl = url;
		
		// Add to history
		if (this._historyIndex < this._history.length - 1) {
			this._history = this._history.slice(0, this._historyIndex + 1);
		}
		this._history.push(url);
		this._historyIndex = this._history.length - 1;
		
		this._updateControls();
		this._panel.title = `Browser Preview - ${this.getDomainFromUrl(url)}`;
		
	}

	private goBack() {
		if (this._historyIndex > 0) {
			this._historyIndex--;
			this._currentUrl = this._history[this._historyIndex];
			this._update();
			this._panel.title = `Browser Preview - ${this.getDomainFromUrl(this._currentUrl)}`;
		}
	}

	private goForward() {
		if (this._historyIndex < this._history.length - 1) {
			this._historyIndex++;
			this._currentUrl = this._history[this._historyIndex];
			this._update();
			this._panel.title = `Browser Preview - ${this.getDomainFromUrl(this._currentUrl)}`;
		}
	}

	private reload() {
		this._update();
	}

	private setZoom(level: number) {
		this._zoomLevel = Math.max(0.5, Math.min(3.0, level));
		this._updateZoom();
	}

	private _updateControls() {
		this._panel.webview.postMessage({
			type: 'updateControls',
			url: this._currentUrl,
			canGoBack: this._historyIndex > 0,
			canGoForward: this._historyIndex < this._history.length - 1,
			zoomLevel: this._zoomLevel
		});
	}

	private _updateZoom() {
		this._panel.webview.postMessage({
			type: 'updateZoom',
			zoomLevel: this._zoomLevel
		});
	}

	private getDomainFromUrl(url: string): string {
		try {
			const urlObj = new URL(url);
			return urlObj.hostname;
		} catch {
			return 'Unknown';
		}
	}

	public navigateToUrl(url: string) {
		this.navigate(url);
	}

	public dispose() {
		BrowserPanel.currentPanel = undefined;

		// Clean up resources
		this._panel.dispose();

		while (this._disposables.length) {
			const x = this._disposables.pop();
			if (x) {
				x.dispose();
			}
		}
	}

	public getState() {
		return {
			currentUrl: this._currentUrl,
			history: this._history,
			historyIndex: this._historyIndex,
			zoomLevel: this._zoomLevel
		};
	}


	private _update() {
		const webview = this._panel.webview;
		this._panel.webview.html = this._getHtmlForWebview(webview);
	}

	private _getHtmlForWebview(webview: vscode.Webview) {
		const nonce = getNonce();

		return `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta http-equiv="Content-Security-Policy" content="default-src 'none'; frame-src https:; style-src ${webview.cspSource} 'unsafe-inline'; script-src 'nonce-${nonce}';">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Browser Preview</title>
				<style>
					:root {
						--vscode-button-background: var(--vscode-button-background, #0e639c);
						--vscode-button-foreground: var(--vscode-button-foreground, #ffffff);
						--vscode-button-hoverBackground: var(--vscode-button-hoverBackground, #1177bb);
						--vscode-input-background: var(--vscode-input-background, #3c3c3c);
						--vscode-input-foreground: var(--vscode-input-foreground, #cccccc);
						--vscode-input-border: var(--vscode-input-border, #3c3c3c);
						--vscode-focusBorder: var(--vscode-focusBorder, #007acc);
						--vscode-foreground: var(--vscode-foreground, #cccccc);
						--vscode-editorWidget-background: var(--vscode-editorWidget-background, #252526);
					}

					* {
						box-sizing: border-box;
					}

					body {
						margin: 0;
						padding: 0;
						height: 100vh;
						font-family: var(--vscode-font-family);
						font-size: var(--vscode-font-size);
						color: var(--vscode-foreground);
						background-color: var(--vscode-editorWidget-background);
						display: flex;
						flex-direction: column;
					}

					.browser-toolbar {
						display: flex;
						align-items: center;
						padding: 8px;
						gap: 8px;
						background-color: var(--vscode-editorWidget-background);
						border-bottom: 1px solid var(--vscode-input-border);
						flex-shrink: 0;
					}

					.nav-button, .refresh-button, .zoom-button {
						background: var(--vscode-button-background);
						color: var(--vscode-button-foreground);
						border: none;
						border-radius: 4px;
						padding: 6px 8px;
						cursor: pointer;
						font-size: 12px;
						min-width: 24px;
						height: 28px;
						display: flex;
						align-items: center;
						justify-content: center;
						transition: background-color 0.2s;
					}

					.nav-button:hover:not(:disabled), .refresh-button:hover, .zoom-button:hover {
						background: var(--vscode-button-hoverBackground);
					}

					.nav-button:disabled {
						opacity: 0.5;
						cursor: default;
					}

					.nav-button:focus, .refresh-button:focus, .zoom-button:focus {
						outline: 1px solid var(--vscode-focusBorder);
						outline-offset: 2px;
					}

					.url-input {
						flex: 1;
						background: var(--vscode-input-background);
						color: var(--vscode-input-foreground);
						border: 1px solid var(--vscode-input-border);
						border-radius: 4px;
						padding: 6px 12px;
						font-size: 13px;
						height: 28px;
					}

					.url-input:focus {
						outline: 1px solid var(--vscode-focusBorder);
						outline-offset: -1px;
					}

					.zoom-controls {
						display: flex;
						align-items: center;
						gap: 4px;
					}

					.zoom-level {
						min-width: 40px;
						text-align: center;
						font-size: 11px;
						color: var(--vscode-foreground);
					}

					.browser-container {
						flex: 1;
						position: relative;
						overflow: auto;
						display: flex;
						align-items: center;
						justify-content: center;
					}

					.browser-frame {
						border: none;
						background: white;
						transform-origin: center center;
						transform: scale(var(--zoom-level, 1));
						width: calc(100% / var(--zoom-level, 1));
						height: calc(100% / var(--zoom-level, 1));
						min-width: 100%;
						min-height: 100%;
					}

					.loading {
						display: flex;
						align-items: center;
						justify-content: center;
						flex: 1;
						color: var(--vscode-foreground);
						font-size: 14px;
					}
				</style>
			</head>
			<body>
				<div class="browser-toolbar">
					<button class="nav-button" id="backBtn" title="Go back">‹</button>
					<button class="nav-button" id="forwardBtn" title="Go forward">›</button>
					<button class="refresh-button" id="refreshBtn" title="Reload">⟳</button>
					<input type="text" class="url-input" id="urlInput" placeholder="Enter URL..." value="${this._currentUrl}">
					<div class="zoom-controls">
						<button class="zoom-button" id="zoomOutBtn" title="Zoom out">−</button>
						<span class="zoom-level" id="zoomLevel">${Math.round(this._zoomLevel * 100)}%</span>
						<button class="zoom-button" id="zoomInBtn" title="Zoom in">+</button>
						<button class="zoom-button" id="zoomResetBtn" title="Reset zoom">⚏</button>
					</div>
				</div>
				<div class="browser-container">
					<iframe id="browserFrame" class="browser-frame" src="${this._currentUrl}" title="Browser Preview" style="--zoom-level: ${this._zoomLevel}" sandbox="allow-same-origin allow-scripts allow-forms allow-top-navigation allow-popups"></iframe>
				</div>

				<script nonce="${nonce}">
					const vscode = acquireVsCodeApi();
					
					const backBtn = document.getElementById('backBtn');
					const forwardBtn = document.getElementById('forwardBtn');
					const refreshBtn = document.getElementById('refreshBtn');
					const urlInput = document.getElementById('urlInput');
					const browserFrame = document.getElementById('browserFrame');
					const zoomInBtn = document.getElementById('zoomInBtn');
					const zoomOutBtn = document.getElementById('zoomOutBtn');
					const zoomResetBtn = document.getElementById('zoomResetBtn');
					const zoomLevel = document.getElementById('zoomLevel');

					let canGoBack = ${this._historyIndex > 0};
					let canGoForward = ${this._historyIndex < this._history.length - 1};
					let currentZoom = ${this._zoomLevel};

					function updateButtons() {
						backBtn.disabled = !canGoBack;
						forwardBtn.disabled = !canGoForward;
					}

					function updateZoom(level) {
						currentZoom = level;
						browserFrame.style.setProperty('--zoom-level', level);
						zoomLevel.textContent = Math.round(level * 100) + '%';
					}

					// Navigation buttons
					backBtn.addEventListener('click', () => {
						vscode.postMessage({ type: 'back' });
					});

					forwardBtn.addEventListener('click', () => {
						vscode.postMessage({ type: 'forward' });
					});

					refreshBtn.addEventListener('click', () => {
						vscode.postMessage({ type: 'reload' });
					});

					// URL input
					urlInput.addEventListener('keydown', (e) => {
						if (e.key === 'Enter') {
							const url = urlInput.value.trim();
							if (url) {
								vscode.postMessage({ type: 'navigate', url: url });
							}
						}
					});

					// Zoom controls
					zoomInBtn.addEventListener('click', () => {
						const newZoom = Math.min(3.0, currentZoom + 0.1);
						vscode.postMessage({ type: 'zoom', level: newZoom });
					});

					zoomOutBtn.addEventListener('click', () => {
						const newZoom = Math.max(0.5, currentZoom - 0.1);
						vscode.postMessage({ type: 'zoom', level: newZoom });
					});

					zoomResetBtn.addEventListener('click', () => {
						vscode.postMessage({ type: 'zoom', level: 1.0 });
					});

					// Simplified iframe navigation tracking
					let lastUrl = '${this._currentUrl}';
					let isUserNavigation = false;

					// Function to safely get iframe URL
					const getIframeUrl = () => {
						try {
							return browserFrame.contentWindow.location.href;
						} catch (e) {
							// Cross-origin restriction - return null to indicate we can't access it
							return null;
						}
					};

					// Function to check if URL changed and notify extension
					const checkUrlChange = () => {
						const currentUrl = getIframeUrl();
						if (currentUrl && 
							currentUrl !== lastUrl && 
							currentUrl !== 'about:blank' && 
							!currentUrl.startsWith('about:') &&
							!isUserNavigation) {
							
							console.log('Iframe navigated from', lastUrl, 'to', currentUrl);
							lastUrl = currentUrl;
							vscode.postMessage({ type: 'iframe-navigate', url: currentUrl });
						}
						isUserNavigation = false;
					};

					// Watch for src attribute changes (user navigation)
					const observer = new MutationObserver((mutations) => {
						mutations.forEach((mutation) => {
							if (mutation.type === 'attributes' && mutation.attributeName === 'src') {
								isUserNavigation = true;
								lastUrl = browserFrame.src;
								setTimeout(() => { isUserNavigation = false; }, 1000);
							}
						});
					});

					observer.observe(browserFrame, { 
						attributes: true, 
						attributeFilter: ['src'] 
					});

					// Listen for iframe load events (works for same-origin content)
					browserFrame.addEventListener('load', () => {
						setTimeout(() => {
							if (!isUserNavigation) {
								checkUrlChange();
							}
						}, 500);
					});

					// Try to inject navigation detection script for same-origin iframes
					browserFrame.addEventListener('load', () => {
						try {
							const iframeDoc = browserFrame.contentDocument;
							if (iframeDoc) {
								// Same-origin iframe - we can inject detection script
								console.log('Same-origin iframe detected, injecting navigation tracking');
								const script = iframeDoc.createElement('script');
								script.textContent = \`
									(function() {
										let currentUrl = window.location.href;
										
										const notifyParent = () => {
											if (window.location.href !== currentUrl) {
												currentUrl = window.location.href;
												console.log('Iframe navigation to:', currentUrl);
												window.parent.postMessage({
													type: 'iframe-navigation',
													url: currentUrl
												}, '*');
											}
										};
										
										// Monitor for navigation changes
										window.addEventListener('popstate', notifyParent);
										window.addEventListener('hashchange', notifyParent);
										
										// Override history methods
										const originalPushState = history.pushState;
										const originalReplaceState = history.replaceState;
										
										history.pushState = function() {
											originalPushState.apply(history, arguments);
											setTimeout(notifyParent, 0);
										};
										
										history.replaceState = function() {
											originalReplaceState.apply(history, arguments);
											setTimeout(notifyParent, 0);
										};
										
										// Periodic check for SPA navigation
										setInterval(notifyParent, 1000);
									})();
								\`;
								iframeDoc.head.appendChild(script);
							}
						} catch (e) {
							// Cross-origin iframe - can't inject script
							console.log('Cross-origin iframe - navigation tracking limited');
						}
					});

					// Listen for messages from injected script
					window.addEventListener('message', (event) => {
						if (event.source === browserFrame.contentWindow && event.data.type === 'iframe-navigation') {
							const newUrl = event.data.url;
							if (newUrl && newUrl !== lastUrl) {
								console.log('Iframe navigation detected via injected script:', newUrl);
								lastUrl = newUrl;
								vscode.postMessage({ type: 'iframe-navigate', url: newUrl });
							}
						}
					});

					// Listen for messages from the extension
					window.addEventListener('message', event => {
						const message = event.data;
						switch (message.type) {
							case 'updateControls':
								urlInput.value = message.url;
								browserFrame.src = message.url;
								canGoBack = message.canGoBack;
								canGoForward = message.canGoForward;
								updateButtons();
								updateZoom(message.zoomLevel);
								lastUrl = message.url;
								break;
							case 'updateZoom':
								updateZoom(message.zoomLevel);
								break;
						}
					});

					// Initialize
					updateButtons();
					updateZoom(currentZoom);
				</script>
			</body>
			</html>`;
	}
}

function getNonce() {
	let text = '';
	const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
	for (let i = 0; i < 32; i++) {
		text += possible.charAt(Math.floor(Math.random() * possible.length));
	}
	return text;
}

export async function activate(context: vscode.ExtensionContext) {
	if (typeof navigator !== 'object') {	// do not run under node.js
		console.error("not running in browser");
		return;
	}

	// Register command to open browser
	const openBrowserCommand = vscode.commands.registerCommand('preview.openBrowser', () => {
		BrowserPanel.createOrShow(context.extensionUri, context);
	});

	context.subscriptions.push(openBrowserCommand);

	// Note: Webview panel serialization disabled to prevent loading issues
	// Users can simply reopen the browser preview after VSCode restart
}

