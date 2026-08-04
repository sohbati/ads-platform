CREATE ROLE ads_platform_user WITH LOGIN PASSWORD 'ads_platform_user';
CREATE DATABASE ads_platform OWNER ads_platform_user;
GRANT ALL PRIVILEGES ON DATABASE ads_platform TO ads_platform_user;
CREATE SCHEMA ads_platform_schema AUTHORIZATION ads_platform_user;
SET search_path TO ads_platform_schema;


CREATE TABLE ads_platform_schema.user (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    mobile VARCHAR(100) NOT NULL UNIQUE,
    national_id VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

