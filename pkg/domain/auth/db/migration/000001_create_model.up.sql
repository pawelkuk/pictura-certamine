CREATE TABLE session (
    token TEXT PRIMARY KEY,
    expiry REAL NOT NULL,
    user_id TEXT,
    FOREIGN KEY(user_id) REFERENCES user(id) ON DELETE CASCADE
);
CREATE INDEX sessions_expiry_idx ON session(expiry);