// /lib/toast.js
//
// Minimal toast helper that matches the exact HTML structure:
//
// <div class="toast">
//	 <svg class="toast-status-icon">…</svg>
//	 <span class="toast-text">Message</span>
//	 <button><svg>…</svg></button>
// </div>

const TOAST_ROOT_ID = "toast-root";
const DEFAULT_TIMEOUT = 4000;

// Ensure #toast-root exists
function ensureToastRoot() {
	let root = document.getElementById(TOAST_ROOT_ID);
	if (root) return root;

	root = document.createElement("div");
	root.id = TOAST_ROOT_ID;
	document.body.appendChild(root);
	return root;
}

// Load and clone an external SVG file
async function loadSvg(url, className) {
	const resp = await fetch(url);
	if (!resp.ok) {
		throw new Error(`Failed to load SVG: ${url}`);
	}

	const text = await resp.text();
	const template = document.createElement("template");
	template.innerHTML = text.trim();

	const svg = template.content.querySelector("svg");
	if (!svg) {
		throw new Error(`No <svg> found in ${url}`);
	}

	if (className) {
		svg.classList.add(className);
	}

	return svg;
}

export async function showToast(
	message,
	{
		timeout = DEFAULT_TIMEOUT,
		icon = "/img/icons/info.svg",
	} = {}
) {
	const root = ensureToastRoot();

	const toast = document.createElement("div");
	toast.className = "toast";

	// Status icon
	let statusIcon;
	try {
		statusIcon = await loadSvg(icon, "toast-status-icon");
	} catch (err) {
		console.error(err);
	}

	// Message
	const text = document.createElement("span");
	text.className = "toast-text";
	text.textContent = message;

	// Close button
	const closeBtn = document.createElement("button");
	closeBtn.type = "button";
	closeBtn.setAttribute("aria-label", "Close");

	let closeIcon;
	try {
		closeIcon = await loadSvg("/img/icons/close.svg");
		closeBtn.appendChild(closeIcon);
	} catch (err) {
		console.error(err);
		closeBtn.textContent = "×";
	}

	closeBtn.addEventListener("click", () => {
		toast.remove();
	});

	// Assemble toast
	if (statusIcon) toast.appendChild(statusIcon);
	toast.appendChild(text);
	toast.appendChild(closeBtn);

	root.appendChild(toast);

	// Auto-dismiss
	if (timeout > 0) {
		setTimeout(() => {
			toast.remove();
		}, timeout);
	}

	return {
		el: toast,
		dismiss() {
			toast.remove();
		},
	};
}