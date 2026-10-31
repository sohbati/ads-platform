SET search_path TO ads_platform_schema;

-- IAM: action catalog, roles, staff-role assignment.
-- staff_group is only a label for categorizing staff (does not grant actions).

CREATE TABLE IF NOT EXISTS ads_platform_schema.action (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.action IS 'عمل مجاز در پنل (مثلاً ad.approve)';
COMMENT ON COLUMN ads_platform_schema.action.code IS 'شناسه پایدار عمل، به شکل resource.verb';
COMMENT ON COLUMN ads_platform_schema.action.name IS 'نام نمایشی عمل';

CREATE TABLE IF NOT EXISTS ads_platform_schema."role" (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema."role" IS 'نقش شغلی (مثلاً moderator)';
COMMENT ON COLUMN ads_platform_schema."role".code IS 'شناسه پایدار نقش';
COMMENT ON COLUMN ads_platform_schema."role".name IS 'نام نمایشی نقش';

CREATE TABLE IF NOT EXISTS ads_platform_schema.role_action (
    role_id BIGINT NOT NULL REFERENCES ads_platform_schema."role"(id) ON DELETE CASCADE,
    action_id BIGINT NOT NULL REFERENCES ads_platform_schema.action(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, action_id)
);

COMMENT ON TABLE ads_platform_schema.role_action IS 'اعمال متعلق به یک نقش';

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff_role (
    staff_id BIGINT NOT NULL REFERENCES ads_platform_schema.staff(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES ads_platform_schema."role"(id) ON DELETE CASCADE,
    PRIMARY KEY (staff_id, role_id)
);

COMMENT ON TABLE ads_platform_schema.staff_role IS 'نقش‌های مستقیم کارکنان؛ تنها منبع دسترسی پنل';

CREATE INDEX IF NOT EXISTS staff_role_role_id_idx ON ads_platform_schema.staff_role (role_id);

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff_group (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.staff_group IS 'دسته‌بندی کارکنان؛ نقش یا عمل نمی‌دهد';

CREATE TABLE IF NOT EXISTS ads_platform_schema.staff_group_member (
    group_id BIGINT NOT NULL REFERENCES ads_platform_schema.staff_group(id) ON DELETE CASCADE,
    staff_id BIGINT NOT NULL REFERENCES ads_platform_schema.staff(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, staff_id)
);

COMMENT ON TABLE ads_platform_schema.staff_group_member IS 'عضویت کارکنان در یک گروه دسته‌بندی';

CREATE INDEX IF NOT EXISTS staff_group_member_staff_id_idx ON ads_platform_schema.staff_group_member (staff_id);

-- Panel action catalog. Codes must match crm-ui/src/lib/permissions/codes.ts.
INSERT INTO ads_platform_schema.action (code, name) VALUES
    ('customer.read', 'مشاهده مشتریان'),
    ('customer.update', 'ویرایش مشتریان'),
    ('customer.block', 'مسدودسازی مشتریان'),
    ('advertisement.read', 'مشاهده آگهی‌ها'),
    ('advertisement.approve', 'تأیید آگهی'),
    ('advertisement.reject', 'رد آگهی'),
    ('advertisement.delete', 'حذف آگهی'),
    ('category.read', 'مشاهده دسته‌بندی‌ها'),
    ('category.create', 'ایجاد دسته‌بندی'),
    ('category.update', 'ویرایش دسته‌بندی'),
    ('category.delete', 'حذف دسته‌بندی'),
    ('location.read', 'مشاهده مکان‌ها'),
    ('report.read', 'مشاهده گزارش‌ها'),
    ('content.read', 'مشاهده محتوا'),
    ('content.update', 'ویرایش محتوا'),
    ('notification.read', 'مشاهده اعلان‌ها'),
    ('notification.send', 'ارسال اعلان'),
    ('staff.read', 'مشاهده کارکنان'),
    ('staff.update', 'مدیریت کارکنان'),
    ('role.read', 'مشاهده نقش‌ها'),
    ('role.update', 'مدیریت نقش‌ها'),
    ('action.read', 'مشاهده عملیات'),
    ('action.update', 'مدیریت عملیات'),
    ('staff_group.read', 'مشاهده گروه‌های کارکنان'),
    ('staff_group.update', 'مدیریت گروه‌های کارکنان'),
    ('settings.read', 'مشاهده تنظیمات'),
    ('settings.update', 'ویرایش تنظیمات')
ON CONFLICT (code) DO NOTHING;
