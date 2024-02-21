-- +goose Up
CREATE TABLE "user" (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    username character varying(127) NOT NULL,
    password character varying(255) NOT NULL,
    "grandAccess" boolean NOT NULL DEFAULT false,
    "createdAt" timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

CREATE TABLE "url" (
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    url character varying(127) NOT NULL,
    "updatedAt" timestamp with time zone NOT NULL DEFAULT now(),
    "createdAt" timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT url_pkey PRIMARY KEY (id)
);

CREATE UNIQUE INDEX user_username_key ON "user" USING btree (username);

-- +goose Down
DROP TABLE "user";
DROP TABLE "url";