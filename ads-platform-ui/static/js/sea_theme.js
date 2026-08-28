(function () {
  var KEY = "ruab-sea";
  var picker = document.getElementById("sea-picker");
  if (!picker) return;

  var root = document.documentElement;

  function currentId() {
    try {
      return localStorage.getItem(KEY) || "ruab";
    } catch (_err) {
      return "ruab";
    }
  }

  function apply(id) {
    if (!id || id === "ruab") {
      root.removeAttribute("data-sea");
      try {
        localStorage.removeItem(KEY);
      } catch (_err) {}
      id = "ruab";
    } else {
      root.setAttribute("data-sea", id);
      try {
        localStorage.setItem(KEY, id);
      } catch (_err) {}
    }

    var buttons = picker.querySelectorAll("[data-sea-id]");
    for (var i = 0; i < buttons.length; i++) {
      var on = buttons[i].getAttribute("data-sea-id") === id;
      buttons[i].classList.toggle("is-selected", on);
      buttons[i].setAttribute("aria-pressed", on ? "true" : "false");
      buttons[i].setAttribute("aria-checked", on ? "true" : "false");
    }
  }

  picker.addEventListener("click", function (ev) {
    var btn = ev.target.closest("[data-sea-id]");
    if (!btn) return;
    apply(btn.getAttribute("data-sea-id"));
  });

  apply(currentId());
})();
