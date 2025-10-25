class Dropdown extends HTMLElement {
    connectedCallback() {
      // Get the first two child elements
      const children = this.children;
      
      if (children.length < 2) {
        console.warn('popover-setup requires at least 2 child elements');
        return;
      }
      
      const trigger = children[0];
      const popover = children[1];
      
      // Generate a unique ID for the popover if it doesn't have one
      if (!popover.id) {
        popover.id = `popover-${Math.random().toString(36).substr(2, 9)}`;
      }
      
      // Set up the trigger element
      trigger.setAttribute('aria-haspopup', 'menu');
      trigger.setAttribute('aria-controls', popover.id);
      trigger.setAttribute('popovertarget', popover.id);
      trigger.setAttribute('popovertargetaction', 'toggle');
      
      // Set up the popover element
      popover.setAttribute('popover', '');
      
      // Set up positioning
      popover.addEventListener('toggle', (e) => {
        if (e.newState !== 'open') return;
        
        const triggerRect = trigger.getBoundingClientRect();
        const popoverRect = popover.getBoundingClientRect();
        
        // Get alignment from attribute (default: left)
        const align = this.getAttribute('align') || 'left';
        
        let marginLeft;
        if (align === 'right') {
          // Align right edges
          marginLeft = triggerRect.right - popoverRect.width;
        } else {
          // Align left edges (default)
          marginLeft = triggerRect.left;
        }
        
        const marginTop = triggerRect.bottom;
        
        popover.style.marginLeft = `${marginLeft}px`;
        popover.style.marginTop = `${marginTop}px`;
      });
    }
  }
  
  customElements.define('apptron-dropdown', Dropdown);