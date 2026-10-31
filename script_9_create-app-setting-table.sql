SET search_path TO ads_platform_schema;

-- Runtime settings that managers change from the CRM (Settings → General) without a deploy.
-- Background jobs read the same rows. Values are text; the backend validates each known key.

CREATE TABLE IF NOT EXISTS ads_platform_schema.app_setting (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT NOT NULL,
    updated_by_staff_id BIGINT REFERENCES ads_platform_schema.staff(id) ON DELETE SET NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.app_setting IS 'تنظیمات قابل‌تغییر در زمان اجرا (از پنل CRM)';
COMMENT ON COLUMN ads_platform_schema.app_setting.key IS 'کلید پایدار تنظیم، مثلاً ads.pending_max_days';
COMMENT ON COLUMN ads_platform_schema.app_setting.value IS 'مقدار تنظیم به‌صورت متن';
COMMENT ON COLUMN ads_platform_schema.app_setting.updated_by_staff_id IS 'آخرین کارمندی که مقدار را تغییر داد';

INSERT INTO ads_platform_schema.app_setting (key, value) VALUES
    ('ads.pending_max_days', '7')
ON CONFLICT (key) DO NOTHING;
