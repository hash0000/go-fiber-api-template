-- +goose Up
-- +goose StatementBegin
ALTER TABLE
    "user"
    ADD
        COLUMN is_payed_tale boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE
    tale DROP COLUMN is_payed_tale;
-- +goose StatementEnd
