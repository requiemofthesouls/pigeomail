-- +goose Up 
CREATE TABLE IF NOT EXISTS pigeomail.telegram_users
(
    id      BIGSERIAL NOT NULL,
    chat_id BIGINT    NOT NULL,
    email   TEXT      NOT NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX IF NOT EXISTS telegram_users_email_uindex
    ON pigeomail.telegram_users (email);

-- +goose Down 
DROP TABLE IF EXISTS pigeomail.telegram_users;
