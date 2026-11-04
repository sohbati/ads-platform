# Advertisements

Moderation UI for listings, backed by `/api/v1/panel/ads`.

- `/advertisements` finds ads of any status. It shows nothing until you search or pick a filter
  (or the Pending status), so the panel never pages through every ad on the platform.
- `/advertisements/review` is the review queue: pending ads only, oldest first.

The search box takes an ad ID, a customer mobile in any form (`0912…`, `+98912…`, Persian
digits), or 3+ letters of the title; it searches on Enter. Category includes its subcategories.
Each filter has a database index (`script_11_add-ad-search-indexes.sql`). Lists page with
Previous/Next cursors instead of page numbers, and the total stops counting at 1000 ("1000+").

Open an ad to see its content, photos, attributes and last review. Rejecting needs a reason for
the customer and takes an optional internal note; the ad returns to the customer's drafts, and
goes back to pending once they edit it. Photo URLs use `NEXT_PUBLIC_MEDIA_BASE_URL`.

## Review claims

Only one reviewer works on an ad at a time. Staff are reviewers when one of their roles grants
`advertisement.review` (the seeded `reviewer` role does; see `script_10_add-ad-review-claims.sql`).

- **Next ad** on the review page claims the oldest waiting ad, or reopens the one the reviewer
  already holds. After approving or rejecting, the next ad opens on its own.
- **Review this ad** in the drawer claims a specific pending ad that nobody else holds.
- Approve (`advertisement.approve`) and reject (`advertisement.reject`) work only on an ad the
  reviewer holds.
- A claim lasts 10 minutes. The open drawer renews it every minute, up to 30 minutes in total; the
  reviewer is warned a minute before it ends and can claim again to continue.
- Closing the drawer, **Return to queue**, and signing out hand the ad back at once.
- Staff with `advertisement.review_manage` can free an ad someone else holds, and see everyone's
  counts for today on the review page.
