CREATE SCHEMA IF NOT EXISTS origin_meta;

CREATE TABLE IF NOT EXISTS origin_meta.assets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id text NOT NULL,
  key text NOT NULL,
  bucket text NOT NULL DEFAULT 'cdn',
  size bigint,
  content_type text,
  etag text,
  cache_ttl int DEFAULT 300,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_assets_user_id ON origin_meta.assets (user_id);
CREATE INDEX IF NOT EXISTS idx_assets_key ON origin_meta.assets (key);

CREATE OR REPLACE FUNCTION origin_meta.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON origin_meta.assets
FOR EACH ROW
EXECUTE PROCEDURE origin_meta.set_updated_at();
