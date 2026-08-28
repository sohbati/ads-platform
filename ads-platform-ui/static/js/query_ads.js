(function () {
  const feed = document.getElementById("ad-feed");
  const grid = document.getElementById("ad-grid");
  const sentinel = document.getElementById("ad-feed-sentinel");
  const statusEl = document.getElementById("ad-feed-status");
  if (!feed || !grid || !sentinel) return;

  const PLACEHOLDER =
    '<svg class="icon ad-card__placeholder" viewBox="0 0 24 24" fill="none" aria-hidden="true">' +
    '<rect x="3" y="5" width="18" height="14" rx="2" stroke="currentColor" stroke-width="1.5"/>' +
    '<circle cx="9" cy="10" r="1.75" stroke="currentColor" stroke-width="1.5"/>' +
    '<path d="m5 17 4.5-4 3.5 3 2.5-2 3.5 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>' +
    "</svg>";

  let page = parseInt(feed.getAttribute("data-page") || "1", 10) || 1;
  let hasMore = feed.getAttribute("data-has-more") === "true";
  let loading = false;
  const query = feed.getAttribute("data-query") || "";
  const category = feed.getAttribute("data-category") || "";
  const feedURL = feed.getAttribute("data-feed-url") || "/api/v1/search";

  function t(key, fallback) {
    return feed.getAttribute(key) || fallback || "";
  }

  function escapeHtml(value) {
    return String(value || "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;");
  }

  function setStatus(html, isError) {
    if (!statusEl) return;
    statusEl.hidden = !html;
    statusEl.classList.toggle("is-error", Boolean(isError));
    statusEl.innerHTML = html || "";
  }

  function cardHTML(ad) {
    const media = ad.thumbnail
      ? '<img class="ad-card__img" src="' + escapeHtml(ad.thumbnail) + '" alt="" loading="lazy" />'
      : PLACEHOLDER;
    const href = ad.href || (ad.id ? "/ad/" + ad.id : "");
    const open = href
      ? '<a class="ad-card is-appended" href="' + escapeHtml(href) + '">'
      : '<div class="ad-card is-appended">';
    const close = href ? "</a>" : "</div>";
    const meta = [ad.location, ad.published_at].filter(Boolean).join(" · ");
    const price = ad.price
      ? '<p class="ad-card__price">' + escapeHtml(ad.price) + "</p>"
      : "";
    return (
      "<li>" +
        open +
        '<div class="ad-card__media">' + media + price + "</div>" +
        '<div class="ad-card__body">' +
          '<h2 class="ad-card__title">' + escapeHtml(ad.title) + "</h2>" +
          '<p class="ad-card__meta">' + escapeHtml(meta) + "</p>" +
        "</div>" +
        close +
      "</li>"
    );
  }

  function nextURL(nextPage) {
    const params = new URLSearchParams();
    if (query) params.set("q", query);
    if (category) params.set("category", category);
    params.set("page", String(nextPage));
    return feedURL + "?" + params.toString();
  }

  function sentinelInView() {
    const rect = sentinel.getBoundingClientRect();
    return rect.top < (window.innerHeight || 0) + 480;
  }

  async function loadMore() {
    if (!hasMore || loading) return;
    loading = true;
    setStatus('<span class="ad-feed__spinner" aria-hidden="true"></span> ' + t("data-i18n-loading", ""), false);

    let loaded = false;
    try {
      const resp = await fetch(nextURL(page + 1), { credentials: "same-origin" });
      if (!resp.ok) throw new Error("feed");
      const data = await resp.json();
      const ads = Array.isArray(data.ads) ? data.ads : [];
      if (ads.length) {
        grid.insertAdjacentHTML("beforeend", ads.map(cardHTML).join(""));
      }
      page = Number(data.page) || page + 1;
      hasMore = Boolean(data.has_more) && ads.length > 0;
      feed.setAttribute("data-page", String(page));
      feed.setAttribute("data-has-more", hasMore ? "true" : "false");
      setStatus("", false);
      loaded = true;
      if (!hasMore && observer) observer.disconnect();
    } catch (_err) {
      setStatus(
        escapeHtml(t("data-i18n-error", "")) +
          ' <button type="button" class="ad-feed__retry" data-ad-retry>' +
          escapeHtml(t("data-i18n-retry", "")) +
          "</button>",
        true
      );
    } finally {
      loading = false;
      if (loaded && hasMore && sentinelInView()) loadMore();
    }
  }

  if (statusEl) {
    statusEl.addEventListener("click", function (ev) {
      const btn = ev.target.closest("[data-ad-retry]");
      if (!btn) return;
      loadMore();
    });
  }

  const observer = new IntersectionObserver(
    function (entries) {
      if (entries.some(function (e) { return e.isIntersecting; })) loadMore();
    },
    { root: null, rootMargin: "480px 0px", threshold: 0 }
  );

  if (hasMore) observer.observe(sentinel);
})();
