-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS config_templates;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS traffic_logs;
DROP TABLE IF EXISTS clients;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
