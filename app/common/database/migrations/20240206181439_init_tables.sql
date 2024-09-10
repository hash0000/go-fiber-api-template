-- +goose Up
CREATE TABLE "user" (
    id bigint NOT NULL,
    token_number smallint NOT NULL DEFAULT 0,
    use_trial boolean NOT NULL DEFAULT false,
    invite_code UUID NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

CREATE TABLE tale (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR NOT NULL,
    file_name VARCHAR NOT NULL,
    is_payed boolean NOT NULL DEFAULT false,
    tale_generation_id VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT book_user_fkey FOREIGN KEY (user_id) REFERENCES "user" (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE "user" CASCADE;

DROP TABLE "tale" CASCADE;