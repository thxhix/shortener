CREATE TABLE IF NOT EXISTS shortener (
    id SERIAL PRIMARY KEY,
    original VARCHAR(512) UNIQUE NOT NULL,
    shorten VARCHAR(10) UNIQUE NOT NULL,
    user_id UUID,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_shortener_user_id ON shortener(user_id);
CREATE INDEX IF NOT EXISTS idx_shortener_is_deleted ON shortener(is_deleted);
