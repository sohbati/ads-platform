SET search_path TO ads_platform_schema;

CREATE TABLE IF NOT EXISTS ads_platform_schema.customer_profile (
    customer_id BIGINT PRIMARY KEY REFERENCES ads_platform_schema.customer(id),
    location_slugs JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.customer_profile IS 'ترجیحات حساب مشتری (مثلاً شهرهای انتخاب‌شده برای جستجو)';
COMMENT ON COLUMN ads_platform_schema.customer_profile.customer_id IS 'شناسه مشتری (یک ردیف به ازای هر حساب)';
COMMENT ON COLUMN ads_platform_schema.customer_profile.location_slugs IS 'آرایهٔ JSON از اسلاگ استان/شهر انتخاب‌شده';
COMMENT ON COLUMN ads_platform_schema.customer_profile.created_at IS 'زمان ایجاد پروفایل';
COMMENT ON COLUMN ads_platform_schema.customer_profile.updated_at IS 'زمان آخرین به‌روزرسانی پروفایل';
