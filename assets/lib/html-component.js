export class HTMLComponent extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: 'open' });
    this._scriptContexts = new Map();
  }
  
  static get observedAttributes() {
    return ['src'];
  }
  
  connectedCallback() {
    const src = this.getAttribute('src');
    if (src) {
      this.load(src);
    }
  }
  
  async load(url) {
    try {
      const response = await fetch(url);
      const html = await response.text();
      
      // Parse and process the template
      const parser = new DOMParser();
      const doc = parser.parseFromString(html, 'text/html');
      const com = doc.querySelector('html');
      
      if (!com) {
        throw new Error('No component found');
      }
      
      // Create a document fragment from template content
      const fragment = document.createDocumentFragment();
      const tempContainer = document.createElement('div');
      tempContainer.innerHTML = com.innerHTML;
      
      // Process each node
      this.processNodes(tempContainer, fragment);
      
      // Clear and append to shadow root
      this.shadowRoot.innerHTML = '';
      this.shadowRoot.appendChild(fragment);
      
      // Initialize any scripts after DOM is ready
      await this.initializeScripts();
      
      // Fire loaded event
      this.dispatchEvent(new CustomEvent('loaded', { 
        detail: { url },
        bubbles: true 
      }));
      
    } catch (error) {
      this.handleError(error);
    }
  }
  
  processNodes(source, target) {
    Array.from(source.childNodes).forEach(node => {
      if (node.nodeType === Node.ELEMENT_NODE) {
        if (node.tagName === 'SCRIPT') {
          // Store script for later execution
          this._scriptContexts.set(node, node.textContent || node.src);
        } else {
          // Clone and append other elements
          const cloned = node.cloneNode(true);
          target.appendChild(cloned);
        }
      } else {
        // Clone text nodes and comments
        target.appendChild(node.cloneNode(true));
      }
    });
  }
  
  async initializeScripts() {
    for (const [scriptNode, content] of this._scriptContexts) {
      if (scriptNode.src) {
        // Load external script
        const script = document.createElement('script');
        script.src = scriptNode.src;
        script.type = scriptNode.type || 'text/javascript';
        this.shadowRoot.appendChild(script);
      } else {
        // Execute inline script with shadow DOM context (supports import and await)
        await this.executeInlineScript(content);
      }
    }
    
    // Clear stored scripts
    this._scriptContexts.clear();
  }
  
  async executeInlineScript(code) {
    try {
      // Create a unique context ID for this script execution
      const contextId = `__htmlComponent_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
      
      // Store context in a temporary global variable
      window[contextId] = {
        shadowRoot: this.shadowRoot,
        component: this,
        $: (selector) => this.shadowRoot.querySelector(selector),
        $$: (selector) => this.shadowRoot.querySelectorAll(selector)
      };
      
      // Resolve relative imports to absolute URLs
      // This is necessary because blob URLs don't have a proper base for resolution
      const resolvedCode = this.resolveImports(code);
      
      // Create a module script that retrieves the context
      // This allows use of import and await
      const moduleCode = `
        const { shadowRoot, component, $, $$ } = window['${contextId}'];
        
        // Clean up the global context immediately after retrieval
        delete window['${contextId}'];
        
        // Execute the original script (with import and await support)
        ${resolvedCode}
      `;
      
      // Create a blob URL for the module
      const blob = new Blob([moduleCode], { type: 'text/javascript' });
      const url = URL.createObjectURL(blob);
      
      try {
        // Import and execute as a module
        await import(url);
      } finally {
        // Clean up the blob URL
        URL.revokeObjectURL(url);
        // Ensure context is cleaned up even if script fails
        delete window[contextId];
      }
      
    } catch (error) {
      console.error('Script execution error:', error);
    }
  }
  
  resolveImports(code) {
    // Replace import statements with absolute URLs
    // Matches: import ... from "path" or import ... from 'path'
    return code.replace(
      /import\s+(?:(?:\{[^}]*\}|\*\s+as\s+\w+|\w+)(?:\s*,\s*(?:\{[^}]*\}|\*\s+as\s+\w+|\w+))*\s+from\s+)?['"]([^'"]+)['"]/g,
      (match, path) => {
        // Convert relative paths to absolute URLs
        const absoluteUrl = new URL(path, window.location.origin + window.location.pathname).href;
        return match.replace(path, absoluteUrl);
      }
    );
  }
  
  handleError(error) {
    console.error('Component loading error:', error);
    this.shadowRoot.innerHTML = `
      <div style="padding: 10px; background: #fee; color: #c00; border: 1px solid #c00; border-radius: 4px;">
        <strong>Error:</strong> ${error.message}
      </div>
    `;
    
    this.dispatchEvent(new CustomEvent('error', {
      detail: { error },
      bubbles: true
    }));
  }
  
  // Utility methods for external access
  getElement(selector) {
    return this.shadowRoot.querySelector(selector);
  }
  
  getElements(selector) {
    return this.shadowRoot.querySelectorAll(selector);
  }
}
  
customElements.define('html-component', HTMLComponent);