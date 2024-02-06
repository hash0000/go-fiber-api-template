-- +goose Up
CREATE TABLE "user"
(
    id          uuid                     NOT NULL DEFAULT gen_random_uuid(),
    name        character varying(127)   NOT NULL,
    phone       character varying(255)   NOT NULL,
    "createdAt" timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX user_name_key
    ON "user" USING btree
    (name);

-- +goose Down
DROP TABLE "user";
~