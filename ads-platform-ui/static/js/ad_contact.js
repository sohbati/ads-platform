(function () {
  const root = document.querySelector("[data-ad-contact]");
  if (!root) return;

  const adId = root.getAttribute("data-ad-id");
  const maskedEl = root.querySelector("[data-ad-phone-masked]");
  const revealedEl = root.querySelector("[data-ad-phone-revealed]");
  const showBtn = root.querySelector("[data-ad-show-phone]");
  const callLink = root.querySelector("[data-ad-call-phone]");

  function revealNext() {
    return window.location.pathname + "?reveal=1";
  }

  function telHref(phone) {
    const digits = String(phone || "").replace(/\D/g, "");
    if (digits.indexOf("09") === 0 && digits.length === 11) {
      return "tel:+98" + digits.slice(1);
    }
    if (digits.indexOf("98") === 0) {
      return "tel:+" + digits;
    }
    return "tel:" + phone;
  }

  function showPhone(phone) {
    root.classList.add("is-revealed");
    if (maskedEl) maskedEl.hidden = true;
    if (showBtn) showBtn.hidden = true;
    if (revealedEl) {
      revealedEl.textContent = phone;
      revealedEl.href = telHref(phone);
      revealedEl.hidden = false;
    }
    if (callLink) {
      callLink.href = telHref(phone);
      callLink.hidden = false;
    }
  }

  function openLogin() {
    if (typeof window.openLoginModal === "function") {
      window.openLoginModal(revealNext());
      return;
    }
    window.location.assign("/login?next=" + encodeURIComponent(revealNext()));
  }

  async function requestPhone() {
    if (root.classList.contains("is-revealed")) return;
    try {
      const resp = await fetch("/api/v1/ads/" + encodeURIComponent(adId) + "/contact", {
        credentials: "same-origin",
      });
      if (resp.status === 401) {
        openLogin();
        return;
      }
      if (!resp.ok) return;
      const data = await resp.json();
      if (data && data.phone) showPhone(data.phone);
    } catch (_err) {
      // Keep the masked number if the request fails.
    }
  }

  function onRevealClick(event) {
    if (root.classList.contains("is-revealed")) return;
    if (!event.target.closest("[data-ad-show-phone], [data-ad-phone-masked]")) return;
    event.preventDefault();
    requestPhone();
  }

  root.addEventListener("click", onRevealClick);

  const params = new URLSearchParams(window.location.search);
  if (params.get("reveal") === "1") {
    params.delete("reveal");
    const clean = params.toString();
    window.history.replaceState({}, "", window.location.pathname + (clean ? "?" + clean : "") + window.location.hash);
    requestPhone();
  }
})();
