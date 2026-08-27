SET search_path TO ads_platform_schema;

-- متادیتای تصاویر آگهی؛ فایل اصلی در object storage (MinIO/S3) نگهداری می‌شود
CREATE TABLE ads_platform_schema.ad_images (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES ads_platform_schema."user"(id),
    ad_id BIGINT REFERENCES ads_platform_schema.ads(id),
    object_key VARCHAR(255) NOT NULL UNIQUE,
    original_filename VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    checksum VARCHAR(64),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uploaded_at TIMESTAMP,
    deleted_at TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.ad_images IS 'متادیتای تصاویر آگهی؛ خود فایل در object storage ذخیره می‌شود و این جدول چرخهٔ عمر آن را دنبال می‌کند';
COMMENT ON COLUMN ads_platform_schema.ad_images.id IS 'شناسه یکتای تصویر';
COMMENT ON COLUMN ads_platform_schema.ad_images.user_id IS 'شناسه کاربر آپلودکننده (ارجاع به جدول user)';
COMMENT ON COLUMN ads_platform_schema.ad_images.ad_id IS 'شناسه آگهی مرتبط؛ تا قبل از ثبت آگهی NULL است';
COMMENT ON COLUMN ads_platform_schema.ad_images.object_key IS 'کلید یکتای فایل در object storage (مثلاً ads/1/1_1.jpg)';
COMMENT ON COLUMN ads_platform_schema.ad_images.original_filename IS 'نام اصلی فایل هنگام آپلود توسط کاربر';
COMMENT ON COLUMN ads_platform_schema.ad_images.content_type IS 'نوع MIME فایل: image/jpeg، image/png، image/webp';
COMMENT ON COLUMN ads_platform_schema.ad_images.file_size IS 'حجم فایل به بایت';
COMMENT ON COLUMN ads_platform_schema.ad_images.status IS 'وضعیت تصویر: pending (در انتظار آپلود)، uploaded (آپلود شده)، deleted (حذف شده)';
COMMENT ON COLUMN ads_platform_schema.ad_images.checksum IS 'اثر انگشت محتوای فایل (مثلاً SHA-256) برای اعتبارسنجی آپلود';
COMMENT ON COLUMN ads_platform_schema.ad_images.created_at IS 'زمان ایجاد رکورد متادیتا';
COMMENT ON COLUMN ads_platform_schema.ad_images.uploaded_at IS 'زمان تکمیل آپلود فایل در object storage';
COMMENT ON COLUMN ads_platform_schema.ad_images.deleted_at IS 'زمان حذف منطقی تصویر؛ فایل بعداً از storage پاک می‌شود';

CREATE INDEX ad_images_user_id_idx ON ads_platform_schema.ad_images (user_id);
CREATE INDEX ad_images_ad_id_idx ON ads_platform_schema.ad_images (ad_id) WHERE ad_id IS NOT NULL;
CREATE INDEX ad_images_pending_created_idx ON ads_platform_schema.ad_images (created_at) WHERE status = 'pending';
