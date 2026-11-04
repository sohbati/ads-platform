SET search_path TO ads_platform_schema;

-- Indexes behind the CRM ads list, so it stays fast with millions of ads. The list never pages
-- through every ad: outside the pending queue it needs an ad ID, a customer mobile, a title search,
-- a category, or a city, and each of those has an index.
--
-- On a large live table, run each CREATE INDEX as CREATE INDEX CONCURRENTLY, one statement at a
-- time and outside a transaction, so writes to ads are not blocked while it builds.

-- Title search: ILIKE '%term%' uses trigrams. Works for Persian text too.
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS ads_title_trgm_idx
    ON ads_platform_schema.ads USING GIN (title gin_trgm_ops);

-- Lists by status, newest activity first, paged by (updated_at, id).
CREATE INDEX IF NOT EXISTS ads_status_updated_idx
    ON ads_platform_schema.ads (status, updated_at DESC, id DESC);

-- Today's per-reviewer counts on the review page.
CREATE INDEX IF NOT EXISTS ads_reviewed_at_idx
    ON ads_platform_schema.ads (reviewed_at) WHERE reviewed_at IS NOT NULL;

-- An early version of script_10 created this duplicate of ads_pending_since_idx.
DROP INDEX IF EXISTS ads_platform_schema.ads_review_queue_idx;
