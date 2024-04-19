CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY,
    name TEXT,
    account_type TEXT,
    login TEXT,
    password TEXT,
    email TEXT,
    email_password TEXT,
    recovery_email TEXT,
    recovery_email_password TEXT,
    cookie TEXT,
    status TEXT,
    created_at TIMESTAMP
);