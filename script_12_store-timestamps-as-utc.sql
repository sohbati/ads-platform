SET search_path TO ads_platform_schema;

-- Store every timestamp as an absolute moment (TIMESTAMPTZ) instead of a bare wall-clock time.
--
-- TIMESTAMP keeps whatever clock time it is given. The backend wrote some columns in UTC
-- (pending_since, reviewed_at, …) and let GORM fill others in the server's local time
-- (created_at, updated_at), then read all of them back as UTC, so dates were off by the server's
-- UTC offset. TIMESTAMPTZ stores the moment itself; clients choose the zone to show it in.
--
-- Existing values are read as UTC. On a machine whose clock is not UTC, rows written before this
-- script keep their created_at/updated_at off by that machine's offset (Tehran: +3:30).
--
-- New tables must use TIMESTAMPTZ.

DO $$
DECLARE
    col record;
BEGIN
    FOR col IN
        SELECT table_name, column_name
        FROM information_schema.columns
        WHERE table_schema = 'ads_platform_schema'
          AND data_type = 'timestamp without time zone'
        ORDER BY table_name, column_name
    LOOP
        EXECUTE format(
            'ALTER TABLE ads_platform_schema.%I ALTER COLUMN %I TYPE TIMESTAMPTZ USING %I AT TIME ZONE ''UTC''',
            col.table_name, col.column_name, col.column_name
        );
        RAISE NOTICE 'converted %.%', col.table_name, col.column_name;
    END LOOP;
END
$$;
