/**
 * InputCapturer - Captures keyboard and mouse input and sends it via MessageChannel.
 * 
 * Usage:
 *   const capturer = new InputCapturer(targetElement);
 *   worker.postMessage({ inputPort: capturer.port }, [capturer.port]);
 * 
 * The port sends messages with these types:
 *   - { type: "keyboard", code: number } - PS/2 scan code
 *   - { type: "mouse-click", left: bool, middle: bool, right: bool }
 *   - { type: "mouse-delta", deltaX: number, deltaY: number }
 *   - { type: "mouse-wheel", deltaX: number, deltaY: number }
 */
export class InputCapturer {
    // PS/2 scan code mapping from event.code
    static CODEMAP = {
        "Escape": 0x0001,
        "Digit1": 0x0002,
        "Digit2": 0x0003,
        "Digit3": 0x0004,
        "Digit4": 0x0005,
        "Digit5": 0x0006,
        "Digit6": 0x0007,
        "Digit7": 0x0008,
        "Digit8": 0x0009,
        "Digit9": 0x000a,
        "Digit0": 0x000b,
        "Minus": 0x000c,
        "Equal": 0x000d,
        "Backspace": 0x000e,
        "Tab": 0x000f,
        "KeyQ": 0x0010,
        "KeyW": 0x0011,
        "KeyE": 0x0012,
        "KeyR": 0x0013,
        "KeyT": 0x0014,
        "KeyY": 0x0015,
        "KeyU": 0x0016,
        "KeyI": 0x0017,
        "KeyO": 0x0018,
        "KeyP": 0x0019,
        "BracketLeft": 0x001a,
        "BracketRight": 0x001b,
        "Enter": 0x001c,
        "ControlLeft": 0x001d,
        "KeyA": 0x001e,
        "KeyS": 0x001f,
        "KeyD": 0x0020,
        "KeyF": 0x0021,
        "KeyG": 0x0022,
        "KeyH": 0x0023,
        "KeyJ": 0x0024,
        "KeyK": 0x0025,
        "KeyL": 0x0026,
        "Semicolon": 0x0027,
        "Quote": 0x0028,
        "Backquote": 0x0029,
        "ShiftLeft": 0x002a,
        "Backslash": 0x002b,
        "KeyZ": 0x002c,
        "KeyX": 0x002d,
        "KeyC": 0x002e,
        "KeyV": 0x002f,
        "KeyB": 0x0030,
        "KeyN": 0x0031,
        "KeyM": 0x0032,
        "Comma": 0x0033,
        "Period": 0x0034,
        "Slash": 0x0035,
        "IntlRo": 0x0035,
        "ShiftRight": 0x0036,
        "NumpadMultiply": 0x0037,
        "AltLeft": 0x0038,
        "Space": 0x0039,
        "CapsLock": 0x003a,
        "F1": 0x003b,
        "F2": 0x003c,
        "F3": 0x003d,
        "F4": 0x003e,
        "F5": 0x003f,
        "F6": 0x0040,
        "F7": 0x0041,
        "F8": 0x0042,
        "F9": 0x0043,
        "F10": 0x0044,
        "NumLock": 0x0045,
        "ScrollLock": 0x0046,
        "Numpad7": 0x0047,
        "Numpad8": 0x0048,
        "Numpad9": 0x0049,
        "NumpadSubtract": 0x004a,
        "Numpad4": 0x004b,
        "Numpad5": 0x004c,
        "Numpad6": 0x004d,
        "NumpadAdd": 0x004e,
        "Numpad1": 0x004f,
        "Numpad2": 0x0050,
        "Numpad3": 0x0051,
        "Numpad0": 0x0052,
        "NumpadDecimal": 0x0053,
        "IntlBackslash": 0x0056,
        "F11": 0x0057,
        "F12": 0x0058,
        "NumpadEnter": 0xe01c,
        "ControlRight": 0xe01d,
        "NumpadDivide": 0xe035,
        "AltRight": 0xe038,
        "Home": 0xe047,
        "ArrowUp": 0xe048,
        "PageUp": 0xe049,
        "ArrowLeft": 0xe04b,
        "ArrowRight": 0xe04d,
        "End": 0xe04f,
        "ArrowDown": 0xe050,
        "PageDown": 0xe051,
        "Insert": 0xe052,
        "Delete": 0xe053,
        "OSLeft": 0xe05b,
        "MetaLeft": 0xe05b,
        "OSRight": 0xe05c,
        "MetaRight": 0xe05c,
        "ContextMenu": 0xe05d,
    };

    static SCAN_CODE_RELEASE = 0x80;

    /**
     * @param {HTMLElement} target - Element to use for pointer lock (typically a canvas)
     * @param {Object} options - Configuration options
     * @param {number} options.doubleEscapeDelay - Ms to wait for second escape (default: 500)
     */
    constructor(target, options = {}) {
        this.target = target;
        this.doubleEscapeDelay = options.doubleEscapeDelay ?? 500;

        // Create message channel
        this._channel = new MessageChannel();
        this._localPort = this._channel.port1;
        this.port = this._channel.port2;

        // State
        this._captured = false;
        this._keysPressed = {};
        this._escapeCount = 0;
        this._escapeTimeout = null;
        this._leftDown = false;
        this._middleDown = false;
        this._rightDown = false;
        this._destroyed = false;

        // Public state
        this.mouseEnabled = true;

        // Callback for capture state changes
        this.onCaptureChange = null;

        // Bind handlers
        this._onKeyDown = this._onKeyDown.bind(this);
        this._onKeyUp = this._onKeyUp.bind(this);
        this._onMouseMove = this._onMouseMove.bind(this);
        this._onMouseDown = this._onMouseDown.bind(this);
        this._onMouseUp = this._onMouseUp.bind(this);
        this._onWheel = this._onWheel.bind(this);
        this._onPointerLockChange = this._onPointerLockChange.bind(this);
        this._onWindowBlur = this._onWindowBlur.bind(this);
        this._onContextMenu = this._onContextMenu.bind(this);
        this._onClick = this._onClick.bind(this);

        this._attach();
    }

    /**
     * Whether input is currently captured
     */
    get captured() {
        return this._captured;
    }

    /**
     * Request input capture (pointer lock)
     */
    capture() {
        if (!this._captured && !this._destroyed) {
            // this.target.requestPointerLock();
            this._setCaptured(true);
        }
    }

    /**
     * Release input capture
     */
    release() {
        if (this._captured) {
            if (document.pointerLockElement === this.target) {
                document.exitPointerLock();
            }
            this._setCaptured(false);
        }
    }

    /**
     * Clean up all event listeners
     */
    destroy() {
        this._destroyed = true;
        this.release();
        this._detach();
        this._localPort.close();
    }

    _attach() {
        this.target.addEventListener("click", this._onClick);
        this.target.addEventListener("contextmenu", this._onContextMenu);
        document.addEventListener("pointerlockchange", this._onPointerLockChange);
        window.addEventListener("keydown", this._onKeyDown);//, true);
        window.addEventListener("keyup", this._onKeyUp);//, true);
        window.addEventListener("mousemove", this._onMouseMove);
        window.addEventListener("mousedown", this._onMouseDown);
        window.addEventListener("mouseup", this._onMouseUp);
        window.addEventListener("wheel", this._onWheel, { passive: false });
        window.addEventListener("blur", this._onWindowBlur);
    }

    _detach() {
        this.target.removeEventListener("click", this._onClick);
        this.target.removeEventListener("contextmenu", this._onContextMenu);
        document.removeEventListener("pointerlockchange", this._onPointerLockChange);
        window.removeEventListener("keydown", this._onKeyDown);//, true);
        window.removeEventListener("keyup", this._onKeyUp);//, true);
        window.removeEventListener("mousemove", this._onMouseMove);
        window.removeEventListener("mousedown", this._onMouseDown);
        window.removeEventListener("mouseup", this._onMouseUp);
        window.removeEventListener("wheel", this._onWheel, { passive: false });
        window.removeEventListener("blur", this._onWindowBlur);
    }

    _setCaptured(captured) {
        if (this._captured !== captured) {
            this._captured = captured;
            
            if (!captured) {
                this._releaseAllKeys();
                this._releaseMouseButtons();
                this._escapeCount = 0;
                clearTimeout(this._escapeTimeout);
            }

            if (this.onCaptureChange) {
                this.onCaptureChange(captured);
            }
        }
    }

    _send(message) {
        this._localPort.postMessage(message);
    }

    _sendKeyCode(code, keydown) {
        this._keysPressed[code] = keydown;

        const RELEASE = InputCapturer.SCAN_CODE_RELEASE;

        if (!keydown) {
            code |= RELEASE;
        }

        // For extended codes (> 0xFF), send prefix first
        if (code > 0xFF) {
            this._send({ type: "keyboard", codes: [code >> 8, code & 0xFF] });
        } else {
            this._send({ type: "keyboard", codes: [code] });
        }
    }

    _releaseAllKeys() {
        for (const code in this._keysPressed) {
            if (this._keysPressed[code]) {
                this._sendKeyCode(parseInt(code), false);
            }
        }
        this._keysPressed = {};
    }

    _releaseMouseButtons() {
        if (this._leftDown || this._middleDown || this._rightDown) {
            this._leftDown = this._middleDown = this._rightDown = false;
            this._send({ type: "mouse-click", left: false, middle: false, right: false });
        }
    }

    _onClick(e) {
        if (!this._captured) {
            this.capture();
        }
    }

    _onContextMenu(e) {
        if (this._captured) {
            e.preventDefault();
        }
    }

    _onPointerLockChange() {
        if (document.pointerLockElement === this.target) {
            this._setCaptured(true);
        } else {
            this._setCaptured(false);
        }
    }

    _onWindowBlur() {
        if (this._captured) {
            this._releaseAllKeys();
            this._releaseMouseButtons();
        }
    }

    _onKeyDown(e) {
        if (!this._captured) return;

        // Double-escape to release
        if (e.code === "Escape") {
            e.preventDefault();
            e.stopPropagation();

            this._escapeCount++;

            if (this._escapeCount >= 2) {
                clearTimeout(this._escapeTimeout);
                this._escapeCount = 0;
                this.release();
                return;
            }

            // First escape - send to VM
            const scancode = InputCapturer.CODEMAP["Escape"];
            this._sendKeyCode(scancode, true);

            clearTimeout(this._escapeTimeout);
            this._escapeTimeout = setTimeout(() => {
                this._escapeCount = 0;
            }, this.doubleEscapeDelay);

            return;
        }

        // Allow browser dev tools
        if (e.shiftKey && e.ctrlKey && (e.keyCode === 73 || e.keyCode === 74 || e.keyCode === 75)) {
            return;
        }

        e.preventDefault();
        e.stopPropagation();

        const scancode = InputCapturer.CODEMAP[e.code];
        if (scancode !== undefined) {
            // Handle key repeat
            if (this._keysPressed[scancode] && !e.repeat) {
                this._sendKeyCode(scancode, false);
            }
            this._sendKeyCode(scancode, true);
        }
    }

    _onKeyUp(e) {
        if (!this._captured) return;

        e.preventDefault();
        e.stopPropagation();

        const scancode = InputCapturer.CODEMAP[e.code];
        if (scancode !== undefined && this._keysPressed[scancode]) {
            this._sendKeyCode(scancode, false);
        }
    }

    _onMouseMove(e) {
        if (!this._captured || !this.mouseEnabled) return;

        let deltaX = e.movementX || 0;
        let deltaY = e.movementY || 0;

        // Invert Y for PS/2 mouse protocol (positive = up)
        deltaY = -deltaY;

        if (deltaX !== 0 || deltaY !== 0) {
            this._send({ type: "mouse-delta", deltaX, deltaY });
        }
    }

    _onMouseDown(e) {
        if (!this._captured || !this.mouseEnabled) return;

        e.preventDefault();

        if (e.button === 0) this._leftDown = true;
        else if (e.button === 1) this._middleDown = true;
        else if (e.button === 2) this._rightDown = true;

        this._send({ 
            type: "mouse-click", 
            left: this._leftDown, 
            middle: this._middleDown, 
            right: this._rightDown 
        });
    }

    _onMouseUp(e) {
        if (!this._captured || !this.mouseEnabled) return;

        e.preventDefault();

        if (e.button === 0) this._leftDown = false;
        else if (e.button === 1) this._middleDown = false;
        else if (e.button === 2) this._rightDown = false;

        this._send({ 
            type: "mouse-click", 
            left: this._leftDown, 
            middle: this._middleDown, 
            right: this._rightDown 
        });
    }

    _onWheel(e) {
        if (!this._captured || !this.mouseEnabled) return;

        e.preventDefault();

        let deltaX = 0;
        let deltaY = 0;

        if (e.deltaY < 0) deltaX = 1;      // Scroll up
        else if (e.deltaY > 0) deltaX = -1; // Scroll down

        this._send({ type: "mouse-wheel", deltaX, deltaY });
    }
}

