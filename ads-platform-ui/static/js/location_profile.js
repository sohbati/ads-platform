(function () {
  const COOKIE_CITY = "city";
  const COOKIE_LOCATIONS = "locations";
  const COOKIE_MAX_AGE = 365 * 24 * 3600;
  const PROFILE_URL = "/api/v1/me/profile";

  function setCookie(name, value) {
    document.cookie = name + "=" + encodeURIComponent(value) + "; Path=/; Max-Age=" + COOKIE_MAX_AGE + "; SameSite=Lax";
  }

  function readCookie(name) {
    const parts = ("; " + document.cookie).split("; " + name + "=");
    if (parts.length < 2) return "";
    return decodeURIComponent(parts.pop().split(";").shift() || "");
  }

  function parseCSV(raw) {
    return String(raw || "")
      .split(",")
      .map(function (s) { return s.trim(); })
      .filter(Boolean);
  }

  function readCookieSlugs() {
    const fromLocations = parseCSV(readCookie(COOKIE_LOCATIONS));
    if (fromLocations.length) return fromLocations;
    const city = readCookie(COOKIE_CITY).trim();
    return city ? [city] : [];
  }

  function writeCookies(slugs) {
    const list = (slugs || []).filter(Boolean);
    setCookie(COOKIE_LOCATIONS, list.join(","));
    setCookie(COOKIE_CITY, list[0] || "tehran");
  }

  async function fetchProfile() {
    const resp = await fetch(PROFILE_URL, { credentials: "same-origin" });
    if (resp.status === 401) return null;
    if (!resp.ok) throw new Error("profile");
    return resp.json();
  }

  async function saveProfile(slugs) {
    const resp = await fetch(PROFILE_URL, {
      method: "PUT",
      credentials: "same-origin",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ location_slugs: slugs || [] }),
    });
    if (resp.status === 401) return false;
    return resp.ok;
  }

  async function save(slugs) {
    try {
      await saveProfile(slugs);
    } catch (_err) {
      // Guest or outage: cookies still apply for this browser.
    }
  }

  async function syncAfterLogin() {
    let profile = null;
    try {
      profile = await fetchProfile();
    } catch (_err) {
      return;
    }
    if (!profile) return;

    const saved = Array.isArray(profile.location_slugs) ? profile.location_slugs.filter(Boolean) : [];
    if (saved.length) {
      writeCookies(saved);
      return;
    }

    const cookies = readCookieSlugs();
    if (cookies.length) {
      await save(cookies);
    }
  }

  window.locationProfile = {
    readCookieSlugs: readCookieSlugs,
    writeCookies: writeCookies,
    save: save,
    syncAfterLogin: syncAfterLogin,
  };
})();
