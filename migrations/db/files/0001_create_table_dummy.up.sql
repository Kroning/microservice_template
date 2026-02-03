CREATE TABLE IF NOT EXISTS dummy
(
    id         BIGSERIAL       PRIMARY KEY,
    name        VARCHAR(255)    NOT NULL,
    created_at  TIMESTAMP       NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP       NOT NULL DEFAULT NOW()
);