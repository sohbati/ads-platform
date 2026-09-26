SET search_path TO ads_platform_schema;

-- IAM: action catalog, roles, user-role assignment.
-- user_group is only a label for categorizing users (does not grant actions).

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

CREATE TABLE IF NOT EXISTS ads_platform_schema.user_role (
    user_id BIGINT NOT NULL REFERENCES ads_platform_schema."user"(id) ON DELETE CASCADE,
    role_id BIGINT NOT NULL REFERENCES ads_platform_schema."role"(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

COMMENT ON TABLE ads_platform_schema.user_role IS 'نقش‌های مستقیم کاربر؛ تنها منبع دسترسی';

CREATE INDEX IF NOT EXISTS user_role_role_id_idx ON ads_platform_schema.user_role (role_id);

CREATE TABLE IF NOT EXISTS ads_platform_schema.user_group (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(80) NOT NULL UNIQUE,
    name VARCHAR(120) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ads_platform_schema.user_group IS 'دسته‌بندی کاربران؛ نقش یا عمل نمی‌دهد';

CREATE TABLE IF NOT EXISTS ads_platform_schema.user_group_member (
    group_id BIGINT NOT NULL REFERENCES ads_platform_schema.user_group(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES ads_platform_schema."user"(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, user_id)
);

COMMENT ON TABLE ads_platform_schema.user_group_member IS 'عضویت کاربر در یک گروه دسته‌بندی';

CREATE INDEX IF NOT EXISTS user_group_member_user_id_idx ON ads_platform_schema.user_group_member (user_id);
