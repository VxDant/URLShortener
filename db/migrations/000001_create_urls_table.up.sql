CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    clicks INT DEFAULT 0
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
CREATE UNIQUE INDEX IF NOT EXISTS idx_long_url ON urls(long_url);