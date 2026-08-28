(function () {
  const root = document.querySelector("[data-ad-gallery]");
  if (!root) return;

  const slides = Array.prototype.slice.call(root.querySelectorAll("[data-ad-slide]"));
  if (slides.length < 2) return;

  const prevBtn = root.querySelector("[data-ad-prev]");
  const nextBtn = root.querySelector("[data-ad-next]");
  const countEl = root.querySelector("[data-ad-count]");
  const countTemplate = root.getAttribute("data-count-template") || "%d / %d";
  let index = 0;
  let startX = 0;

  function formatCount(current, total) {
    let used = 0;
    return countTemplate.replace(/%d/g, function () {
      used += 1;
      return used === 1 ? String(current) : String(total);
    });
  }

  function show(nextIndex) {
    index = (nextIndex + slides.length) % slides.length;
    slides.forEach(function (el, i) {
      el.classList.toggle("is-active", i === index);
    });
    if (countEl) countEl.textContent = formatCount(index + 1, slides.length);
  }

  function prev() {
    show(index - 1);
  }

  function next() {
    show(index + 1);
  }

  if (prevBtn) prevBtn.addEventListener("click", prev);
  if (nextBtn) nextBtn.addEventListener("click", next);

  document.addEventListener("keydown", function (ev) {
    if (ev.defaultPrevented) return;
    const tag = (ev.target && ev.target.tagName) || "";
    if (tag === "INPUT" || tag === "TEXTAREA" || tag === "SELECT") return;
    if (ev.key === "ArrowLeft") prev();
    if (ev.key === "ArrowRight") next();
  });

  const stage = root.querySelector(".ad-gallery__stage");
  if (stage) {
    stage.addEventListener("touchstart", function (ev) {
      if (!ev.changedTouches || !ev.changedTouches[0]) return;
      startX = ev.changedTouches[0].clientX;
    }, { passive: true });

    stage.addEventListener("touchend", function (ev) {
      if (!ev.changedTouches || !ev.changedTouches[0]) return;
      const dx = ev.changedTouches[0].clientX - startX;
      if (Math.abs(dx) < 40) return;
      if (dx > 0) prev();
      else next();
    }, { passive: true });
  }
})();
