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
    picturesUnavailable: false,
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
  const photoGrid = document.getElementById("new-ad-photo-grid");
  const addPhotosBtn = document.getElementById("new-ad-add-photos");
  const picturesHint = document.getElementById("new-ad-pictures-hint");
  const picturesUnavailableHint = document.getElementById("new-ad-pictures-unavailable");
  const storageContinue = document.getElementById("new-ad-storage-continue");
  const continueWithoutPhotosBtn = document.getElementById("new-ad-continue-without-photos");
  const lightbox = document.getElementById("new-ad-lightbox");
  const lightboxImg = lightbox ? lightbox.querySelector("[data-photo-full]") : null;
  const lightboxCount = lightbox ? lightbox.querySelector("[data-photo-count]") : null;
  const lightboxPrev = lightbox ? lightbox.querySelector("[data-photo-prev]") : null;
  const lightboxNext = lightbox ? lightbox.querySelector("[data-photo-next]") : null;

  const pendingPhotos = [];
  let existingPhotos = [];
  let lightboxIndex = 0;
  let pinLat = null;
  let pinLng = null;
  let adMap = null;
  let adMarker = null;
  let userEditedNeighborhood = false;
  let reverseTimer = 0;
  let publishWithoutPhotos = false;

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
          opt.dataset.slug = leaf.slug || "";
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

  function updateTitlePlaceholder() {
    const titleEl = document.getElementById("new-ad-title");
    if (!titleEl) return;
    const opt = categorySelect.options[categorySelect.selectedIndex];
    const slug = opt && opt.dataset.slug ? opt.dataset.slug : "";
    const map = boot.titlePlaceholders || {};
    titleEl.placeholder = (slug && map[slug]) || boot.titlePlaceholder || titleEl.placeholder;
  }

  function onCategoryChange() {
    renderAttrs();
    updateTitlePlaceholder();
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
        select.dataset.type = (prop.type === "integer" || prop.type === "number") ? prop.type : "string";
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
          else input.step = "any";
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
  bindMap();
  categorySelect.addEventListener("change", onCategoryChange);
  applyPrefill();
  updateTitlePlaceholder();
  bindPriceField();
  if (pinLat == null) centerMapOnCity(false);
  syncPictureAvailability();

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

    const files = publishWithoutPhotos
      ? []
      : pendingPhotos.map(function (item) { return item.file; });
    const maxPics = boot.maxPictures || 8;
    if (files.length + existingPhotos.length > maxPics) {
      setMessage(form.dataset.tooMany || "", true);
      return;
    }

    if (files.length && !publishWithoutPhotos) {
      const uploadOk = await picturesUploadAvailable();
      if (!uploadOk) {
        showStorageContinue();
        return;
      }
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
    if (pinLat != null && pinLng != null) {
      payload.latitude = pinLat;
      payload.longitude = pinLng;
    }
    if (priceAmount != null) {
      payload.price_amount = priceAmount;
    }
    if (isEdit) {
      payload.keep_media = existingPhotos.map(function (item) { return item.storedUrl; }).filter(Boolean);
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
        if (err.code === "AD_STORAGE_UNAVAILABLE" && pendingPhotos.length && !publishWithoutPhotos) {
          showStorageContinue();
          return;
        }
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
          userEditedNeighborhood = false;
          const neighborhoodEl = document.getElementById("new-ad-neighborhood");
          if (neighborhoodEl) neighborhoodEl.value = "";
          centerMapOnCity(true);
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

  function bindMap() {
    if (adMap) return;
    if (typeof L === "undefined") {
      window.addEventListener("load", bindMap);
      return;
    }
    const el = document.getElementById("new-ad-map");
    if (!el) return;
    // Leaflet prefixes imagePath onto iconUrl. Hashed script URLs (?v=)
    // break auto-detect, so use absolute icon URLs and skip that prefix.
    delete L.Icon.Default.prototype._getIconUrl;
    L.Icon.Default.mergeOptions({
      iconRetinaUrl: "/static/vendor/leaflet/images/marker-icon-2x.png",
      iconUrl: "/static/vendor/leaflet/images/marker-icon.png",
      shadowUrl: "/static/vendor/leaflet/images/marker-shadow.png",
    });
    adMap = L.map(el, { scrollWheelZoom: true }).setView([35.6892, 51.3890], 11);
    adMap.attributionControl.setPrefix(
      "<a href=\"https://leafletjs.com\" target=\"_blank\" rel=\"noopener noreferrer\">Leaflet</a>"
    );
    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 19,
      attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\" target=\"_blank\" rel=\"noopener noreferrer\">OpenStreetMap</a>",
    }).addTo(adMap);
    openAttributionInNewTab(el);
    if (typeof window.addMapLocateButton === "function") {
      window.addMapLocateButton(adMap, el, function (latlng) {
        setMessage("");
        setMapPin(latlng.lat, latlng.lng, true);
        adMap.setView([latlng.lat, latlng.lng], 16);
      }, function (message) {
        if (message) setMessage(message, true);
      });
    }
    adMap.on("click", function (e) {
      setMapPin(e.latlng.lat, e.latlng.lng, true);
    });
    const neighborhoodEl = document.getElementById("new-ad-neighborhood");
    if (neighborhoodEl) {
      neighborhoodEl.addEventListener("input", function () {
        userEditedNeighborhood = true;
      });
    }
    window.setTimeout(function () {
      if (adMap) adMap.invalidateSize();
    }, 0);
    if (pinLat != null && pinLng != null) {
      setMapPin(pinLat, pinLng, false);
    } else {
      centerMapOnCity(false);
    }
  }

  function openAttributionInNewTab(mapEl) {
    if (!mapEl) return;
    mapEl.querySelectorAll(".leaflet-control-attribution a").forEach(function (link) {
      link.setAttribute("target", "_blank");
      link.setAttribute("rel", "noopener noreferrer");
    });
  }

  function clearMapPin() {
    pinLat = null;
    pinLng = null;
    if (adMarker && adMap) {
      adMap.removeLayer(adMarker);
      adMarker = null;
    }
  }

  function setMapPin(lat, lng, reverse) {
    lat = Number(lat);
    lng = Number(lng);
    if (!Number.isFinite(lat) || !Number.isFinite(lng)) return;
    pinLat = lat;
    pinLng = lng;
    if (!adMap || typeof L === "undefined") return;
    if (!adMarker) {
      adMarker = L.marker([lat, lng], { draggable: true }).addTo(adMap);
      adMarker.on("dragend", function () {
        const ll = adMarker.getLatLng();
        setMapPin(ll.lat, ll.lng, true);
      });
      adMap.setView([lat, lng], 15);
    } else {
      adMarker.setLatLng([lat, lng]);
      if (!reverse) adMap.setView([lat, lng], 15);
    }
    if (reverse) scheduleReverse(lat, lng);
  }

  function centerMapOnCity(clearPin) {
    if (clearPin) clearMapPin();
    if (!adMap) return;
    const qs = "slug=" + encodeURIComponent(boot.citySlug || "") +
      "&name=" + encodeURIComponent(boot.cityName || "") +
      "&lang=" + encodeURIComponent(boot.locale || "fa");
    fetch("/api/v1/geo/city?" + qs, { credentials: "same-origin" })
      .then(function (resp) { return resp.ok ? resp.json() : null; })
      .then(function (pt) {
        if (!pt || pt.lat == null || pt.lng == null) return;
        if (!clearPin && pinLat != null) return;
        adMap.setView([pt.lat, pt.lng], 12);
        adMap.invalidateSize();
      })
      .catch(function () {});
  }

  function scheduleReverse(lat, lng) {
    if (reverseTimer) window.clearTimeout(reverseTimer);
    reverseTimer = window.setTimeout(function () {
      const qs = "lat=" + encodeURIComponent(lat) +
        "&lng=" + encodeURIComponent(lng) +
        "&lang=" + encodeURIComponent(boot.locale || "fa");
      fetch("/api/v1/geo/reverse?" + qs, { credentials: "same-origin" })
        .then(function (resp) { return resp.ok ? resp.json() : null; })
        .then(function (data) {
          if (!data) return;
          const neighborhoodEl = document.getElementById("new-ad-neighborhood");
          if (!neighborhoodEl) return;
          if (userEditedNeighborhood && neighborhoodEl.value.trim()) return;
          if (data.neighborhood) neighborhoodEl.value = data.neighborhood;
        })
        .catch(function () {});
    }, 400);
  }

  function applyPrefill() {
    const p = boot.prefill;
    if (!p) return;
    if (p.category_id) {
      categorySelect.value = String(p.category_id);
      onCategoryChange();
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
    if (p.latitude != null && p.longitude != null) {
      setMapPin(Number(p.latitude), Number(p.longitude), false);
    }
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
    existingPhotos = [];
    (media || []).forEach(function (item) {
      const full = (item && (item.url || item.thumb)) || "";
      const thumb = (item && (item.thumb || item.url)) || "";
      const storedUrl = (item && (item.stored_url || item.storedUrl)) || "";
      if (!full && !thumb) return;
      existingPhotos.push({
        thumb: thumb || full,
        full: full || thumb,
        storedUrl: storedUrl,
      });
    });
    renderPhotoGrid();
  }

  async function picturesUploadAvailable() {
    try {
      const resp = await fetch("/api/v1/media/status", { credentials: "same-origin" });
      if (!resp.ok) return false;
      const body = await resp.json();
      return body && body.pictures_upload === true;
    } catch (_err) {
      return false;
    }
  }

  function showStorageContinue() {
    if (submitBtn) submitBtn.hidden = true;
    if (storageContinue) storageContinue.hidden = false;
    if (addPhotosBtn) addPhotosBtn.disabled = true;
    if (continueWithoutPhotosBtn) continueWithoutPhotosBtn.focus();
  }

  function applyPictureAvailability(ok) {
    boot.picturesUnavailable = !ok;
    if (picturesUnavailableHint) picturesUnavailableHint.hidden = ok;
    if (picturesHint) picturesHint.hidden = !ok;
    if (picturesInput) picturesInput.disabled = !ok;
    if (addPhotosBtn) {
      addPhotosBtn.hidden = !ok;
      addPhotosBtn.disabled = !ok ||
        (pendingPhotos.length + existingPhotos.length) >= (boot.maxPictures || 8);
    }
  }

  async function syncPictureAvailability() {
    applyPictureAvailability(await picturesUploadAvailable());
  }

  function tAttr(name, fallback) {
    return form.getAttribute(name) || fallback || "";
  }

  function isAllowedImage(file) {
    const type = (file.type || "").toLowerCase();
    if (type === "image/jpeg" || type === "image/png" || type === "image/webp") return true;
    return /\.(jpe?g|png|webp)$/i.test(file.name || "");
  }

  function alreadyPending(file) {
    return pendingPhotos.some(function (item) {
      return item.file.name === file.name &&
        item.file.size === file.size &&
        item.file.lastModified === file.lastModified;
    });
  }

  function addPendingFiles(fileList) {
    if (boot.picturesUnavailable) return;
    const maxPics = boot.maxPictures || 8;
    const incoming = Array.prototype.slice.call(fileList || []);
    let skipped = false;
    incoming.forEach(function (file) {
      if (!isAllowedImage(file) || alreadyPending(file)) return;
      if (pendingPhotos.length + existingPhotos.length >= maxPics) {
        skipped = true;
        return;
      }
      pendingPhotos.push({
        file: file,
        url: URL.createObjectURL(file),
      });
    });
    if (picturesInput) picturesInput.value = "";
    renderPhotoGrid();
    if (skipped) setMessage(form.dataset.tooMany || "", true);
    else setMessage("");
  }

  function removePendingAt(index) {
    const item = pendingPhotos[index];
    if (!item) return;
    if (item.url) URL.revokeObjectURL(item.url);
    pendingPhotos.splice(index, 1);
    afterPhotoRemoved();
  }

  function removeExistingAt(index) {
    if (!existingPhotos[index]) return;
    existingPhotos.splice(index, 1);
    afterPhotoRemoved();
  }

  function afterPhotoRemoved() {
    if (lightbox && !lightbox.hidden) {
      const sources = lightboxSources();
      if (!sources.length) closeLightbox();
      else showLightbox(Math.min(lightboxIndex, sources.length - 1));
    }
    renderPhotoGrid();
  }

  function visiblePhotos() {
    const out = [];
    existingPhotos.forEach(function (item, index) {
      out.push({ src: item.thumb, full: item.full, pending: false, existingIndex: index });
    });
    pendingPhotos.forEach(function (item, index) {
      out.push({ src: item.url, full: item.url, pending: true, pendingIndex: index });
    });
    return out;
  }

  function lightboxSources() {
    return visiblePhotos().map(function (item) { return item.full; });
  }

  function formatCounter(current, total) {
    let used = 0;
    const template = tAttr("data-photo-counter", "%d / %d");
    return template.replace(/%d/g, function () {
      used += 1;
      return used === 1 ? String(current) : String(total);
    });
  }

  function renderPhotoGrid() {
    if (!photoGrid) return;
    photoGrid.innerHTML = "";
    const photos = visiblePhotos();
    photos.forEach(function (item, index) {
      const tile = document.createElement("div");
      tile.className = "new-ad__photo";

      const openBtn = document.createElement("button");
      openBtn.type = "button";
      openBtn.className = "new-ad__photo-open";
      openBtn.setAttribute("aria-label", tAttr("data-picture-view", ""));
      const img = document.createElement("img");
      img.src = item.src;
      img.alt = "";
      openBtn.appendChild(img);
      openBtn.addEventListener("click", function () { showLightbox(index); });
      tile.appendChild(openBtn);

      const removeBtn = document.createElement("button");
      removeBtn.type = "button";
      removeBtn.className = "new-ad__photo-remove";
      removeBtn.setAttribute("aria-label", tAttr("data-picture-remove", ""));
      removeBtn.innerHTML = '<svg viewBox="0 0 24 24" fill="none" aria-hidden="true"><path d="M6 6l12 12M18 6 6 18" stroke="currentColor" stroke-width="2.25" stroke-linecap="round"/></svg>';
      removeBtn.addEventListener("click", function (ev) {
        ev.preventDefault();
        ev.stopPropagation();
        if (item.pending) removePendingAt(item.pendingIndex);
        else removeExistingAt(item.existingIndex);
      });
      tile.appendChild(removeBtn);
      photoGrid.appendChild(tile);
    });

    const maxPics = boot.maxPictures || 8;
    if (addPhotosBtn) {
      addPhotosBtn.hidden = !!boot.picturesUnavailable;
      addPhotosBtn.disabled = boot.picturesUnavailable ||
        (pendingPhotos.length + existingPhotos.length) >= maxPics;
    }
  }

  function showLightbox(index) {
    const sources = lightboxSources();
    if (!lightbox || !lightboxImg || !sources.length) return;
    lightboxIndex = (index + sources.length) % sources.length;
    lightboxImg.src = sources[lightboxIndex];
    if (lightboxCount) {
      lightboxCount.textContent = formatCounter(lightboxIndex + 1, sources.length);
      lightboxCount.hidden = sources.length < 2;
    }
    const many = sources.length > 1;
    if (lightboxPrev) lightboxPrev.hidden = !many;
    if (lightboxNext) lightboxNext.hidden = !many;
    lightbox.hidden = false;
    document.body.style.overflow = "hidden";
  }

  function closeLightbox() {
    if (!lightbox) return;
    lightbox.hidden = true;
    if (lightboxImg) lightboxImg.removeAttribute("src");
    document.body.style.overflow = "";
  }

  function stepLightbox(delta) {
    const sources = lightboxSources();
    if (sources.length < 2) return;
    showLightbox(lightboxIndex + delta);
  }

  if (addPhotosBtn && picturesInput) {
    addPhotosBtn.addEventListener("click", function () {
      if (addPhotosBtn.disabled) return;
      picturesInput.click();
    });
  }
  if (picturesInput) {
    picturesInput.addEventListener("change", function () {
      addPendingFiles(picturesInput.files);
    });
  }
  if (continueWithoutPhotosBtn) {
    continueWithoutPhotosBtn.addEventListener("click", function () {
      publishWithoutPhotos = true;
      if (storageContinue) storageContinue.hidden = true;
      if (submitBtn) submitBtn.hidden = false;
      if (typeof form.requestSubmit === "function") form.requestSubmit();
      else form.dispatchEvent(new Event("submit", { cancelable: true, bubbles: true }));
    });
  }
  if (lightbox) {
    const closeBtn = lightbox.querySelector("[data-photo-close]");
    if (closeBtn) closeBtn.addEventListener("click", closeLightbox);
    if (lightboxPrev) lightboxPrev.addEventListener("click", function () { stepLightbox(-1); });
    if (lightboxNext) lightboxNext.addEventListener("click", function () { stepLightbox(1); });
    lightbox.addEventListener("click", function (ev) {
      if (ev.target === lightbox) closeLightbox();
    });
  }
  document.addEventListener("keydown", function (ev) {
    if (!lightbox || lightbox.hidden) return;
    if (ev.key === "Escape") closeLightbox();
    if (ev.key === "ArrowLeft") stepLightbox(-1);
    if (ev.key === "ArrowRight") stepLightbox(1);
  });

  renderPhotoGrid();
})();
