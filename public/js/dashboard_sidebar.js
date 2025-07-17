//@ts-check

document.addEventListener("DOMContentLoaded", main);

function main() {
    const toggleButton = document.querySelector("[data-sidebar-toggle]");
    const sidebar = document.querySelector("[data-sidebar]");
    if (!toggleButton || !sidebar) return;

    toggleButton.addEventListener("click", (e) => {
        sidebar.classList.toggle("hidden");
    });
}
