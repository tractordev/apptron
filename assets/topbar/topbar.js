import { register } from "https://esm.run/@teamhanko/hanko-elements";

class TopBar extends HTMLElement {
    constructor() {
        super();
        this.attachShadow({ mode: 'open' });
        this.loadTemplate();
    }

    async setup() {
        const { hanko } = await register("https://ad6044b5-53c2-4cb5-8542-9fdaef75f771.hanko.io");
        window.hanko = hanko;
        hanko.onSessionCreated(() => {
            window.location.reload();
        });
        const session = await hanko.validateSession();
        console.log("Session valid:", session.is_valid, "Claims:", session.claims);
        if (session.is_valid) {
            document.body.classList.add('signedin');
            this.shadowRoot.getElementById('header-bar').classList.add("signedin");
        }



        const signinBtn = this.shadowRoot.querySelector('#signin');
        const signinDialog = this.shadowRoot.getElementById('signin-dialog');
        const signinCloseBtn = this.shadowRoot.querySelector('.signin-close');
        
        // Open dialog
        signinBtn.addEventListener('click', () => {
            signinDialog.showModal();
        });
        
        // Close dialog when clicking X
        signinCloseBtn.addEventListener('click', () => {
            signinDialog.close();
        });
        
        // Optional: Close when clicking backdrop (outside dialog)
        signinDialog.addEventListener('click', (e) => {
            if (e.target === signinDialog) {
                signinDialog.close();
            }
        });



        const accountBtn = this.shadowRoot.querySelector('#account');
        const accountDialog = this.shadowRoot.getElementById('account-dialog');
        const accountCloseBtn = this.shadowRoot.querySelector('.account-close');
        
        // Open dialog
        accountBtn.addEventListener('click', () => {
            accountDialog.showModal();
        });
        
        // Close dialog when clicking X
        accountCloseBtn.addEventListener('click', () => {
            accountDialog.close();
        });
        
        // Optional: Close when clicking backdrop (outside dialog)
        accountDialog.addEventListener('click', (e) => {
            if (e.target === accountDialog) {
                accountDialog.close();
            }
        });



        const logoutBtn = this.shadowRoot.querySelector('#logout');
        logoutBtn.addEventListener('click', async () => {
            console.log("logout");
            await window.hanko.user.logout();
            window.location.reload();
        });
    }
    
    async loadTemplate() {
        try {
            // Fetch the template from external file
            const response = await fetch('/topbar/topbar.html');
            const text = await response.text();
            
            // Parse the template
            const parser = new DOMParser();
            const doc = parser.parseFromString(text, 'text/html');
            const template = doc.getElementById('topbar');
            
            if (template) {
                // Clone and attach the template content to shadow DOM
                this.shadowRoot.appendChild(template.content.cloneNode(true));
                
                // Setup dialog functionality
                this.setup();
            } else {
                throw new Error('Template not found');
            }
        } catch (error) {
            console.error('Failed to load top-bar template:', error);
            // Fallback content if template fails to load
            this.shadowRoot.innerHTML = `
                <style>
                    :host { display: block; }
                    .error { 
                        background: #333; 
                        color: white; 
                        padding: 1rem; 
                        text-align: center;
                    }
                </style>
                <div class="error">Top bar failed to load</div>
            `;
        }
    }

}

// Register the custom element
customElements.define('apptron-topbar', TopBar);