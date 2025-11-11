CREATE TABLE IF NOT EXISTS articles (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    content     TEXT NOT NULL,
    author_id   INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    published   BOOLEAN NOT NULL DEFAULT false,
    views       INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_published ON articles(published);

-- SELECT last_value FROM articles_id_seq;

-- ALTER SEQUENCE articles_id_seq RESTART WITH 1;

-- SELECT setval('articles_id_seq', COALESCE((SELECT MAX(id) FROM articles), 1));