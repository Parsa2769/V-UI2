-- Drop trigger and function
DROP TRIGGER IF EXISTS inbounds_updated_at ON inbounds;
DROP FUNCTION IF EXISTS update_inbounds_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS idx_inbounds_deleted_at;
DROP INDEX IF EXISTS idx_inbounds_enable;
DROP INDEX IF EXISTS idx_inbounds_port;
DROP INDEX IF EXISTS idx_inbounds_tag;
DROP INDEX IF EXISTS idx_inbounds_user_id;

-- Drop table
DROP TABLE IF EXISTS inbounds;
