// Package i18n loads per-language JSON bundles and resolves locale from the selected city.
//
// City data is fetched from ads-platform-cdn over HTTP (see internal/core/cities).
//
// Add a new language:
//  1. Create bundles/<locale>.json (copy en.json or fa.json as a template).
//  2. Set "locale", "dir" (rtl|ltr), "label", and fill "messages".
//  3. Rebuild the app — bundles are embedded at compile time.
package i18n
