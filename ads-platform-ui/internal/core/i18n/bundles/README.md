# Language bundles

Each file in this folder is one UI language. The filename must match the locale code (e.g. `fa.json`, `en.json`).

## Add a new language

1. Copy `en.json` to `<locale>.json` (e.g. `ar.json`).
2. Set top-level fields:
   - `"locale"`: same as filename without `.json`
   - `"dir"`: `"rtl"` or `"ltr"`
   - `"label"`: name shown in UI (optional, for future switchers)
3. Translate every string under `"messages"`.
4. Add city slugs under `messages.city_names` if needed.
5. Rebuild the app (`make run` or `make build`) — bundles are embedded at compile time.

City names on the home page are fetched from ads-platform-cdn (`GET {CDN_BASE_URL}/json/cities.json`). Optional overrides per locale live under `messages.city_names` in each bundle (e.g. English names).

## JSON structure

```json
{
  "locale": "fa",
  "dir": "rtl",
  "label": "فارسی",
  "messages": {
    "meta": { "site_description": "..." },
    "header": { "categories": "...", ... },
    "hero": { "title": "… %s …", "subtitle": "..." },
    "categories": { ... },
    "cities": { ... },
    "cta": { ... },
    "footer": { ... },
    "error": { ... },
    "category_items": {
      "real-estate": { "name": "...", "description": "..." }
    },
    "city_names": {
      "tehran": "تهران"
    }
  }
}
```

Use `%s` in strings that need a city or app name; templates call `{{ format .T.Hero.Title .CityDisplayName }}`.
