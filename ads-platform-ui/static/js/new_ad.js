(function () {
  const form = document.getElementById("new-ad-form");
  if (!form) return;

  const bootstrapEl = document.getElementById("new-ad-bootstrap");
  let boot = {
    locale: "fa",
    cityId: 0,
    maxPictures: 8,
    categories: [],
    schemas: [],
    enums: {},
    successHref: "/my-info/user-ads",
  };
  try {
    boot = Object.assign(boot, JSON.parse(bootstrapEl.textContent || "{}"));
  } catch (_err) {
    // Keep defaults.
  }

  const categorySelect = document.getElementById("new-ad-category");
  const attrsRoot = document.getElementById("new-ad-attrs");
  const messageEl = document.getElementById("new-ad-message");
  const submitBtn = document.getElementById("new-ad-submit");
  const picturesInput = document.getElementById("new-ad-pictures");

  const resolveError = typeof window.resolveApiError === "function"
    ? window.resolveApiError
    : function (_code, _params, fallback) {
        return fallback || "";
      };

  const schemasByName = {};
  (boot.schemas || []).forEach(function (item) {
    if (item && item.name) schemasByName[item.name] = item;
  });

  function setMessage(text, isError) {
    if (!messageEl) return;
    messageEl.textContent = text || "";
    messageEl.classList.toggle("is-error", Boolean(isError));
  }

  function asciiDigits(value) {
    return String(value || "")
      .replace(/[۰-۹]/g, function (d) { return String(d.charCodeAt(0) - 1776); })
      .replace(/[٠-٩]/g, function (d) { return String(d.charCodeAt(0) - 1632); })
      .replace(/\D/g, "");
  }

  function formatPriceGrouped(digits) {
    if (!digits) return "";
    digits = String(digits).replace(/^0+(?=\d)/, "");
    return digits.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  }

  function parsePriceAmount(value) {
    const digits = asciiDigits(value);
    if (!digits) return null;
    const n = parseInt(digits, 10);
    return Number.isNaN(n) ? null : n;
  }

  function formatPriceInput(el) {
    if (!el) return;
    const caret = el.selectionStart;
    const old = el.value;
    const digitsBefore = asciiDigits(old.slice(0, caret)).length;
    const formatted = formatPriceGrouped(asciiDigits(old));
    if (formatted === old) return;
    el.value = formatted;
    if (document.activeElement !== el || caret == null) return;
    if (digitsBefore === 0) {
      el.setSelectionRange(0, 0);
      return;
    }
    let seen = 0;
    let pos = formatted.length;
    for (let i = 0; i < formatted.length; i++) {
      if (formatted[i] >= "0" && formatted[i] <= "9") {
        seen++;
        if (seen === digitsBefore) {
          pos = i + 1;
          break;
        }
      }
    }
    el.setSelectionRange(pos, pos);
  }

  function bindPriceField() {
    const el = document.getElementById("new-ad-price");
    if (!el) return;
    el.addEventListener("input", function () { formatPriceInput(el); });
  }

  function parentTitle(categories, parentId) {
    if (parentId == null) return "";
    for (let i = 0; i < categories.length; i++) {
      if (categories[i].id === parentId) return categories[i].title || "";
    }
    return "";
  }

  function fillCategories() {
    const categories = boot.categories || [];
    const leaves = categories.filter(function (c) { return c.isLeaf; });
    const groups = {};
    leaves.forEach(function (leaf) {
      const group = parentTitle(categories, leaf.parent) || leaf.title;
      if (!groups[group]) groups[group] = [];
      groups[group].push(leaf);
    });
    Object.keys(groups).sort().forEach(function (group) {
      const optgroup = document.createElement("optgroup");
      optgroup.label = group;
      groups[group]
        .slice()
        .sort(function (a, b) { return (a.order || 0) - (b.order || 0); })
        .forEach(function (leaf) {
          const opt = document.createElement("option");
          opt.value = String(leaf.id);
          opt.textContent = leaf.title;
          opt.dataset.template = leaf.adsAttrsJsonSchemaTemplateName || leaf.slug || "";
          optgroup.appendChild(opt);
        });
      categorySelect.appendChild(optgroup);
    });
  }

  function enumLabel(vocab, token) {
    const loc = boot.locale || "fa";
    const enums = boot.enums || {};
    const byVocab = enums[vocab];
    if (byVocab && byVocab[token]) {
      return byVocab[token][loc] || byVocab[token].fa || byVocab[token].en || token;
    }
    return token;
  }

  function parseSchema(raw) {
    if (!raw) return null;
    if (typeof raw === "string") {
      try { return JSON.parse(raw); } catch (_e) { return null; }
    }
    return raw;
  }

  function selectedSchema() {
    const opt = categorySelect.options[categorySelect.selectedIndex];
    if (!opt || !opt.dataset.template) return null;
    const item = schemasByName[opt.dataset.template];
    if (!item) return null;
    return parseSchema(item.jsonSchema);
  }

  function renderAttrs() {
    attrsRoot.innerHTML = "";
    const schema = selectedSchema();
    if (!schema || !schema.properties) return;
    const required = {};
    (schema.required || []).forEach(function (name) { required[name] = true; });

    Object.keys(schema.properties).forEach(function (name) {
      const prop = schema.properties[name] || {};
      const wrap = document.createElement("div");
      wrap.className = "new-ad__field";
      const id = "attr-" + name;
      const label = document.createElement("label");
      label.className = "new-ad__label";
      label.setAttribute("for", id);
      label.textContent = (prop.title || name) + (required[name] ? " *" : "");
      wrap.appendChild(label);

      if (prop.type === "boolean") {
        wrap.className = "new-ad__check";
        const input = document.createElement("input");
        input.type = "checkbox";
        input.id = id;
        input.dataset.attr = name;
        input.dataset.type = "boolean";
        wrap.insertBefore(input, label);
        label.removeAttribute("class");
      } else if (Array.isArray(prop.enum)) {
        const select = document.createElement("select");
        select.id = id;
        select.dataset.attr = name;
        select.dataset.type = "string";
        if (required[name]) select.required = true;
        const empty = document.createElement("option");
        empty.value = "";
        empty.textContent = "—";
        select.appendChild(empty);
        const vocab = prop["x-enumVocab"];
        prop.enum.forEach(function (token) {
          const opt = document.createElement("option");
          opt.value = String(token);
          opt.textContent = vocab ? enumLabel(vocab, token) : String(token);
          select.appendChild(opt);
        });
        wrap.appendChild(select);
      } else {
        const input = document.createElement("input");
        input.id = id;
        input.dataset.attr = name;
        if (prop.type === "integer" || prop.type === "number") {
          input.type = "number";
          input.dataset.type = prop.type;
          if (prop.minimum != null) input.min = prop.minimum;
          if (prop.maximum != null) input.max = prop.maximum;
          if (prop.type === "integer") input.step = "1";
        } else {
          input.type = "text";
          input.dataset.type = "string";
          if (prop.maxLength) input.maxLength = prop.maxLength;
        }
        if (required[name]) input.required = true;
        wrap.appendChild(input);
      }
      attrsRoot.appendChild(wrap);
    });
  }

  function collectAttrs(schema) {
    const attrs = {};
    if (!schema || !schema.properties) return attrs;
    Object.keys(schema.properties).forEach(function (name) {
      const el = document.getElementById("attr-" + name);
      if (!el) return;
      const type = el.dataset.type;
      if (type === "boolean") {
        attrs[name] = el.checked;
        return;
      }
      const value = (el.value || "").trim();
      if (value === "") return;
      if (type === "integer") {
        const n = parseInt(value, 10);
        if (!Number.isNaN(n)) attrs[name] = n;
      } else if (type === "number") {
        const n = parseFloat(value);
        if (!Number.isNaN(n)) attrs[name] = n;
      } else {
        attrs[name] = value;
      }
    });
    return attrs;
  }

  async function readError(resp) {
    const text = await resp.text();
    try {
      const data = JSON.parse(text);
      if (data && typeof data.error === "string") {
        return {
          code: data.error,
          params: Array.isArray(data.params) ? data.params : [],
        };
      }
    } catch (_err) {}
    return { code: "", params: [] };
  }

  fillCategories();
  categorySelect.addEventListener("change", renderAttrs);
  applyPrefill();
  bindPriceField();

  form.addEventListener("submit", async function (event) {
    event.preventDefault();
    setMessage("");

    if (!boot.cityId) {
      setMessage(form.dataset.cityRequired || "", true);
      return;
    }
    const categoryId = parseInt(categorySelect.value, 10);
    if (!categoryId) {
      setMessage(form.dataset.needCategory || "", true);
      return;
    }
    const title = (document.getElementById("new-ad-title").value || "").trim();
    const description = (document.getElementById("new-ad-description").value || "").trim();
    if (!title) {
      setMessage(form.dataset.needTitle || "", true);
      return;
    }
    if (!description) {
      setMessage(form.dataset.needDescription || "", true);
      return;
    }

    const files = picturesInput && picturesInput.files ? Array.from(picturesInput.files) : [];
    const maxPics = boot.maxPictures || 8;
    if (files.length > maxPics) {
      setMessage(form.dataset.tooMany || "", true);
      return;
    }

    const isEdit = boot.mode === "edit" && boot.adId;
    const submitPath = isEdit ? "/api/v1/ads/" + boot.adId : "/api/v1/ads";
    const method = isEdit ? "PUT" : "POST";
    const loginNext = isEdit ? "/edit-ad/" + boot.adId : "/new-ad";

    const priceAmount = parsePriceAmount(document.getElementById("new-ad-price").value);
    const payload = {
      category_id: categoryId,
      city_id: boot.cityId,
      title: title,
      description: description,
      price_type: document.getElementById("new-ad-price-type").value || "fixed",
      currency: "IRR",
      neighborhood: (document.getElementById("new-ad-neighborhood").value || "").trim(),
      attrs: collectAttrs(selectedSchema()),
    };
    if (priceAmount != null) {
      payload.price_amount = priceAmount;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = form.dataset.submitting || submitBtn.textContent;

    try {
      let resp;
      if (files.length) {
        const fd = new FormData();
        fd.append("payload", JSON.stringify(payload));
        files.forEach(function (file) { fd.append("pictures", file); });
        resp = await fetch(submitPath, { method: method, credentials: "same-origin", body: fd });
      } else {
        resp = await fetch(submitPath, {
          method: method,
          credentials: "same-origin",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(payload),
        });
      }

      if (resp.status === 401) {
        if (typeof window.openLoginModal === "function") {
          window.openLoginModal(loginNext);
        } else {
          window.location.assign("/query-ads?login=1&next=" + encodeURIComponent(loginNext));
        }
        return;
      }

      if (!resp.ok) {
        const err = await readError(resp);
        setMessage(resolveError(err.code, err.params, resolveError("_default", [], err.code)), true);
        return;
      }

      window.location.assign(boot.successHref || "/my-info/user-ads");
    } catch (_err) {
      setMessage(resolveError("NETWORK_ERROR", [], ""), true);
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = form.dataset.submitLabel || submitBtn.textContent;
    }
  });

  const changeCityBtn = document.getElementById("new-ad-change-city");
  const cityNameEl = document.getElementById("new-ad-city-name");
  if (changeCityBtn) {
    changeCityBtn.addEventListener("click", function () {
      if (typeof window.openLocationModal !== "function") return;
      window.openLocationModal({
        mode: "single",
        selected: boot.citySlug ? [boot.citySlug] : [],
        onApply: function (items) {
          const city = pickCity(items);
          if (!city || !city.id) {
            setMessage(form.dataset.cityRequired || "", true);
            return;
          }
          boot.cityId = city.id;
          boot.citySlug = city.slug;
          boot.cityName = city.name;
          if (cityNameEl) cityNameEl.textContent = city.name;
        },
      });
    });
  }

  function pickCity(items) {
    if (!items || !items.length) return null;
    for (let i = 0; i < items.length; i++) {
      if (items[i] && items[i].type === "2") return items[i];
    }
    const first = items[0];
    if (first && first.firstCity) return first.firstCity;
    return first;
  }

  function applyPrefill() {
    const p = boot.prefill;
    if (!p) return;
    if (p.category_id) {
      categorySelect.value = String(p.category_id);
      renderAttrs();
    }
    const titleEl = document.getElementById("new-ad-title");
    const descEl = document.getElementById("new-ad-description");
    const priceEl = document.getElementById("new-ad-price");
    const priceTypeEl = document.getElementById("new-ad-price-type");
    const neighborhoodEl = document.getElementById("new-ad-neighborhood");
    if (titleEl && p.title) titleEl.value = p.title;
    if (descEl && p.description) descEl.value = p.description;
    if (priceEl && p.price_amount != null) priceEl.value = formatPriceGrouped(String(p.price_amount));
    if (priceTypeEl && p.price_type) priceTypeEl.value = p.price_type;
    if (neighborhoodEl && p.neighborhood) neighborhoodEl.value = p.neighborhood;
    fillAttrs(p.attrs);
    renderExistingMedia(p.media);
  }

  function fillAttrs(attrs) {
    if (!attrs || typeof attrs !== "object") return;
    Object.keys(attrs).forEach(function (name) {
      const el = document.getElementById("attr-" + name);
      if (!el) return;
      const value = attrs[name];
      if (el.dataset.type === "boolean") {
        el.checked = Boolean(value);
        return;
      }
      if (value != null && value !== "") {
        el.value = String(value);
      }
    });
  }

  function renderExistingMedia(media) {
    const root = document.getElementById("new-ad-existing-pictures");
    if (!root || !media || !media.length) return;
    root.innerHTML = "";
    media.forEach(function (item) {
      const src = (item && (item.thumb || item.url)) || "";
      if (!src) return;
      const img = document.createElement("img");
      img.src = src;
      img.alt = "";
      root.appendChild(img);
    });
    if (root.childNodes.length) root.hidden = false;
  }
})();
