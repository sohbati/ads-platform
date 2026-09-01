(function () {
  const COOKIE = "ruab-vid";
  const YEAR = 365 * 24 * 60 * 60;

  function readCookie(name) {
    const prefix = name + "=";
    const parts = document.cookie.split(";");
    for (let i = 0; i < parts.length; i++) {
      const part = parts[i].trim();
      if (part.indexOf(prefix) === 0) {
        return decodeURIComponent(part.slice(prefix.length));
      }
    }
    return "";
  }

  function viewerId() {
    const existing = readCookie(COOKIE);
    if (existing) return existing;
    const id = crypto.randomUUID();
    document.cookie = COOKIE + "=" + encodeURIComponent(id) + "; Max-Age=" + YEAR + "; Path=/; SameSite=Lax";
    return id;
  }

  function send(eventName, adId) {
    const id = Number(adId);
    if (!id || !eventName) return;
    const payload = JSON.stringify({
      ad_id: id,
      event: eventName,
      viewer_id: viewerId(),
      occurred_at: new Date().toISOString(),
    });
    const blob = new Blob([payload], { type: "application/json" });
    if (navigator.sendBeacon) {
      navigator.sendBeacon("/api/v1/stats/events", blob);
      return;
    }
    fetch("/api/v1/stats/events", {
      method: "POST",
      credentials: "same-origin",
      keepalive: true,
      headers: { "Content-Type": "application/json" },
      body: payload,
    }).catch(function () {});
  }

  window.RuabStats = {
    send: send,
    viewerId: viewerId,
  };

  const root = document.querySelector(".ad-detail[data-ad-id]");
  if (!root) return;
  const adId = root.getAttribute("data-ad-id");
  if (!adId) return;

  let sent = false;
  let timer = 0;

  function cancelTimer() {
    if (timer) {
      window.clearTimeout(timer);
      timer = 0;
    }
  }

  function armTimer() {
    cancelTimer();
    if (sent || document.hidden) return;
    timer = window.setTimeout(function () {
      timer = 0;
      if (sent || document.hidden) return;
      sent = true;
      send("view", adId);
    }, 1000);
  }

  document.addEventListener("visibilitychange", function () {
    if (document.hidden) {
      cancelTimer();
      return;
    }
    armTimer();
  });

  armTimer();
})();
