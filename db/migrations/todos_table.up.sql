-- Table: public.todos

CREATE TABLE IF NOT EXISTS public.todos
(
    id text COLLATE pg_catalog."default" NOT NULL,
    title text COLLATE pg_catalog."default" NOT NULL,
    completed boolean NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
    )

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.todos
    OWNER to rill;

CREATE INDEX IF NOT EXISTS idx_todos_deleted_at
    ON public.todos USING btree
    (deleted_at ASC NULLS LAST)
    TABLESPACE pg_default;
