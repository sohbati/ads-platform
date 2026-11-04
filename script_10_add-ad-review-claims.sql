SET search_path TO ads_platform_schema;

-- Review claims: a reviewer takes one pending ad at a time, so two reviewers never work on the
-- same ad. A claim lasts until claim_expires_at; the panel extends it while the ad stays open,
-- and an abandoned claim simply expires and the ad returns to the queue.

ALTER TABLE ads_platform_schema.ads
    ADD COLUMN IF NOT EXISTS claimed_by_staff_id BIGINT REFERENCES ads_platform_schema.staff(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS claimed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS claim_expires_at TIMESTAMP;

COMMENT ON COLUMN ads_platform_schema.ads.claimed_by_staff_id IS 'بازبینی که آگهی را برای بررسی برداشته است';
COMMENT ON COLUMN ads_platform_schema.ads.claimed_at IS 'زمان برداشتن آگهی برای بررسی';
COMMENT ON COLUMN ads_platform_schema.ads.claim_expires_at IS 'پایان اعتبار برداشت؛ پس از آن آگهی به صف برمی‌گردد';

-- The queue (pending ads, oldest first) uses ads_pending_since_idx from script_1.
CREATE INDEX IF NOT EXISTS ads_claimed_by_staff_id_idx
    ON ads_platform_schema.ads (claimed_by_staff_id) WHERE claimed_by_staff_id IS NOT NULL;

-- Codes must match crm-ui/src/lib/permissions/codes.ts.
INSERT INTO ads_platform_schema.action (code, name) VALUES
    ('advertisement.review', 'بررسی آگهی (برداشتن از صف)'),
    ('advertisement.review_manage', 'مدیریت صف بررسی (آزادسازی آگهی دیگران)')
ON CONFLICT (code) DO NOTHING;

-- Give staff this role to make them reviewers.
INSERT INTO ads_platform_schema."role" (code, name) VALUES
    ('reviewer', 'بازبین آگهی')
ON CONFLICT (code) DO NOTHING;

INSERT INTO ads_platform_schema.role_action (role_id, action_id)
SELECT r.id, a.id
FROM ads_platform_schema."role" r
JOIN ads_platform_schema.action a ON a.code IN (
    'advertisement.read', 'advertisement.review', 'advertisement.approve', 'advertisement.reject'
)
WHERE r.code = 'reviewer'
ON CONFLICT DO NOTHING;
