(function () {
  const modal = document.getElementById("location-modal");
  if (!modal) return;

  const TYPE_PROVINCE = "1";
  const TYPE_CITY = "2";
  const COOKIE_CITY = "city";
  const COOKIE_LOCATIONS = "locations";
  const COOKIE_MAX_AGE = 365 * 24 * 3600;

  const searchInput = document.getElementById("location-modal-search");
  const chipsEl = document.getElementById("location-modal-chips");
  const treeEl = document.getElementById("location-modal-tree");
  const badgeEl = document.getElementById("location-modal-badge");
  const countEl = document.getElementById("location-modal-count");
  const namesEl = document.getElementById("location-modal-names");
  const statusEl = document.getElementById("location-modal-status");
  const applyBtn = document.getElementById("location-modal-apply");
  const clearBtn = document.getElementById("location-modal-clear");
  const subtitleEl = modal.querySelector("[data-loc-subtitle-multi]");

  let catalog = null;
  let loadPromise = null;
  let mode = "multiple";
  let selected = new Set();
  let expanded = new Set();
  let query = "";
  let onApply = null;
  let lastFocused = null;
  let treeRenderTimer = 0;

  function t(key, fallback) {
    return modal.getAttribute(key) || fallback || "";
  }

  function toLocaleDigits(value) {
    const s = String(value);
    if (document.documentElement.lang !== "fa") return s;
    return s.replace(/\d/g, function (d) { return "۰۱۲۳۴۵۶۷۸۹"[d]; });
  }

  function formatCount(template, n) {
    return toLocaleDigits(String(template || "").replace("%d", String(n)));
  }

  function setCookie(name, value) {
    document.cookie = name + "=" + encodeURIComponent(value) + "; Path=/; Max-Age=" + COOKIE_MAX_AGE + "; SameSite=Lax";
  }

  function parseCSV(raw) {
    return String(raw || "")
      .split(",")
      .map(function (s) { return s.trim(); })
      .filter(Boolean);
  }

  function setStatus(text) {
    if (statusEl) statusEl.textContent = text || "";
  }

  function byId(id) {
    return catalog && catalog.byId[id];
  }

  function displayName(item) {
    if (!item) return "";
    if (item.type === TYPE_CITY && item.parent != null) {
      const parent = byId(item.parent);
      if (parent && parent.name === item.name) {
        const suffix = t("data-i18n-city-suffix", "");
        return suffix ? item.name + " (" + suffix + ")" : item.name;
      }
    }
    return item.name || item.slug || "";
  }

  function selectionLabel(item) {
    if (!item) return "";
    if (item.type === TYPE_CITY && item.parent != null) {
      const parent = byId(item.parent);
      if (parent && parent.name !== item.name) {
        return parent.name + " (" + item.name + ")";
      }
    }
    return displayName(item);
  }

  function buildCatalog(rows) {
    const byIdMap = {};
    const bySlug = {};
    const children = {};
    const provinces = [];
    (rows || []).forEach(function (row) {
      if (!row || !row.id) return;
      byIdMap[row.id] = row;
      if (row.slug) bySlug[String(row.slug).toLowerCase()] = row;
      if (row.parent != null) {
        if (!children[row.parent]) children[row.parent] = [];
        children[row.parent].push(row);
      }
      if (row.type === TYPE_PROVINCE) provinces.push(row);
    });
    return { byId: byIdMap, bySlug: bySlug, children: children, provinces: provinces };
  }

  function loadCities() {
    if (catalog) return Promise.resolve(catalog);
    if (loadPromise) return loadPromise;
    const url = modal.getAttribute("data-cities-url") || "/api/v1/cities";
    loadPromise = fetch(url, { credentials: "same-origin" })
      .then(function (resp) {
        if (!resp.ok) throw new Error("cities");
        return resp.json();
      })
      .then(function (rows) {
        catalog = buildCatalog(rows);
        return catalog;
      })
      .catch(function (err) {
        loadPromise = null;
        throw err;
      });
    return loadPromise;
  }

  function popularItems() {
    if (!catalog) return [];
    return parseCSV(modal.getAttribute("data-popular")).map(function (slug) {
      return catalog.bySlug[slug];
    }).filter(Boolean);
  }

  function cityChildren(provinceId) {
    return (catalog.children[provinceId] || []).filter(function (c) { return c.type === TYPE_CITY; });
  }

  function matchesQuery(item, q) {
    if (!q) return true;
    const name = (item.name || "").toLowerCase();
    const slug = (item.slug || "").toLowerCase();
    return name.indexOf(q) !== -1 || slug.indexOf(q) !== -1;
  }

  function visibleTree() {
    const q = query.trim().toLowerCase();
    const out = [];
    catalog.provinces.forEach(function (prov) {
      const cities = cityChildren(prov.id);
      const provMatch = matchesQuery(prov, q);
      const matchedCities = q
        ? cities.filter(function (city) { return matchesQuery(city, q); })
        : cities;
      if (!q) {
        out.push({ province: prov, cities: cities, forceOpen: false });
        return;
      }
      if (provMatch || matchedCities.length) {
        out.push({
          province: prov,
          cities: provMatch && matchedCities.length === 0 ? cities : matchedCities,
          forceOpen: true,
        });
      }
    });
    return out;
  }

  function chevronSvg() {
    return '<svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M9 6l6 6-6 6" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"/></svg>';
  }

  function checkRow(item, extraClass) {
    const row = document.createElement("div");
    row.className = "loc-tree__row" + (extraClass ? " " + extraClass : "");
    row.setAttribute("role", "treeitem");
    row.dataset.slug = item.slug || "";

    const kids = item.type === TYPE_PROVINCE ? cityChildren(item.id) : [];
    const toggle = document.createElement("button");
    toggle.type = "button";
    toggle.className = "loc-tree__toggle" + (kids.length ? "" : " is-leaf");
    toggle.tabIndex = kids.length ? 0 : -1;
    toggle.dataset.id = String(item.id);
    toggle.setAttribute("aria-label", displayName(item));
    toggle.innerHTML = chevronSvg();
    if (kids.length) {
      const open = expanded.has(item.id);
      toggle.classList.toggle("is-open", open);
      toggle.setAttribute("aria-expanded", open ? "true" : "false");
    }
    row.appendChild(toggle);

    const label = document.createElement("label");
    label.className = "loc-check";
    const input = document.createElement("input");
    input.type = "checkbox";
    input.tabIndex = 0;
    input.checked = selected.has(item.slug);
    const box = document.createElement("span");
    box.className = "loc-check__box";
    const text = document.createElement("span");
    text.className = "loc-check__label";
    text.textContent = displayName(item);
    label.appendChild(input);
    label.appendChild(box);
    label.appendChild(text);
    row.appendChild(label);
    return row;
  }

  function renderChips() {
    const frag = document.createDocumentFragment();
    popularItems().forEach(function (item) {
      const btn = document.createElement("button");
      btn.type = "button";
      btn.className = "loc-modal__chip" + (selected.has(item.slug) ? " is-active" : "");
      btn.dataset.slug = item.slug;
      btn.textContent = item.name;
      frag.appendChild(btn);
    });
    chipsEl.replaceChildren(frag);
  }

  function renderTree() {
    const frag = document.createDocumentFragment();
    const groups = visibleTree();
    if (!groups.length) {
      const empty = document.createElement("p");
      empty.className = "loc-tree__row";
      empty.textContent = t("data-i18n-empty-search", "");
      frag.appendChild(empty);
      treeEl.replaceChildren(frag);
      return;
    }
    groups.forEach(function (group) {
      const open = group.forceOpen || expanded.has(group.province.id);
      if (group.forceOpen) expanded.add(group.province.id);
      frag.appendChild(checkRow(group.province, ""));
      if (open) {
        group.cities.forEach(function (city) {
          frag.appendChild(checkRow(city, "is-child"));
        });
      }
    });
    treeEl.replaceChildren(frag);
  }

  function scheduleTreeRender() {
    if (treeRenderTimer) window.cancelAnimationFrame(treeRenderTimer);
    treeRenderTimer = window.requestAnimationFrame(function () {
      treeRenderTimer = 0;
      renderTree();
    });
  }

  function selectedItems() {
    if (!catalog) return [];
    const items = [];
    selected.forEach(function (slug) {
      const item = catalog.bySlug[slug];
      if (item) items.push(item);
    });
    return items;
  }

  function firstCityOf(item) {
    if (!item) return null;
    if (item.type === TYPE_CITY) return item;
    if (item.type === TYPE_PROVINCE) {
      const kids = cityChildren(item.id);
      return kids[0] || null;
    }
    return null;
  }

  function updateSummary() {
    const items = selectedItems();
    const n = items.length;
    badgeEl.textContent = toLocaleDigits(n);
    countEl.textContent = formatCount(t("data-i18n-selected-count"), n);
    namesEl.textContent = items.map(selectionLabel).join("، ");
    applyBtn.textContent = formatCount(t("data-i18n-apply"), n);
  }

  function syncSelectionUI() {
    treeEl.querySelectorAll("[data-slug]").forEach(function (row) {
      const input = row.querySelector("input[type='checkbox']");
      if (input) input.checked = selected.has(row.dataset.slug);
    });
    chipsEl.querySelectorAll("[data-slug]").forEach(function (chip) {
      chip.classList.toggle("is-active", selected.has(chip.dataset.slug));
    });
    updateSummary();
  }

  function render() {
    renderChips();
    renderTree();
    updateSummary();
  }

  function toggleSelect(slug) {
    if (!slug) return;
    if (mode === "single") {
      selected = selected.has(slug) ? new Set() : new Set([slug]);
    } else if (selected.has(slug)) {
      selected.delete(slug);
    } else {
      selected.add(slug);
    }
    syncSelectionUI();
  }

  function expandSelectedParents() {
    selected.forEach(function (slug) {
      const item = catalog.bySlug[slug];
      if (!item) return;
      if (item.type === TYPE_PROVINCE) expanded.add(item.id);
      else if (item.parent != null) expanded.add(item.parent);
    });
  }

  function applyDefault(items) {
    const slugs = items.map(function (item) { return item.slug; });
    const primary = primaryCitySlug(items);
    setCookie(COOKIE_LOCATIONS, slugs.join(","));
    setCookie(COOKIE_CITY, primary || "tehran");
    const reload = function () { window.location.reload(); };
    if (window.locationProfile && typeof window.locationProfile.save === "function") {
      window.locationProfile.save(slugs).then(reload, reload);
      return;
    }
    reload();
  }

  function primaryCitySlug(items) {
    for (let i = 0; i < items.length; i++) {
      const city = firstCityOf(items[i]);
      if (city && city.slug) return city.slug;
    }
    return "";
  }

  function openModal(options) {
    options = options || {};
    if (document.body.classList.contains("category-modal-open") && typeof window.closeCategoryModal === "function") {
      window.closeCategoryModal();
    }
    if (document.body.classList.contains("login-modal-open") && typeof window.closeLoginModal === "function") {
      window.closeLoginModal();
    }
    mode = options.mode === "single" ? "single" : "multiple";
    onApply = typeof options.onApply === "function" ? options.onApply : null;
    modal.dataset.mode = mode;
    lastFocused = document.activeElement;
    setStatus("");
    query = "";
    if (searchInput) searchInput.value = "";
    if (subtitleEl) {
      subtitleEl.textContent = mode === "single"
        ? subtitleEl.getAttribute("data-loc-subtitle-single") || ""
        : subtitleEl.getAttribute("data-loc-subtitle-multi") || "";
    }

    const initial = Array.isArray(options.selected)
      ? options.selected
      : parseCSV(modal.getAttribute("data-selected"));
    selected = new Set(initial.filter(Boolean));
    expanded = new Set();

    modal.hidden = false;
    modal.classList.add("is-open");
    modal.setAttribute("aria-hidden", "false");
    document.body.classList.add("location-modal-open");

    loadCities()
      .then(function () {
        expandSelectedParents();
        render();
        if (searchInput) searchInput.focus();
      })
      .catch(function () {
        setStatus(t("data-i18n-load-error"));
      });
  }

  function closeModal() {
    modal.classList.remove("is-open");
    modal.hidden = true;
    modal.setAttribute("aria-hidden", "true");
    document.body.classList.remove("location-modal-open");
    onApply = null;
    if (lastFocused && typeof lastFocused.focus === "function") lastFocused.focus();
  }

  chipsEl.addEventListener("click", function (ev) {
    const chip = ev.target.closest("[data-slug]");
    if (!chip || !chipsEl.contains(chip)) return;
    ev.preventDefault();
    const item = catalog && catalog.bySlug[chip.dataset.slug];
    if (item && item.parent != null) expanded.add(item.parent);
    toggleSelect(chip.dataset.slug);
    scheduleTreeRender();
  });

  treeEl.addEventListener("click", function (ev) {
    const toggle = ev.target.closest(".loc-tree__toggle");
    if (toggle && treeEl.contains(toggle) && !toggle.classList.contains("is-leaf")) {
      ev.preventDefault();
      ev.stopPropagation();
      const id = Number(toggle.dataset.id);
      if (!id) return;
      if (expanded.has(id)) expanded.delete(id);
      else expanded.add(id);
      scheduleTreeRender();
      return;
    }

    const check = ev.target.closest(".loc-check");
    if (!check || !treeEl.contains(check)) return;
    ev.preventDefault();
    const row = check.closest("[data-slug]");
    if (row && row.dataset.slug) toggleSelect(row.dataset.slug);
  });

  modal.querySelectorAll("[data-loc-close]").forEach(function (el) {
    el.addEventListener("click", closeModal);
  });

  document.addEventListener("keydown", function (ev) {
    if (ev.key === "Escape" && modal.classList.contains("is-open")) {
      ev.preventDefault();
      closeModal();
    }
  });

  if (searchInput) {
    searchInput.addEventListener("input", function () {
      query = searchInput.value || "";
      if (catalog) scheduleTreeRender();
    });
  }

  if (clearBtn) {
    clearBtn.addEventListener("click", function () {
      selected = new Set();
      syncSelectionUI();
    });
  }

  if (applyBtn) {
    applyBtn.addEventListener("click", function () {
      const items = selectedItems();
      if (onApply) {
        onApply(items.map(function (item) {
          return Object.assign({}, item, { firstCity: firstCityOf(item) });
        }), mode);
        closeModal();
        return;
      }
      applyDefault(items);
    });
  }

  document.querySelectorAll("[data-location-open]").forEach(function (el) {
    el.addEventListener("click", function (ev) {
      ev.preventDefault();
      const openMode = el.getAttribute("data-location-open") || "multiple";
      const selectedAttr = el.getAttribute("data-location-selected");
      openModal({
        mode: openMode,
        selected: selectedAttr != null ? parseCSV(selectedAttr) : undefined,
      });
    });
  });

  window.openLocationModal = openModal;
  window.closeLocationModal = closeModal;
})();
