SET search_path TO ads_platform_schema;

-- قالب JSON Schema برای ویژگی‌های وابسته به دسته (ستون ads.attrs)
-- نام قالب با adsAttrsJsonSchemaTemplateName در category.json متناظر است
CREATE TABLE ads_platform_schema.ads_attrs_json_schema_template (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    json_schema JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.ads_attrs_json_schema_template IS 'قالب JSON Schema برای فرم ویژگی‌های آگهی (attrs)؛ هر دسته می‌تواند با adsAttrsJsonSchemaTemplateName به یک قالب اشاره کند';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.id IS 'شناسه یکتای قالب';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.name IS 'نام یکتای قالب؛ همان مقداری که در category.json در adsAttrsJsonSchemaTemplateName ذخیره می‌شود';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.title IS 'عنوان نمایشی قالب برای ادمین و UI';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.description IS 'توضیح اختیاری دربارهٔ کاربرد قالب';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.json_schema IS 'سند JSON Schema (draft) که فیلدهای attrs و اعتبارسنجی آن‌ها را تعریف می‌کند';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.created_at IS 'زمان ایجاد قالب';
COMMENT ON COLUMN ads_platform_schema.ads_attrs_json_schema_template.updated_at IS 'زمان آخرین به‌روزرسانی قالب';
