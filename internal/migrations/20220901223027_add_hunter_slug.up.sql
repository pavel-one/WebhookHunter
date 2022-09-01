ALTER TABLE hunters ADD COLUMN slug VARCHAR(255) CONSTRAINT unique_hunters_slug UNIQUE;

CREATE INDEX idx_hunters_slug ON hunters (slug)