// @ts-check

document.addEventListener("DOMContentLoaded", main);
function main() {
    const anchors = document.querySelectorAll("a");
    if (anchors.length == 0) return;

    for (const anchor of anchors) {
        if (hasHTMXAttributes(anchor)) continue;

        anchor.addEventListener("click", async (e) => {
            if (e.detail["htmxIgnore"]) {
                return;
            }

            e.preventDefault();
            let content = "";
            let href = "";

            try {
                const anchor = /** @type {HTMLAnchorElement | undefined} */ (
                    e.currentTarget
                );
                if (!anchor) {
                    throw new Error("failed to acquire anchor");
                }

                href = anchor.href;
                const res = await fetch(href);
                content = await res.text();
            } catch (err) {
                console.warn("failed to highjack event", err);
                e.target?.dispatchEvent(
                    new CustomEvent("click", { detail: { htmxIgnore: true } }),
                );
                return;
            }

            setPageContent(content);
            window.history.pushState(content, "", href);
            document.dispatchEvent(new Event("DOMContentLoaded"));
        });
    }

    // make the back button work properly.
    window.addEventListener("popstate", (e) => {
        if (!e.state) return;
        setPageContent(e.state);
    });
}

/**
 * @param {string} state
 */
function setPageContent(state) {
    const parser = new DOMParser();
    const dom = parser.parseFromString(state, "text/html");
    document.body = dom.body;
    document.title = dom.title;
}

/**
 * @param {HTMLElement} element
 * @returns {boolean}
 */
function hasHTMXAttributes(element) {
    for (const attribute of element.attributes) {
        if (attribute.name.startsWith("hx-")) {
            return true;
        }
    }
    return false;
}
