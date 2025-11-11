CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY, 
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

-- SELECT last_value FROM users_id_seq;

-- -- Сбросить на 1
-- ALTER SEQUENCE users_id_seq RESTART WITH 1;