-- +goose Up
-- +goose StatementBegin
ALTER TABLE "user" ALTER COLUMN use_trial SET DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "user" ALTER COLUMN use_trial SET DEFAULT true;
-- +goose StatementEnd
