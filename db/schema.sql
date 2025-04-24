CREATE SCHEMA events;

SET search_path TO events, public;

CREATE TABLE events.actions (
    id SERIAL PRIMARY KEY,
    hostname TEXT UNIQUE NOT NULL,
    type INT NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_actions_hostname ON events.actions (hostname);
CREATE INDEX idx_actions_created_at ON events.actions (created_at);
