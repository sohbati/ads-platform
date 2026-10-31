SET search_path TO ads_platform_schema;

-- Default CRM super admin for a fresh install (full access to every panel action).
--
--   Mobile:   09120000000
--   Password: Ruab@Admin1405
--
-- Change this password right after the first sign-in (Staff page → Set password), or reset it with:
--   cd ads-platform-back && go run ./cmd/create-admin -mobile 09120000000
--
-- password_hash is bcrypt (cost 12) of the password above. Mobile is stored in E.164 form, as the
-- backend normalizes it. Re-running this script never overwrites an existing admin or their password.

INSERT INTO ads_platform_schema.staff (name, mobile, national_id, password_hash, is_super_admin, is_active)
VALUES (
    'مدیر سیستم',
    '+989120000000',
    '0000000000',
    '$2a$12$ohZcn4l5eFmpX37LjQUz0efe7XeW9mQCNl0nCvuZCaFO/7ubfr2v2',
    TRUE,
    TRUE
)
ON CONFLICT (mobile) DO NOTHING;
