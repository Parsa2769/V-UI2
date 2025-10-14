-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create inbounds table
CREATE TABLE IF NOT EXISTS inbounds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag VARCHAR(255) NOT NULL UNIQUE,
    protocol VARCHAR(50) NOT NULL,
    port INTEGER NOT NULL,
    listen VARCHAR(255) DEFAULT '0.0.0.0',
    enable BOOLEAN DEFAULT TRUE,
    settings JSONB NOT NULL DEFAULT '{}',
    stream_settings JSONB NOT NULL DEFAULT '{}',
    sniffing JSONB NOT NULL DEFAULT '{}',
    allocate JSONB NOT NULL DEFAULT '{}',
    remark TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_inbounds_user_id ON inbounds(user_id);
CREATE INDEX IF NOT EXISTS idx_inbounds_tag ON inbounds(tag);
CREATE INDEX IF NOT EXISTS idx_inbounds_port ON inbounds(port);
CREATE INDEX IF NOT EXISTS idx_inbounds_enable ON inbounds(enable);
CREATE INDEX IF NOT EXISTS idx_inbounds_deleted_at ON inbounds(deleted_at);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_inbounds_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER inbounds_updated_at
BEFORE UPDATE ON inbounds
FOR EACH ROW
EXECUTE FUNCTION update_inbounds_updated_at();
