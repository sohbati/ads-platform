(function () {
  const modal = document.getElementById("login-modal");
  if (!modal) return;

  const mobileInput = document.getElementById("login-mobile");
  const otpInput = document.getElementById("login-otp");
  const sendBtn = document.getElementById("login-send-otp");
  const verifyBtn = document.getElementById("login-verify-otp");
  const resendBtn = document.getElementById("login-resend-otp");
  const changeMobileBtn = document.getElementById("login-change-mobile");
  const messageEl = document.getElementById("login-message");
  const stepMobile = document.getElementById("login-step-mobile");
  const stepOtp = document.getElementById("login-step-otp");
  const nextInput = document.getElementById("login-next");
  const otpHint = document.getElementById("login-otp-hint");
  const accountTrigger = document.getElementById("header-login-trigger");
  const loginForm = document.getElementById("login-form");
  const defaultCountryCode = (loginForm && loginForm.dataset.countryCode) || "+98";

  const resolveError = typeof window.resolveApiError === "function"
    ? window.resolveApiError
    : function (_code, _params, fallback) {
        return fallback || "";
      };

  let lastFocused = null;

  function setMessage(text, isError) {
    if (!messageEl) return;
    messageEl.textContent = text || "";
    messageEl.classList.toggle("is-error", Boolean(isError));
  }

  function hideResendOtp() {
    if (!resendBtn) return;
    resendBtn.hidden = true;
    resendBtn.classList.add("is-hidden");
  }

  function showResendOtp() {
    if (!resendBtn) return;
    resendBtn.hidden = false;
    resendBtn.classList.remove("is-hidden");
  }

  function mobileValue() {
    const raw = (mobileInput && mobileInput.value || "").trim();
    if (!raw) {
      return { ok: false, error: "MOBILE_EMPTY" };
    }
    if (typeof window.normalizeMobile !== "function") {
      return { ok: true, mobile: raw };
    }
    return window.normalizeMobile(raw, defaultCountryCode);
  }

  function mobileError(result) {
    setMessage(messageForError(result.error || "INVALID_MOBILE", []), true);
  }

  async function readErrorResponse(resp) {
    const text = await resp.text();
    try {
      const data = JSON.parse(text);
      if (data && typeof data.error === "string") {
        return {
          code: data.error,
          params: Array.isArray(data.params) ? data.params : [],
        };
      }
    } catch (_err) {
      // Plain-text error body.
    }
    return { code: "", params: [] };
  }

  function messageForError(code, params) {
    return resolveError(code, params, resolveError("_default", [], code));
  }

  async function requestOtp(mobile) {
    const resp = await fetch("/api/v1/auth/otp/" + encodeURIComponent(mobile) + "/send", {
      method: "POST",
      credentials: "same-origin",
    });

    if (!resp.ok) {
      const err = await readErrorResponse(resp);
      throw Object.assign(new Error(messageForError(err.code, err.params)), { code: err.code });
    }
  }

  function showStepMobile() {
    stepMobile.classList.remove("is-hidden");
    stepMobile.hidden = false;
    stepOtp.classList.add("is-hidden");
    stepOtp.hidden = true;
    if (otpInput) otpInput.value = "";
    hideResendOtp();
    setMessage("");
  }

  function showStepOtp(mobile) {
    stepMobile.classList.add("is-hidden");
    stepMobile.hidden = true;
    stepOtp.classList.remove("is-hidden");
    stepOtp.hidden = false;
    hideResendOtp();
    if (otpHint && mobile) {
      otpHint.textContent = otpHint.dataset.template
        ? otpHint.dataset.template.replace("%s", mobile)
        : otpHint.textContent;
    }
    if (otpInput) otpInput.focus();
  }

  function openModal(nextPath) {
    if (document.body.classList.contains("category-modal-open") && typeof window.closeCategoryModal === "function") {
      window.closeCategoryModal();
    }
    if (document.body.classList.contains("location-modal-open") && typeof window.closeLocationModal === "function") {
      window.closeLocationModal();
    }
    if (nextInput && nextPath) {
      nextInput.value = nextPath;
    }
    showStepMobile();
    lastFocused = document.activeElement;
    modal.hidden = false;
    modal.setAttribute("aria-hidden", "false");
    modal.classList.add("is-open");
    document.body.classList.add("login-modal-open");
    if (mobileInput) mobileInput.focus();
  }

  function closeModal() {
    modal.classList.remove("is-open");
    modal.hidden = true;
    modal.setAttribute("aria-hidden", "true");
    document.body.classList.remove("login-modal-open");
    setMessage("");
    showStepMobile();
    if (lastFocused && typeof lastFocused.focus === "function") {
      lastFocused.focus();
    }
  }

  async function isAuthenticated() {
    try {
      const resp = await fetch("/api/v1/auth/me", { credentials: "same-origin" });
      if (!resp.ok) return false;
      const data = await resp.json();
      return Boolean(data && data.authenticated);
    } catch (_err) {
      return false;
    }
  }

  async function logout() {
    await fetch("/api/v1/auth/logout", { method: "POST", credentials: "same-origin" });
    window.location.assign("/query-ads");
  }

  window.openLoginModal = openModal;
  window.closeLoginModal = closeModal;
  window.logout = logout;

  document.querySelectorAll("[data-logout]").forEach(function (btn) {
    btn.addEventListener("click", logout);
  });

  modal.querySelectorAll("[data-login-close]").forEach(function (el) {
    el.addEventListener("click", closeModal);
  });

  document.addEventListener("keydown", function (event) {
    if (event.key === "Escape" && modal.classList.contains("is-open")) {
      closeModal();
    }
  });

  const accountMenu = document.getElementById("account-menu");

  function setAccountMenuAuth(authed) {
    if (!accountMenu) return;
    accountMenu.querySelectorAll("[data-account-authed]").forEach(function (el) {
      el.hidden = !authed;
    });
    accountMenu.querySelectorAll("[data-account-guest]").forEach(function (el) {
      el.hidden = authed;
    });
  }

  function closeAccountMenu() {
    if (!accountMenu) return;
    accountMenu.hidden = true;
    if (accountTrigger) accountTrigger.setAttribute("aria-expanded", "false");
  }

  function openAccountMenu() {
    if (!accountMenu) return;
    accountMenu.hidden = false;
    if (accountTrigger) accountTrigger.setAttribute("aria-expanded", "true");
  }

  isAuthenticated().then(setAccountMenuAuth);

  if (accountTrigger && accountMenu) {
    accountTrigger.addEventListener("click", async function (event) {
      event.stopPropagation();
      const authed = await isAuthenticated();
      setAccountMenuAuth(authed);
      if (accountMenu.hidden) openAccountMenu();
      else closeAccountMenu();
    });

    document.addEventListener("click", function (event) {
      if (accountMenu.hidden) return;
      if (accountMenu.contains(event.target) || accountTrigger.contains(event.target)) return;
      closeAccountMenu();
    });

    document.addEventListener("keydown", function (event) {
      if (event.key === "Escape") closeAccountMenu();
    });
  }

  const accountLogin = document.querySelector("[data-account-login]");
  if (accountLogin) {
    accountLogin.addEventListener("click", function () {
      closeAccountMenu();
      openModal("/my-info");
    });
  }

  if (changeMobileBtn) {
    changeMobileBtn.addEventListener("click", showStepMobile);
  }

  async function handleSendOtp() {
    const mobileResult = mobileValue();
    if (!mobileResult.ok) {
      mobileError(mobileResult);
      return;
    }
    const mobile = mobileResult.mobile;

    setMessage("");
    sendBtn.disabled = true;

    try {
      await requestOtp(mobile);
      showStepOtp(mobile);
      setMessage(messageEl.dataset.sent || resolveError("OTP_SENT", [], "Code sent."), false);
    } catch (err) {
      setMessage(err.message || messageForError("_default", []), true);
    } finally {
      sendBtn.disabled = false;
    }
  }

  if (sendBtn) {
    sendBtn.addEventListener("click", handleSendOtp);
  }

  async function handleResendOtp() {
    const mobileResult = mobileValue();
    if (!mobileResult.ok) {
      mobileError(mobileResult);
      return;
    }
    const mobile = mobileResult.mobile;

    setMessage("");
    resendBtn.disabled = true;

    try {
      await requestOtp(mobile);
      if (otpInput) otpInput.value = "";
      hideResendOtp();
      setMessage(messageEl.dataset.sent || resolveError("OTP_SENT", [], "Code sent."), false);
      if (otpInput) otpInput.focus();
    } catch (err) {
      setMessage(err.message || messageForError("_default", []), true);
    } finally {
      resendBtn.disabled = false;
    }
  }

  if (resendBtn) {
    resendBtn.addEventListener("click", handleResendOtp);
  }

  if (verifyBtn) {
    verifyBtn.addEventListener("click", async function () {
      const mobileResult = mobileValue();
      if (!mobileResult.ok) {
        mobileError(mobileResult);
        return;
      }
      const mobile = mobileResult.mobile;
      const otp = (otpInput && otpInput.value || "").trim();

      if (otp.length !== 6) {
        setMessage(messageForError("INVALID_OTP", []), true);
        return;
      }

      setMessage("");
      hideResendOtp();
      verifyBtn.disabled = true;

      try {
        const resp = await fetch("/api/v1/auth/otp/" + encodeURIComponent(mobile) + "/verify", {
          method: "POST",
          credentials: "same-origin",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ otp: otp }),
        });

        if (!resp.ok) {
          const err = await readErrorResponse(resp);
          setMessage(messageForError(err.code, err.params), true);
          if (err.code === "OTP_EXPIRED" || err.code === "OTP_NOT_FOUND") {
            showResendOtp();
          }
          return;
        }

        if (window.locationProfile && typeof window.locationProfile.syncAfterLogin === "function") {
          try {
            await window.locationProfile.syncAfterLogin();
          } catch (_syncErr) {
            // Cookies still apply for this visit if profile sync fails.
          }
        }

        closeModal();
        const next = (nextInput && nextInput.value) || "/my-info";
        window.location.assign(next);
      } catch (_err) {
        setMessage(messageForError("NETWORK_ERROR", []), true);
      } finally {
        verifyBtn.disabled = false;
      }
    });
  }

  const params = new URLSearchParams(window.location.search);
  if (params.get("login") === "1") {
    const next = params.get("next") || "/my-info";
    openModal(next);
    params.delete("login");
    params.delete("next");
    const clean = params.toString();
    const newUrl = window.location.pathname + (clean ? "?" + clean : "") + window.location.hash;
    window.history.replaceState({}, "", newUrl);
  }
})();
