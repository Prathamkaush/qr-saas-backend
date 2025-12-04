CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    avatar_url TEXT,
    provider TEXT NOT NULL,
    provider_id TEXT,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now()
);
