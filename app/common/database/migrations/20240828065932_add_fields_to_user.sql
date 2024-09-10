-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    "user"
    ADD
        COLUMN first_name        varchar(255),
    ADD
        COLUMN last_name         varchar(255),
    ADD
        COLUMN telegram_username varchar(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    tale
    DROP COLUMN first_name,
    DROP COLUMN last_name,
    DROP COLUMN telegram_username;
-- +goose StatementEnd
