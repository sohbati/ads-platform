(function () {
  var KEY = "ruab-theme";
  var LEGACY = "ruab-sea";
  var ALLOWED = { light: true, dark: true, tide: true };
  var root = document.documentElement;
  var picker = document.getElementById("theme-picker");

  function currentId() {
    try {
      var id = localStorage.getItem(KEY);
      return ALLOWED[id] ? id : "light";
    } catch (_err) {
      return "light";
    }
  }

  function apply(id) {
    if (!ALLOWED[id]) {
      id = "light";
    }
    if (id === "light") {
      root.removeAttribute("data-theme");
      try {
        localStorage.removeItem(KEY);
      } catch (_err) {}
    } else {
      root.setAttribute("data-theme", id);
      try {
        localStorage.setItem(KEY, id);
      } catch (_err) {}
    }

    if (!picker) return;
    var buttons = picker.querySelectorAll("[data-theme-id]");
    for (var i = 0; i < buttons.length; i++) {
      var on = buttons[i].getAttribute("data-theme-id") === id;
      buttons[i].classList.toggle("is-selected", on);
      buttons[i].setAttribute("aria-pressed", on ? "true" : "false");
      buttons[i].setAttribute("aria-checked", on ? "true" : "false");
    }
  }

  try {
    localStorage.removeItem(LEGACY);
  } catch (_err) {}

  if (picker) {
    picker.addEventListener("click", function (ev) {
      var btn = ev.target.closest("[data-theme-id]");
      if (!btn) return;
      apply(btn.getAttribute("data-theme-id"));
    });
  }

  apply(currentId());
})();
