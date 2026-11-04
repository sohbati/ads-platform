// Shows dates in the visitor's own time zone and language. The server sends the exact moment in
// <time datetime="…" data-local-date> with a UTC date as fallback text; this replaces the text.
// Persian pages use the Jalali calendar (۱۳ آبان ۱۴۰۵), others their own (4 Nov 2026).
(function () {
  const lang = (document.documentElement.lang || "fa").toLowerCase();
  const locale = lang.startsWith("fa") ? "fa-IR-u-ca-persian" : lang;

  let formatter;
  try {
    formatter = new Intl.DateTimeFormat(locale, { year: "numeric", month: "long", day: "numeric" });
  } catch (_) {
    return; // keep the server's fallback text
  }

  function format(iso) {
    const date = new Date(iso);
    return Number.isNaN(date.getTime()) ? "" : formatter.format(date);
  }

  function localize(root) {
    (root || document).querySelectorAll("time[data-local-date]").forEach(function (el) {
      const text = format(el.getAttribute("datetime"));
      if (text) el.textContent = text;
    });
  }

  window.LocalTime = { format: format, localize: localize };
  localize(document);
})();
