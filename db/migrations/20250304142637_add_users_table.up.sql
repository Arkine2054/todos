CREATE TABLE IF NOT EXISTS users (
    id SERIAL primary key,
    email varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    is_admin bool NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );