(function () {
  function initMegaCat(root) {
    if (!root || root.dataset.megaReady === "1") return;
    root.dataset.megaReady = "1";

    const navItems = root.querySelectorAll("[data-mega-id]");
    const panels = root.querySelectorAll("[data-mega-panel]");

    function activate(id) {
      const next = String(id);
      root.setAttribute("data-active", next);
      navItems.forEach(function (btn) {
        const on = btn.getAttribute("data-mega-id") === next;
        btn.classList.toggle("is-active", on);
        btn.setAttribute("aria-pressed", on ? "true" : "false");
      });
      panels.forEach(function (panel) {
        const match = panel.getAttribute("data-mega-panel") === next;
        panel.classList.toggle("is-active", match);
        if (match) {
          panel.removeAttribute("hidden");
        } else {
          panel.setAttribute("hidden", "");
        }
      });
    }

    root.addEventListener("click", function (event) {
      const btn = event.target.closest("[data-mega-id]");
      if (!btn || !root.contains(btn)) return;
      event.preventDefault();
      activate(btn.getAttribute("data-mega-id"));
    });

    root.addEventListener("mouseover", function (event) {
      const btn = event.target.closest("[data-mega-id]");
      if (!btn || !root.contains(btn)) return;
      activate(btn.getAttribute("data-mega-id"));
    });

    root.addEventListener("focusin", function (event) {
      const btn = event.target.closest("[data-mega-id]");
      if (!btn || !root.contains(btn)) return;
      activate(btn.getAttribute("data-mega-id"));
    });

    const initial = root.getAttribute("data-active");
    if (initial) activate(initial);
  }

  document.querySelectorAll("[data-mega-cat]").forEach(initMegaCat);
})();
