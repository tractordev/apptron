class Dropdown extends HTMLElement {
  connectedCallback() {
    const children = this.children;

    if (children.length < 2) {
      console.warn('popover-setup requires at least 2 child elements');
      return;
    }

    const trigger = children[0];
    const popover = children[1];

    if (!popover.id) {
      popover.id = `popover-${Math.random().toString(36).substr(2, 9)}`;
    }

    trigger.setAttribute('aria-haspopup', 'menu');
    trigger.setAttribute('aria-controls', popover.id);
    trigger.setAttribute('popovertarget', popover.id);
    trigger.setAttribute('popovertargetaction', 'toggle');

    popover.setAttribute('popover', '');

    // ===== click-outside handler =====
    // Store handler on the instance so we can remove it later
    this._outsideClickHandler = (event) => {
      const path = event.composedPath();
      // If click is inside trigger or popover, ignore
      if (path.includes(trigger) || path.includes(popover)) return;
      // Otherwise hide the popover
      if (typeof popover.hidePopover === 'function') {
        popover.hidePopover();
      } else {
        // Fallback if HTML popover API isn't present
        popover.removeAttribute('open');
      }
    };

    popover.addEventListener('toggle', (e) => {
      if (e.newState === 'open') {
        // Positioning logic
        const triggerRect = trigger.getBoundingClientRect();
        const popoverRect = popover.getBoundingClientRect();

        const align = this.getAttribute('align') || 'left';

        let marginLeft;
        if (align === 'right') {
          marginLeft = triggerRect.right - popoverRect.width;
        } else {
          marginLeft = triggerRect.left;
        }

        const marginTop = triggerRect.bottom;

        popover.style.marginLeft = `${marginLeft}px`;
        popover.style.marginTop = `${marginTop}px`;

        // Start listening for outside clicks (capture so we run early)
        document.addEventListener('pointerdown', this._outsideClickHandler, true);
      } else {
        // Popover closed: stop listening
        document.removeEventListener('pointerdown', this._outsideClickHandler, true);
      }
    });
  }

  disconnectedCallback() {
    // Clean up if the element is removed while open
    if (this._outsideClickHandler) {
      document.removeEventListener('pointerdown', this._outsideClickHandler, true);
    }
  }
}

customElements.define('apptron-dropdown', Dropdown);