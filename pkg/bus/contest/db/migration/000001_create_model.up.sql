CREATE TABLE contestant (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    surname TEXT NOT NULL,
    birthdate DATE NOT NULL,
    policy_accepted BOOLEAN NOT NULL CHECK (policy_accepted IN (0, 1))
);
CREATE TABLE contest (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    start_time TEXT NOT NULL,
    end_time TEXT NOT NULL,
    is_active BOOLEAN NOT NULL CHECK (is_active IN (0, 1))
);
CREATE TABLE entry (
    id TEXT PRIMARY KEY,
    contestant_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    status TEXT NOT NULL CHECK (
        status IN (
            'Pending',
            'Submitted',
            'ConfirmationEmailSent',
            'Confirmed'
        )
    ),
    FOREIGN KEY (contestant_id) REFERENCES contestant(id) ON DELETE CASCADE
);
CREATE TABLE art_piece (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    entry_id TEXT NOT NULL,
    key TEXT NOT NULL,
    created_at TEXT NOT NULL,
    foreign KEY (entry_id) REFERENCES entry(id) ON DELETE CASCADE
);