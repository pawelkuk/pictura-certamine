CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    authorization_token TEXT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL CHECK (is_active IN (0, 1)),
    activation_token TEXT NOT NULL,
    password_reset_token TEXT NOT NULL
);