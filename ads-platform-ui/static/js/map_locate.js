(function (window) {
  const locateIcon =
    '<svg viewBox="0 0 24 24" width="18" height="18" aria-hidden="true" fill="none">' +
    '<path d="M12 3v2.5M12 18.5V21M3 12h2.5M18.5 12H21" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/>' +
    '<circle cx="12" cy="12" r="4" stroke="currentColor" stroke-width="1.8"/>' +
    '<circle cx="12" cy="12" r="1.2" fill="currentColor"/>' +
    "</svg>";

  function labelsFrom(el) {
    return {
      locate: (el && el.getAttribute("data-locate")) || "",
      denied: (el && el.getAttribute("data-locate-denied")) || "",
      failed: (el && el.getAttribute("data-locate-failed")) || "",
    };
  }

  function addMapLocateButton(map, mapEl, onFound, onError) {
    if (!map || typeof L === "undefined" || !L.Control) return;
    const copy = labelsFrom(mapEl);

    const Control = L.Control.extend({
      onAdd: function () {
        const wrap = L.DomUtil.create("div", "leaflet-bar map-locate-bar");
        const btn = L.DomUtil.create("button", "map-locate", wrap);
        btn.type = "button";
        btn.innerHTML = locateIcon;
        if (copy.locate) {
          btn.title = copy.locate;
          btn.setAttribute("aria-label", copy.locate);
        }
        L.DomEvent.disableClickPropagation(wrap);
        L.DomEvent.disableScrollPropagation(wrap);
        L.DomEvent.on(btn, "click", function (event) {
          L.DomEvent.stop(event);
          locate(map, btn, copy, onFound, onError);
        });
        return wrap;
      },
    });

    map.addControl(new Control({ position: "bottomright" }));
  }

  function locate(map, btn, copy, onFound, onError) {
    if (!navigator.geolocation) {
      if (onError) onError(copy.failed);
      return;
    }
    btn.disabled = true;
    btn.setAttribute("aria-busy", "true");
    navigator.geolocation.getCurrentPosition(
      function (pos) {
        btn.disabled = false;
        btn.removeAttribute("aria-busy");
        const latlng = {
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
        };
        if (onFound) onFound(latlng, pos.coords.accuracy);
        else map.setView(latlng, Math.max(map.getZoom(), 16));
      },
      function (err) {
        btn.disabled = false;
        btn.removeAttribute("aria-busy");
        const denied = err && err.code === 1;
        if (onError) onError(denied ? copy.denied : copy.failed);
      },
      { enableHighAccuracy: true, timeout: 12000, maximumAge: 10000 }
    );
  }

  window.addMapLocateButton = addMapLocateButton;
})(window);
