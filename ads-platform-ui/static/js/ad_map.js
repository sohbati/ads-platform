(function () {
  const el = document.getElementById("ad-detail-map");
  if (!el) return;

  function bindMap() {
    if (typeof L === "undefined") {
      window.addEventListener("load", bindMap);
      return;
    }
    const lat = parseFloat(el.getAttribute("data-lat") || "");
    const lng = parseFloat(el.getAttribute("data-lng") || "");
    if (!Number.isFinite(lat) || !Number.isFinite(lng)) return;

    const map = L.map(el, {
      scrollWheelZoom: true,
      dragging: true,
      doubleClickZoom: true,
    }).setView([lat, lng], 14);
    map.attributionControl.setPrefix(
      "<a href=\"https://leafletjs.com\" target=\"_blank\" rel=\"noopener noreferrer\">Leaflet</a>"
    );
    L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
      maxZoom: 19,
      attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\" target=\"_blank\" rel=\"noopener noreferrer\">OpenStreetMap</a>",
    }).addTo(map);
    el.querySelectorAll(".leaflet-control-attribution a").forEach(function (link) {
      link.setAttribute("target", "_blank");
      link.setAttribute("rel", "noopener noreferrer");
    });

    const color = (window.getComputedStyle(document.documentElement).getPropertyValue("--color-primary") || "#0d9488").trim();
    const circle = L.circle([lat, lng], {
      radius: 400,
      color: color,
      fillColor: color,
      fillOpacity: 0.22,
      weight: 2,
      interactive: false,
    }).addTo(map);
    map.fitBounds(circle.getBounds(), { padding: [24, 24], maxZoom: 15 });
    let youMarker = null;
    if (typeof window.addMapLocateButton === "function") {
      window.addMapLocateButton(map, el, function (latlng) {
        const status = document.getElementById("ad-detail-map-status");
        if (status) status.textContent = "";
        if (youMarker) map.removeLayer(youMarker);
        youMarker = L.circleMarker(latlng, {
          radius: 8,
          color: "#fff",
          weight: 2,
          fillColor: color,
          fillOpacity: 1,
        }).addTo(map);
        map.setView(latlng, Math.max(map.getZoom(), 15));
      }, function (message) {
        const status = document.getElementById("ad-detail-map-status");
        if (status) status.textContent = message || "";
      });
    }
    window.setTimeout(function () {
      map.invalidateSize();
    }, 0);
  }

  bindMap();
})();
