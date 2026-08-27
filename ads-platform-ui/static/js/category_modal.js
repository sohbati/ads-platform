(function () {
  const modal = document.getElementById("category-modal");
  if (!modal) return;

  const navEl = document.getElementById("category-modal-nav");
  const panelEl = document.getElementById("category-modal-panel");
  const statusEl = document.getElementById("category-modal-status");
  const TONE_COUNT = 4;

  const ICONS = {
    home: '<path d="M4 10.5 12 4l8 6.5V20a1 1 0 0 1-1 1h-5v-6H10v6H5a1 1 0 0 1-1-1v-9.5Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    building: '<path d="M4 20V6l6-2 6 2v14M10 20v-6h4v6M8 9h.01M8 12h.01M16 9h.01M16 12h.01M4 20h16" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
    key: '<path d="M8 14a4 4 0 1 1 3.5-6H20l2 2-2 2h-2l-1 1v2h-2v2h-3.5A4 4 0 0 1 8 14Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    calendar: '<rect x="4" y="5" width="16" height="15" rx="2" stroke="currentColor" stroke-width="1.5"/><path d="M8 3v4M16 3v4M4 10h16" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    hardhat: '<path d="M4 14h16v3H4v-3Zm1-1a7 7 0 0 1 14 0M12 5v3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
    car: '<path d="M5 17h14M6 17l-1-4 2-5h10l2 5-1 4M8 17a1.5 1.5 0 1 0 0 .01M16 17a1.5 1.5 0 1 0 0 .01M7 13h10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
    motorcycle: '<circle cx="6.5" cy="16.5" r="2.5" stroke="currentColor" stroke-width="1.5"/><circle cx="17.5" cy="16.5" r="2.5" stroke="currentColor" stroke-width="1.5"/><path d="M9 16.5h5L12 10h4l2 4M8 11h3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
    cog: '<circle cx="12" cy="12" r="3" stroke="currentColor" stroke-width="1.5"/><path d="M12 4v2M12 18v2M4 12h2M18 12h2M6.3 6.3l1.4 1.4M16.3 16.3l1.4 1.4M17.7 6.3l-1.4 1.4M7.7 16.3l-1.4 1.4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    boat: '<path d="M4 14s2 4 8 4 8-4 8-4H4Zm4-3 4-6 4 6" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    device: '<rect x="7" y="3" width="10" height="18" rx="2" stroke="currentColor" stroke-width="1.5"/><path d="M11 18h2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    laptop: '<path d="M5 6h14v9H5V6Zm-2 11h18" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    desktop: '<rect x="4" y="4" width="16" height="11" rx="1.5" stroke="currentColor" stroke-width="1.5"/><path d="M9 20h6M12 15v5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    gamepad: '<rect x="3" y="8" width="18" height="9" rx="3" stroke="currentColor" stroke-width="1.5"/><path d="M8 12.5h3M9.5 11v3M16 12h.01M18 14h.01" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    tv: '<rect x="3" y="5" width="18" height="12" rx="2" stroke="currentColor" stroke-width="1.5"/><path d="M8 21h8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    camera: '<path d="M4 8h3l2-2h6l2 2h3v11H4V8Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/><circle cx="12" cy="13" r="3" stroke="currentColor" stroke-width="1.5"/>',
    phone: '<rect x="8" y="3" width="8" height="18" rx="2" stroke="currentColor" stroke-width="1.5"/><path d="M11 18h2" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    sofa: '<path d="M4 14V11a3 3 0 0 1 3-3h10a3 3 0 0 1 3 3v3M4 14v2a2 2 0 0 0 2 2h1M20 14v2a2 2 0 0 1-2 2h-1M4 14h16" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    blender: '<path d="M8 4h8l-1 8H9L8 4Zm1 8v7h6v-7M9 21h6" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    bed: '<path d="M4 18V10h7a4 4 0 0 1 4 4v4M4 14h16v4M4 18h16" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
    lamp: '<path d="M9 21h6M12 17v4M8 17h8l-1.5-8h-5L8 17Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    fan: '<circle cx="12" cy="12" r="2" stroke="currentColor" stroke-width="1.5"/><path d="M12 4c3 2 4 5 2.5 6.5C13 9 12 7 12 4Zm8 8c-2 3-5 4-6.5 2.5C15 13 17 12 20 12ZM12 20c-3-2-4-5-2.5-6.5C11 15 12 17 12 20ZM4 12c2-3 5-4 6.5-2.5C9 11 7 12 4 12Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    wrench: '<path d="M14.7 6.3a4 4 0 0 0-5.4 5.4L4 17l3 3 5.3-5.3a4 4 0 0 0 5.4-5.4l-2.2 2.2-3.1-3.1 2.2-2.2Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    book: '<path d="M5 5h9a3 3 0 0 1 3 3v12H8a3 3 0 0 0-3 3V5Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/><path d="M5 5a3 3 0 0 1 3-3h9" stroke="currentColor" stroke-width="1.5"/>',
    truck: '<path d="M3 16V8h11v8M14 11h4l3 3v2h-7" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/><circle cx="6.5" cy="16.5" r="1.5" stroke="currentColor" stroke-width="1.5"/><circle cx="16.5" cy="16.5" r="1.5" stroke="currentColor" stroke-width="1.5"/>',
    spark: '<path d="M12 3v4M12 17v4M5 7l2.5 2.5M16.5 14.5 19 17M5 17l2.5-2.5M16.5 9.5 19 7" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/><circle cx="12" cy="12" r="2.5" stroke="currentColor" stroke-width="1.5"/>',
    broom: '<path d="M8 20 18 6M6 14l8 6 2-8-8-4-2 6Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    briefcase: '<rect x="3" y="8" width="18" height="12" rx="2" stroke="currentColor" stroke-width="1.5"/><path d="M8 8V6a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" stroke="currentColor" stroke-width="1.5"/>',
    chart: '<path d="M4 19h16M7 16V10M12 16V7M17 16v-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    code: '<path d="m8 9-4 3 4 3M16 9l4 3-4 3M13 6l-2 12" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>',
    coin: '<circle cx="12" cy="12" r="8" stroke="currentColor" stroke-width="1.5"/><path d="M12 8v8M9.5 10.5c.5-1 1.5-1.5 2.5-1.5s2 .6 2.5 1.5M9.5 13.5c.5 1 1.5 1.5 2.5 1.5s2-.6 2.5-1.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    shirt: '<path d="M16 4l2 3-2 2v11H8V9L6 7l2-3 4 2 4-2Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    gem: '<path d="M8 4h8l4 5-8 11L4 9l4-5Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/><path d="M4 9h16M10 4 8 9l4 11 4-11-2-5" stroke="currentColor" stroke-width="1.5"/>',
    bottle: '<path d="M10 6V4h4v2l1 2v11a2 2 0 0 1-2 2h-2a2 2 0 0 1-2-2V8l1-2Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    ball: '<circle cx="12" cy="12" r="9" stroke="currentColor" stroke-width="1.5"/><path d="M12 3c2.5 3 2.5 15 0 18M3 12h18" stroke="currentColor" stroke-width="1.5"/>',
    bike: '<circle cx="6.5" cy="16" r="3" stroke="currentColor" stroke-width="1.5"/><circle cx="17.5" cy="16" r="3" stroke="currentColor" stroke-width="1.5"/><path d="M6.5 16 10 9h5l2.5 7M10 9l2 7" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    ticket: '<path d="M4 8a2 2 0 0 0 2-2h12a2 2 0 0 0 2 2v2a2 2 0 0 1 0 4v2a2 2 0 0 0-2 2H6a2 2 0 0 0-2-2v-2a2 2 0 0 1 0-4V8Z" stroke="currentColor" stroke-width="1.5"/><path d="M12 7v10" stroke="currentColor" stroke-width="1.5" stroke-dasharray="2 2"/>',
    users: '<circle cx="9" cy="8" r="3" stroke="currentColor" stroke-width="1.5"/><circle cx="16" cy="9" r="2.5" stroke="currentColor" stroke-width="1.5"/><path d="M3 19c1.2-3 3.5-4.5 6-4.5s4.8 1.5 6 4.5M14 14.5c1.5-.3 3 .4 4 2.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>',
    factory: '<path d="M3 20h18M4 20V10l5 3V10l5 3V8h6v12" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    tools: '<path d="M7 20 3 16l6-6 2 2M14 4l6 6-3 3-6-6 3-3Z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>',
    grid: '<rect x="4" y="4" width="7" height="7" rx="1.5" stroke="currentColor" stroke-width="1.5"/><rect x="13" y="4" width="7" height="7" rx="1.5" stroke="currentColor" stroke-width="1.5"/><rect x="4" y="13" width="7" height="7" rx="1.5" stroke="currentColor" stroke-width="1.5"/><rect x="13" y="13" width="7" height="7" rx="1.5" stroke="currentColor" stroke-width="1.5"/>',
  };

  const ICON_BY_SLUG = {
    "real-estate": "home",
    "residential-sale": "home",
    "residential-rent": "key",
    "commercial-sale": "building",
    "commercial-rent": "building",
    "short-term-rent": "calendar",
    construction: "hardhat",
    vehicles: "car",
    cars: "car",
    motorcycles: "motorcycle",
    "auto-parts": "cog",
    boats: "boat",
    "car-rental": "key",
    digital: "device",
    "mobile-tablet": "device",
    laptop: "laptop",
    computers: "desktop",
    gaming: "gamepad",
    av: "tv",
    cameras: "camera",
    phones: "phone",
    home: "sofa",
    appliances: "blender",
    kitchen: "blender",
    "bedroom-bath": "bed",
    decor: "lamp",
    hvac: "fan",
    services: "wrench",
    education: "book",
    transport: "truck",
    beauty: "spark",
    cleaning: "broom",
    jobs: "briefcase",
    "admin-jobs": "briefcase",
    "sales-jobs": "chart",
    "it-jobs": "code",
    "finance-jobs": "coin",
    personal: "shirt",
    fashion: "shirt",
    jewelry: "gem",
    cosmetics: "bottle",
    leisure: "ball",
    sports: "ball",
    books: "book",
    bikes: "bike",
    events: "ticket",
    social: "users",
    volunteer: "users",
    industrial: "factory",
    "industrial-tools": "tools",
    machinery: "factory",
  };

  let groups = [];
  let loadPromise = null;
  let lastFocused = null;
  let activeId = 0;

  function t(key, fallback) {
    return modal.getAttribute(key) || fallback || "";
  }

  function iconSvg(slug) {
    const key = ICON_BY_SLUG[slug] || "grid";
    return '<svg class="cat-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true">' + (ICONS[key] || ICONS.grid) + "</svg>";
  }

  function searchHref(slug) {
    return "/query-ads?category=" + encodeURIComponent(slug);
  }

  function formatViewAll(title) {
    return String(t("data-i18n-view-all", "%s")).replace("%s", title);
  }

  function setStatus(text) {
    if (statusEl) statusEl.textContent = text || "";
  }

  function buildGroups(items) {
    const childrenByParent = {};
    const roots = [];
    (items || []).forEach(function (item) {
      if (item.parent == null) {
        roots.push(item);
        return;
      }
      const pid = item.parent;
      if (!childrenByParent[pid]) childrenByParent[pid] = [];
      childrenByParent[pid].push(item);
    });
    roots.sort(function (a, b) { return (a.order || 0) - (b.order || 0); });
    return roots.map(function (root, index) {
      const children = (childrenByParent[root.id] || []).slice().sort(function (a, b) {
        return (a.order || 0) - (b.order || 0);
      });
      return {
        id: root.id,
        title: root.title,
        slug: root.slug,
        tone: String(index % TONE_COUNT),
        children: children,
      };
    });
  }

  function render() {
    if (!navEl || !panelEl) return;
    navEl.innerHTML = groups.map(function (group) {
      const on = group.id === activeId;
      return (
        '<li><button type="button" class="cat-modal__nav-item' + (on ? " is-active" : "") + '" data-tone="' + group.tone + '" data-cat-id="' + group.id + '" aria-pressed="' + (on ? "true" : "false") + '">' +
          '<span class="cat-modal__nav-icon">' + iconSvg(group.slug) + "</span>" +
          "<span>" + escapeHtml(group.title) + "</span>" +
        "</button></li>"
      );
    }).join("");

    panelEl.innerHTML = groups.map(function (group) {
      const on = group.id === activeId;
      let cards = "";
      if (group.children.length) {
        cards = '<div class="cat-modal__cards">' + group.children.map(function (child) {
          return (
            '<a class="cat-modal__card" data-tone="' + group.tone + '" href="' + searchHref(child.slug) + '">' +
              '<span class="cat-modal__card-icon">' + iconSvg(child.slug) + "</span>" +
              '<span class="cat-modal__card-title">' + escapeHtml(child.title) + "</span>" +
            "</a>"
          );
        }).join("") + "</div>";
      } else {
        cards = '<p class="cat-modal__empty">' + escapeHtml(t("data-i18n-empty", "")) + "</p>";
      }
      return (
        '<div class="cat-modal__panel-body" data-tone="' + group.tone + '" data-cat-panel="' + group.id + '"' + (on ? "" : " hidden") + ">" +
          '<a class="cat-modal__view-all" href="' + searchHref(group.slug) + '">' +
            escapeHtml(formatViewAll(group.title)) +
            '<svg class="cat-modal__view-all-icon" viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="m9 6 6 6-6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/></svg>' +
          "</a>" +
          cards +
        "</div>"
      );
    }).join("");
  }

  function escapeHtml(value) {
    return String(value || "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function activate(id) {
    const next = Number(id);
    if (!next) return;
    activeId = next;
    navEl.querySelectorAll("[data-cat-id]").forEach(function (btn) {
      const on = Number(btn.getAttribute("data-cat-id")) === next;
      btn.classList.toggle("is-active", on);
      btn.setAttribute("aria-pressed", on ? "true" : "false");
    });
    panelEl.querySelectorAll("[data-cat-panel]").forEach(function (panel) {
      const match = Number(panel.getAttribute("data-cat-panel")) === next;
      panel.hidden = !match;
    });
  }

  function loadCategories() {
    if (groups.length) return Promise.resolve(groups);
    if (loadPromise) return loadPromise;
    const url = modal.getAttribute("data-categories-url") || "/api/v1/categories";
    loadPromise = fetch(url, { credentials: "same-origin" })
      .then(function (resp) {
        if (!resp.ok) throw new Error("categories");
        return resp.json();
      })
      .then(function (items) {
        groups = buildGroups(items);
        if (!activeId && groups.length) activeId = groups[0].id;
        render();
        return groups;
      })
      .catch(function (err) {
        loadPromise = null;
        throw err;
      });
    return loadPromise;
  }

  function openModal(slug) {
    lastFocused = document.activeElement;
    setStatus("");
    modal.hidden = false;
    modal.classList.add("is-open");
    modal.setAttribute("aria-hidden", "false");
    document.body.classList.add("category-modal-open");

    loadCategories()
      .then(function () {
        if (slug) {
          const match = groups.find(function (g) { return g.slug === slug; });
          if (match) activeId = match.id;
        }
        if (!activeId && groups.length) activeId = groups[0].id;
        render();
      })
      .catch(function () {
        setStatus(t("data-i18n-load-error"));
      });
  }

  function closeModal() {
    modal.classList.remove("is-open");
    modal.hidden = true;
    modal.setAttribute("aria-hidden", "true");
    document.body.classList.remove("category-modal-open");
    if (lastFocused && typeof lastFocused.focus === "function") lastFocused.focus();
  }

  navEl.addEventListener("click", function (ev) {
    const btn = ev.target.closest("[data-cat-id]");
    if (!btn || !navEl.contains(btn)) return;
    ev.preventDefault();
    activate(btn.getAttribute("data-cat-id"));
  });

  navEl.addEventListener("mouseover", function (ev) {
    const btn = ev.target.closest("[data-cat-id]");
    if (!btn || !navEl.contains(btn)) return;
    activate(btn.getAttribute("data-cat-id"));
  });

  modal.querySelectorAll("[data-cat-close]").forEach(function (el) {
    el.addEventListener("click", closeModal);
  });

  document.addEventListener("keydown", function (ev) {
    if (ev.key === "Escape" && modal.classList.contains("is-open")) {
      ev.preventDefault();
      closeModal();
    }
  });

  document.querySelectorAll("[data-category-open]").forEach(function (el) {
    el.addEventListener("click", function (ev) {
      ev.preventDefault();
      openModal(el.getAttribute("data-category-open") || "");
    });
  });

  window.openCategoryModal = openModal;
  window.closeCategoryModal = closeModal;

  const params = new URLSearchParams(window.location.search);
  if (params.get("open") === "category") {
    openModal(params.get("cat") || "");
    params.delete("open");
    params.delete("cat");
    const clean = params.toString();
    window.history.replaceState({}, "", window.location.pathname + (clean ? "?" + clean : "") + window.location.hash);
  }
})();
