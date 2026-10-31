CREATE SCHEMA IF NOT EXISTS ads_platform_schema;

SET search_path TO ads_platform_schema;

CREATE TABLE IF NOT EXISTS ads_platform_schema.customer (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    mobile VARCHAR(100) NOT NULL UNIQUE,
    national_id VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.customer IS 'ناشران آگهی (ورود OTP و مالکیت آگهی)';
COMMENT ON COLUMN ads_platform_schema.customer.id IS 'شناسه یکتای مشتری';
COMMENT ON COLUMN ads_platform_schema.customer.name IS 'نام نمایشی مشتری';
COMMENT ON COLUMN ads_platform_schema.customer.mobile IS 'شماره موبایل مشتری (یکتا؛ برای ورود و احراز هویت)';
COMMENT ON COLUMN ads_platform_schema.customer.national_id IS 'کد ملی مشتری (یکتا)';
COMMENT ON COLUMN ads_platform_schema.customer.created_at IS 'زمان ایجاد رکورد مشتری';
COMMENT ON COLUMN ads_platform_schema.customer.updated_at IS 'زمان آخرین به‌روزرسانی رکورد مشتری';

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    mobile VARCHAR(100) NOT NULL UNIQUE,
    national_id VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL DEFAULT '',
    is_super_admin BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.staff IS 'کارکنان پنل CRM؛ دسترسی فقط از طریق نقش';

-- فیلدهای اصلی فیلترپذیر به‌صورت ستون؛ ویژگی‌های وابسته به دسته‌بندی در JSONB
CREATE TABLE IF NOT EXISTS ads_platform_schema.ads (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES ads_platform_schema.customer(id),
    category_id INTEGER NOT NULL,
    city_id INTEGER NOT NULL,
    title VARCHAR(120) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    price_amount BIGINT,
    price_type VARCHAR(20) NOT NULL DEFAULT 'fixed',
    currency CHAR(3) NOT NULL DEFAULT 'IRR',
    attrs JSONB NOT NULL DEFAULT '{}'::jsonb,
    media JSONB NOT NULL DEFAULT '[]'::jsonb,
    contact JSONB NOT NULL DEFAULT '{}'::jsonb,
    location JSONB NOT NULL DEFAULT '{}'::jsonb,
    slug VARCHAR(160),
    published_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.ads IS 'آگهی‌های ثبت‌شده در پلتفرم (منبع اصلی داده؛ جستجوی پیشرفته می‌تواند در Elastic هم ایندکس شود)';
COMMENT ON COLUMN ads_platform_schema.ads.id IS 'شناسه یکتای آگهی';
COMMENT ON COLUMN ads_platform_schema.ads.customer_id IS 'شناسه ناشر آگهی (ارجاع به جدول customer)';
COMMENT ON COLUMN ads_platform_schema.ads.category_id IS 'شناسه دسته‌بندی برگ (مثلاً فروش مسکونی، خودرو)';
COMMENT ON COLUMN ads_platform_schema.ads.city_id IS 'شناسه شهر محل آگهی';
COMMENT ON COLUMN ads_platform_schema.ads.title IS 'عنوان کوتاه آگهی برای نمایش در لیست و جزئیات';
COMMENT ON COLUMN ads_platform_schema.ads.description IS 'توضیحات کامل آگهی';
COMMENT ON COLUMN ads_platform_schema.ads.status IS 'وضعیت آگهی: draft، pending، active، rejected، expired، deleted';
COMMENT ON COLUMN ads_platform_schema.ads.price_amount IS 'مبلغ قیمت به کوچک‌ترین واحد پول (مثلاً ریال)';
COMMENT ON COLUMN ads_platform_schema.ads.price_type IS 'نوع قیمت: fixed (قطعی)، negotiable (توافقی)، free (رایگان)، salary (حقوق)';
COMMENT ON COLUMN ads_platform_schema.ads.currency IS 'کد ارز سه‌حرفی ISO (پیش‌فرض IRR)';
COMMENT ON COLUMN ads_platform_schema.ads.attrs IS 'ویژگی‌های وابسته به دسته به‌صورت JSONB (مثلاً متراژ، اتاق، برند، کارکرد)';
COMMENT ON COLUMN ads_platform_schema.ads.media IS 'آرایه JSON تصاویر/رسانه آگهی (آدرس، بندانگشتی، ترتیب، تصویر اصلی)';
COMMENT ON COLUMN ads_platform_schema.ads.contact IS 'اطلاعات تماس قابل‌نمایش (شماره، چت، نمایش/مخفی بودن تلفن)';
COMMENT ON COLUMN ads_platform_schema.ads.location IS 'موقعیت مکانی اختیاری (مختصات، محله، نمایش عمومی آدرس)';
COMMENT ON COLUMN ads_platform_schema.ads.slug IS 'نامک یکتا برای URL دوستانه آگهی';
COMMENT ON COLUMN ads_platform_schema.ads.published_at IS 'زمان انتشار عمومی آگهی';
COMMENT ON COLUMN ads_platform_schema.ads.expires_at IS 'زمان انقضای آگهی؛ پس از آن از حالت active خارج می‌شود';
COMMENT ON COLUMN ads_platform_schema.ads.created_at IS 'زمان ایجاد رکورد آگهی';
COMMENT ON COLUMN ads_platform_schema.ads.updated_at IS 'زمان آخرین به‌روزرسانی رکورد آگهی';

CREATE INDEX IF NOT EXISTS ads_customer_id_idx ON ads_platform_schema.ads (customer_id);
CREATE INDEX IF NOT EXISTS ads_category_id_idx ON ads_platform_schema.ads (category_id);
CREATE INDEX IF NOT EXISTS ads_city_status_idx ON ads_platform_schema.ads (city_id, status);
CREATE INDEX IF NOT EXISTS ads_price_idx ON ads_platform_schema.ads (price_amount) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS ads_published_at_idx ON ads_platform_schema.ads (published_at DESC) WHERE status = 'active';
CREATE INDEX IF NOT EXISTS ads_attrs_gin_idx ON ads_platform_schema.ads USING GIN (attrs);
CREATE INDEX IF NOT EXISTS ads_media_gin_idx ON ads_platform_schema.ads USING GIN (media);

-- متادیتای تصاویر آگهی؛ فایل اصلی در object storage (MinIO/S3) نگهداری می‌شود
CREATE TABLE IF NOT EXISTS ads_platform_schema.ad_images (
    id BIGSERIAL PRIMARY KEY,
    customer_id BIGINT NOT NULL REFERENCES ads_platform_schema.customer(id),
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
COMMENT ON COLUMN ads_platform_schema.ad_images.status IS 'وضعیت تصویر: pending (در انتظار آپلود)، uploaded (آپلود شده)، deleted (حذف شده)';

CREATE INDEX IF NOT EXISTS ad_images_customer_id_idx ON ads_platform_schema.ad_images (customer_id);
CREATE INDEX IF NOT EXISTS ad_images_ad_id_idx ON ads_platform_schema.ad_images (ad_id) WHERE ad_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS ad_images_pending_created_idx ON ads_platform_schema.ad_images (created_at) WHERE status = 'pending';

CREATE TABLE IF NOT EXISTS ads_platform_schema.customer_profile (
    customer_id BIGINT PRIMARY KEY REFERENCES ads_platform_schema.customer(id),
    location_slugs JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ads_platform_schema.action (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ads_platform_schema."role" (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ads_platform_schema.role_action (
    role_id BIGINT NOT NULL REFERENCES ads_platform_schema."role"(id) ON DELETE CASCADE,
    action_id BIGINT NOT NULL REFERENCES ads_platform_schema.action(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, action_id)
);

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff_role (
    staff_id BIGINT NOT NULL REFERENCES ads_platform_schema.staff(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES ads_platform_schema."role"(id) ON DELETE CASCADE,
    PRIMARY KEY (staff_id, role_id)
);

CREATE INDEX IF NOT EXISTS staff_role_role_id_idx ON ads_platform_schema.staff_role (role_id);

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff_group (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff_group_member (
    group_id BIGINT NOT NULL REFERENCES ads_platform_schema.staff_group(id) ON DELETE CASCADE,
    staff_id BIGINT NOT NULL REFERENCES ads_platform_schema.staff(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, staff_id)
);

CREATE INDEX IF NOT EXISTS staff_group_member_staff_id_idx ON ads_platform_schema.staff_group_member (staff_id);
