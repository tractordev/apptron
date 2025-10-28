
export function modalDialog(el) {
    const closers = el.querySelectorAll('[data-action="close"]');
    closers.forEach(closer => {
        closer.addEventListener("click", () => el.close());
    });
    el.addEventListener("click", (e) => {
        const r = el.getBoundingClientRect();
        const inBounds = (
            e.clientX >= r.left &&
            e.clientX <= r.right &&
            e.clientY >= r.top &&
            e.clientY <= r.bottom
        );
        if (!inBounds) el.close();
    });
    return el;
}