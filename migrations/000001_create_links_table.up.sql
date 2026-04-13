CREATE TABLE IF NOT EXISTS links
(
    id           SERIAL PRIMARY KEY,
    short_code   VARCHAR UNIQUE NOT NULL,
    original_url VARCHAR        NOT NULL,
    created_at   TIMESTAMP      NOT NULL,
    visits       INT DEFAULT 0
);