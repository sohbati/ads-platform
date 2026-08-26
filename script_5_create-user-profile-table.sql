SET search_path TO ads_platform_schema;

CREATE TABLE IF NOT EXISTS ads_platform_schema.user_profile (
    user_id BIGINT PRIMARY KEY REFERENCES ads_platform_schema."user"(id),
    location_slugs JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.user_profile IS 'ترجیحات حساب کاربر (مثلاً شهرهای انتخاب‌شده برای جستجو)';
COMMENT ON COLUMN ads_platform_schema.user_profile.user_id IS 'شناسه کاربر (یک ردیف به ازای هر حساب)';
COMMENT ON COLUMN ads_platform_schema.user_profile.location_slugs IS 'آرایهٔ JSON از اسلاگ استان/شهر انتخاب‌شده';
COMMENT ON COLUMN ads_platform_schema.user_profile.created_at IS 'زمان ایجاد پروفایل';
COMMENT ON COLUMN ads_platform_schema.user_profile.updated_at IS 'زمان آخرین به‌روزرسانی پروفایل';
