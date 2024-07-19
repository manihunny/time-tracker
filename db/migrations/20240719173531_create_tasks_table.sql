-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks
(
    id integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    user_id integer NOT NULL,
    title text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    started_at timestamp without time zone NOT NULL,
    finished_at timestamp without time zone,
    CONSTRAINT tasks_pkey PRIMARY KEY (id),
    CONSTRAINT tasks_user_id_fkey FOREIGN KEY (user_id)
        REFERENCES public.people (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.tasks;
-- +goose StatementEnd
