CREATE TABLE users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    fullname TEXT NOT NULL,
    birthdate TEXT NOT NULL,
    avatarname TEXT,
    avatar BLOB,
    nickname TEXT UNIQUE,
    aboutme TEXT,
    privacy INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    secret_key TEXT
);
