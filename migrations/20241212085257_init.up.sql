BEGIN;

CREATE TABLE IF NOT EXISTS links(
    id     UUID PRIMARY KEY,
    alias  TEXT UNIQUE,
    url    TEXT,

    expired_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_links_url ON links USING hash (url);
CREATE INDEX IF NOT EXISTS idx_links_alias ON links USING hash (alias);

COMMIT;