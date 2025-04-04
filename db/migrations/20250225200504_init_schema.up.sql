CREATE TABLE IF NOT EXISTS users (
    id SERIAL primary key,
    email varchar(255) NOT NULL,
    password_hash varchar(255) NOT NULL,
    is_admin bool NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL primary key,
    user_id INTEGER,
    title text,
    description text,
    completed bool NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id)
    );