SET search_path TO ads_platform_schema;

CREATE TABLE IF NOT EXISTS ads_platform_schema.ad_events (
    id               BIGSERIAL PRIMARY KEY,
    ad_id            BIGINT NOT NULL REFERENCES ads_platform_schema.ads(id),
    event            VARCHAR(32) NOT NULL,
    viewer_id        UUID NOT NULL,
    session_user_id  BIGINT,
    occurred_at      TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS ad_events_ad_occurred_at
    ON ads_platform_schema.ad_events (ad_id, occurred_at);

CREATE TABLE IF NOT EXISTS ads_platform_schema.ad_stats_daily (
    ad_id            BIGINT NOT NULL REFERENCES ads_platform_schema.ads(id),
    day              DATE NOT NULL,
    views            INTEGER NOT NULL DEFAULT 0,
    unique_viewers   INTEGER NOT NULL DEFAULT 0,
    contact_reveals  INTEGER NOT NULL DEFAULT 0,
    calls            INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (ad_id, day)
);

CREATE TABLE IF NOT EXISTS ads_platform_schema.ad_view_uniques (
    ad_id     BIGINT NOT NULL REFERENCES ads_platform_schema.ads(id),
    viewer_id UUID NOT NULL,
    day       DATE NOT NULL,
    PRIMARY KEY (ad_id, viewer_id, day)
);

COMMENT ON TABLE ads_platform_schema.ad_events IS 'Append-only ad view/contact events for rebuild and debug';
COMMENT ON TABLE ads_platform_schema.ad_stats_daily IS 'Per-ad daily rollups for seller dashboards';
COMMENT ON TABLE ads_platform_schema.ad_view_uniques IS 'One unique viewer per ad per Tehran calendar day';
