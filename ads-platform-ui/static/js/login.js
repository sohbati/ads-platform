(function () {
  const mobileInput = document.getElementById("login-mobile");
  const otpInput = document.getElementById("login-otp");
  const sendBtn = document.getElementById("login-send-otp");
  const verifyBtn = document.getElementById("login-verify-otp");
  const messageEl = document.getElementById("login-message");
  const stepMobile = document.getElementById("login-step-mobile");
  const stepOtp = document.getElementById("login-step-otp");
  const nextInput = document.getElementById("login-next");

  if (!mobileInput || !sendBtn || !verifyBtn) {
    return;
  }

  function setMessage(text, isError) {
    messageEl.textContent = text || "";
    messageEl.classList.toggle("is-error", Boolean(isError));
  }

  function mobileValue() {
    return (mobileInput.value || "").trim();
  }

  sendBtn.addEventListener("click", async function () {
    const mobile = mobileValue();
    if (!mobile) {
      setMessage("Enter your mobile number.", true);
      return;
    }

    setMessage("");
    sendBtn.disabled = true;

    try {
      const resp = await fetch("/api/v1/auth/otp/" + encodeURIComponent(mobile) + "/send", {
        method: "POST",
        credentials: "same-origin",
      });

      if (!resp.ok) {
        const body = await resp.text();
        setMessage(body || "Failed to send OTP.", true);
        return;
      }

      stepMobile.classList.add("is-hidden");
      stepOtp.classList.remove("is-hidden");
      otpInput.focus();
      setMessage("OTP sent.", false);
    } catch (err) {
      setMessage("Network error while sending OTP.", true);
    } finally {
      sendBtn.disabled = false;
    }
  });

  verifyBtn.addEventListener("click", async function () {
    const mobile = mobileValue();
    const otp = (otpInput.value || "").trim();

    if (!mobile || otp.length !== 6) {
      setMessage("Enter the 6-digit OTP.", true);
      return;
    }

    setMessage("");
    verifyBtn.disabled = true;

    try {
      const resp = await fetch("/api/v1/auth/otp/" + encodeURIComponent(mobile) + "/verify", {
        method: "POST",
        credentials: "same-origin",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ otp: otp }),
      });

      if (!resp.ok) {
        const body = await resp.text();
        setMessage(body || "Login failed.", true);
        return;
      }

      const next = (nextInput && nextInput.value) || "/my-info";
      window.location.assign(next);
    } catch (err) {
      setMessage("Network error while verifying OTP.", true);
    } finally {
      verifyBtn.disabled = false;
    }
  });
})();
