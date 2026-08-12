(function () {
  function loadCatalog() {
    const el = document.getElementById("api-error-messages");
    if (!el) return {};
    try {
      return JSON.parse(el.textContent || "{}");
    } catch (_err) {
      return {};
    }
  }

  const catalog = loadCatalog();

  function resolveApiError(code, params, fallback) {
    let pattern = catalog[code] || catalog._default || fallback || code;
    if (!params || !params.length) {
      return pattern;
    }

    let message = pattern;
    for (let i = 0; i < params.length; i++) {
      message = message.replace("%s", params[i]);
    }
    return message;
  }

  window.resolveApiError = resolveApiError;
})();
