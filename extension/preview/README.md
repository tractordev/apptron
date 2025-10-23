# Browser Preview Extension

A VSCode extension that provides a webview-based browser panel with navigation controls, designed to be consistent with VSCode themes. Opens as a full editor tab rather than a sidebar.

## Features

- **Editor-based Browser Panel**: Opens as a full editor tab for maximum screen real estate
- **Iframe-based Browser**: Uses iframe for web content rendering to bypass VSCode webview restrictions
- **Navigation Controls**: Back, forward, and reload buttons with full history management
- **Iframe Navigation Tracking**: Automatically tracks when users click links within the iframe and updates history
- **URL Bar**: Enter URLs or search terms directly
- **Zoom Controls**: Zoom in/out and reset zoom for better content viewing
- **Theme Integration**: Styled to match VSCode's current theme using CSS variables
- **Context Menu Support**: Right-click on selected text to open URLs in the browser preview
- **Smart URL Handling**: 
  - Automatically adds `https://` protocol if missing
  - Treats non-URL text as Google search queries
  - Validates URLs before navigation
- **Panel Persistence**: Maintains state when VSCode is restarted

## Installation

1. Build the extension:
   ```bash
   npm install
   npm run compile-web
   ```

2. The extension will be compiled to `dist/web/extension.js`

## Usage

### Opening the Browser Preview

1. **Command Palette**: Press `Cmd+Shift+P` (macOS) or `Ctrl+Shift+P` (Windows/Linux) and type "Open Browser Preview"
2. **Context Menu**: Select a URL in any file, right-click, and choose "Open in Browser Preview"

The browser preview will open as a new editor tab alongside your code.

### Navigation

- **URL Bar**: Type any URL or search term and press Enter
- **Back/Forward**: Use the `‹` and `›` buttons to navigate through history
- **Reload**: Use the `⟳` button to refresh the current page
- **Sync URL**: Use the `⟲` button to manually sync the URL bar with the current iframe location
- **Zoom Controls**: Use `−` and `+` buttons to zoom out/in, or `⚏` to reset zoom to 100%
- **Link Tracking**: Automatic for same-origin sites, manual sync available for cross-origin sites

### Smart URL Handling

The extension intelligently handles different types of input:

- **Full URLs**: `https://example.com` → Opens directly
- **Domain names**: `example.com` → Automatically adds `https://`
- **Search queries**: `javascript tutorial` → Searches on Google
- **Selected text**: Highlight any text and use the context menu to open it

## Architecture

The extension uses:

- **WebView Panel API**: Creates a custom webview panel in VSCode that opens as an editor tab
- **Iframe**: Embeds web content using iframe to bypass security restrictions
- **Message Passing**: Communication between extension and webview via postMessage
- **History Management**: Tracks navigation history for back/forward functionality, including iframe navigation
- **Zoom Transformation**: CSS transforms to scale iframe content
- **Navigation Monitoring**: Periodic checking and event listening for iframe URL changes
- **CSS Variables**: Uses VSCode theme variables for consistent styling
- **Panel Persistence**: Automatically restores browser panels when VSCode restarts

## Files

- `src/web/extension.ts` - Main extension logic and webview provider
- `package.json` - Extension manifest with commands and contributions
- `dist/web/extension.js` - Compiled extension bundle

## Commands

- `preview.openBrowser` - Opens the browser preview panel
- `preview.openUrlInBrowser` - Opens selected text as URL in browser preview

## Theme Integration

The extension automatically adapts to your current VSCode theme using CSS variables:

- `--vscode-button-background`
- `--vscode-button-foreground`
- `--vscode-input-background`
- `--vscode-input-foreground`
- `--vscode-focusBorder`
- And more...

## Development

To modify the extension:

1. Edit `src/web/extension.ts`
2. Run `npm run compile-web` to rebuild
3. Reload the VSCode window to test changes

## Inspiration

This extension takes inspiration from the browser preview functionality in the [vscode-livepreview](https://github.com/microsoft/vscode-livepreview) extension, focusing specifically on the browser panel component with enhanced navigation controls and theme integration.

## Browser Compatibility

The iframe-based approach works with most websites, though some sites may have security policies that prevent embedding. This is a limitation of iframe technology, not the extension itself.

**Navigation Tracking Limitations**: 

Due to browser security restrictions, navigation tracking has different capabilities depending on the website:

- **Same-Origin Sites** (localhost, file://, etc.): Full automatic navigation tracking works perfectly
- **Cross-Origin Sites** (most external websites): Cannot automatically detect navigation due to browser security policies
- **Manual Sync**: Use the `⟲` sync button to manually update the URL bar with the current page location
- **Workaround**: For cross-origin sites, you can manually edit the URL bar or use the sync button after navigating

## Key Improvements

### Editor vs Sidebar
- **Full Screen Experience**: Opens as an editor tab, giving you maximum screen real estate for browsing
- **Better Integration**: Behaves like any other VSCode editor tab with proper focus management

### Enhanced Navigation
- **Link Click Tracking**: Automatically detects when you click links in the iframe and updates the browser controls
- **Comprehensive History**: Maintains full browsing history including iframe navigation
- **Smart URL Updates**: URL bar automatically updates as you navigate

### Zoom Features
- **Flexible Zooming**: Zoom range from 50% to 300% with 10% increments
- **Visual Feedback**: Real-time zoom percentage display
- **Reset Functionality**: Quick reset to 100% zoom