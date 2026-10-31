# Advertisements

Moderation UI for listings, backed by `/api/v1/panel/ads`.

- `/advertisements` lists every status, with a status filter.
- `/advertisements/review` is the review queue: pending ads only, oldest first.

Open an ad to see its content, photos, attributes and last review. Pending ads can be approved
(`advertisement.approve`) or rejected (`advertisement.reject`). Rejecting needs a reason for the
customer and takes an optional internal note; the ad returns to the customer's drafts, and goes
back to pending once they edit it. Photo URLs use `NEXT_PUBLIC_MEDIA_BASE_URL`.
