-- +goose Up
ALTER TABLE "user" ALTER COLUMN use_trial SET DEFAULT true;

-- +goose Down
ALTER TABLE "user" ALTER COLUMN use_trial SET DEFAULT false;
