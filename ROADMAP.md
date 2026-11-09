# Ads Platform Roadmap

Status: MVP. Last reviewed: 2026-11-04.

Legend: [x] done · [ ] not started · [~] started or partly done

Phases are in order: finish **Phase 1** before launch, **Phase 2** in the first months after,
**Phase 3** once there is steady traffic. Within a phase, items are roughly in priority order.

---

## Already built

- [x] Customer sign-in with mobile + OTP (BFF sessions)
- [x] Post, edit, and browse ads; categories, cities, map location, photos via media CDN
- [x] Public search by place and category, with ad attributes per category
- [x] Ad view / contact-reveal stats (ads-platform-stats)
- [x] CRM: staff sign-in, staff, roles, actions, staff groups (IAM)
- [x] Ad review: pending queue, approve / reject with customer and internal reasons
- [x] Review claims: one reviewer per ad, "Next ad", reviewer role, today's counts
- [x] CRM ads search that scales: ID / mobile / title, category, city, cursor paging
- [x] Category tree in the CRM (read-only)
- [x] Runtime settings (`ads.pending_max_days`)
- [x] Timestamps stored as UTC, shown in the viewer's zone (Jalali for Persian)

---

## Phase 1 — Finish the MVP (before launch)

Gaps where the UI, a permission, or a setting already exists but nothing works behind it.

- [ ] **Real SMS provider** (e.g. Kavenegar, SMS.ir) in ads-platform-notification.
      OTP SMS are only logged today, so nobody outside development can sign in.
- [ ] **Pending-timeout job.** `ads.pending_max_days` is stored but nothing acts on it.
      Decide what happens to an ad that waits too long (back to draft, or escalate to a manager).
- [ ] **Ad expiry.** `ads.expires_at` exists but is never set. Set it on approval
      (e.g. 30 days), hide expired ads from search, and let the owner renew.
- [ ] **Owner controls:** delete an ad, and mark it "sold" / take it down. Customers can only
      edit today.
- [ ] **Notify the customer about review results** (SMS or in-app): approved, or rejected
      with the reason. Today they only find out by opening "my ads".
- [ ] **Bookmarks ("marked ads").** The page exists but shows "coming soon".
- [ ] **CRM customers page.** Calls a placeholder endpoint and is always empty.
      List and search customers, see their ads, and block them (`customer.block` exists).
- [ ] **CRM ad delete.** `advertisement.delete` exists without an endpoint; needed for
      removing illegal or spam ads that are already live.
- [ ] **CRM dashboard numbers.** Returns zeros: pending count, oldest waiting ad,
      ads approved today, new customers.
- [ ] **Category management.** Categories are a JSON file in the CDN; the
      `category.create/update/delete` permissions are unused. At minimum, document how to
      change categories safely; ideally edit them from the CRM.
- [ ] **Deployment:** Dockerfiles and a compose file for all services, environment-specific
      config, and one command to start everything (today: 11 shell scripts).
- [ ] **Database migrations runner.** SQL scripts (`script_1` … `script_12`) are run by hand;
      track which ones each database has applied (golang-migrate, goose, or a simple
      `schema_migrations` table).
- [ ] **Backups** of PostgreSQL and the media store, with a tested restore.

### Marketplace basics

What buyers and sellers expect from a classifieds site on day one.

- [~] **Search filters on the site.** The search API already supports price range, "with photo",
      sort (newest / cheapest / most expensive), several cities, and neighborhood; the public
      pages expose none of them. Add a filter panel and sort menu.
- [ ] **Category-specific filters:** e.g. cars by year, mileage, gearbox; property by area, rooms,
      floor, parking; phones by brand and storage. Built from the existing attribute schemas.
- [ ] **Relative times and place on cards:** "15 minutes ago in Tehran, Saadat Abad" instead of a
      date; it tells buyers how fresh an ad is at a glance.
- [ ] **Report an ad** (moved up from Phase 2): wrong price, sold, scam, wrong category,
      offensive. Reports go to a CRM queue.
- [ ] **Safety notice when revealing contact info:** don't pay in advance, meet in person,
      check documents. Shown once before the phone number appears.
- [ ] **Hide my number / chat only** option for sellers; needs in-app chat (see Phase 3),
      or at least a masked number until the buyer signs in.
- [ ] **Similar ads** on the ad page (same category and city, close price).
- [ ] **Share an ad:** copy link and the phone's native share sheet.
- [ ] **Recently viewed ads** and **recent searches** for the visitor.
- [ ] **Personal note on a bookmarked ad** ("called, call back Friday"), visible only to the owner.
- [ ] **Ad flags set by the seller:** urgent, exchange accepted, negotiable price (negotiable
      exists), shown as badges on cards and as search filters.
- [ ] **Map view of results** for property ads, using the stored location.
- [ ] **Posting flow helpers:** suggest the category from the title, save a draft while posting,
      reorder photos and choose the cover, and show posting rules before submit.
- [ ] **Renew or re-post** an expired or sold ad in one step, without retyping it.
- [ ] **Seller's other ads** link on the ad page.

## Phase 2 — Launch readiness and trust (first months)

- [ ] **Rejection reason templates.** Common reasons as a dropdown in the reject dialog;
      faster reviews and consistent wording. Keep them in runtime settings.
- [ ] **Staff audit log.** Who approved, rejected, deleted, blocked, or changed settings and
      roles, and when. Approvals are recorded today; nothing else is.
- [ ] **Rate limits and abuse protection** on public endpoints: OTP send per mobile and per IP,
      posting ads per customer per day, contact reveals.
- [ ] **Duplicate and spam detection.** Flag ads whose title/photos match a recent ad, or
      customers posting the same ad repeatedly, for the reviewer.
- [ ] **SEO:** `sitemap.xml`, `robots.txt`, canonical URLs, Open Graph tags on ad pages,
      structured data (schema.org Offer/Product), readable URLs with the slug.
- [ ] **Monitoring:** structured logs, health dashboard for all services, alerts when the
      review queue is older than N hours or a service is down.
- [ ] **CI:** run Go tests, crm-ui tests/lint/type-check, and integration tests on every push.
- [ ] **CRM date-range filter** with a Jalali date picker (left out of ads search for now).
- [ ] **Legal pages:** terms of use, privacy policy, posting rules (linked from the reject
      reasons).

## Phase 3 — Growth

- [ ] **In-app chat** between buyer and seller, so sellers don't have to expose their number.
      Consider moving it earlier: on many classifieds sites, chat is how most buyers first
      contact a seller.
- [ ] **Paid features:** bump an ad to the top, featured/highlighted ads, extra photos;
      requires a payment gateway (e.g. Zarinpal, IDPay) and an orders/invoices model.
- [ ] **Saved searches and alerts:** notify a customer when a new ad matches their search.
- [ ] **Search engine (Elasticsearch/OpenSearch)** for public and CRM search: typo tolerance,
      Persian normalization (ی/ي, ک/ك), attribute filters, price ranges, sort by distance.
- [ ] **Seller profiles and ratings**, verified sellers (national ID check), business accounts
      with many ads.
- [ ] **Photo moderation:** automatic checks for unsafe images, faces, phone numbers in photos.
- [ ] **Customer insights in the CRM** (the Reports page is a placeholder): ads per day, approval
      rate, review time, popular categories and cities, views and contacts per ad.
- [ ] **Mobile apps or PWA** with push notifications.
- [ ] **Archive old ads** (expired/deleted) out of the live table once it reaches tens of
      millions of rows; read replica for CRM and reporting.

---

## Known technical debt

- [ ] Local development rows written before `script_12` have `created_at`/`updated_at`
      3.5 hours late (one-off fix available; production unaffected).
- [ ] CRM ads search no longer matches customer names (only ID, mobile, title). Add a trigram
      index on `customer.name` if staff need it.
- [ ] CRM "content", "notifications", "reports", and "locations" pages are placeholders.
- [ ] `ads-platform-stats` and other services have no README.
