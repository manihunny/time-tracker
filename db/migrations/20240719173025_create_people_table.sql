-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.people
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    surname text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    name text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    patronymic text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    address text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    passport_number text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    CONSTRAINT people_pkey PRIMARY KEY (id)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.people;
-- +goose StatementEnd
