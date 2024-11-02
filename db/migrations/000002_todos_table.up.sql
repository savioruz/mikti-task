-- Table: public.todos

CREATE TABLE IF NOT EXISTS todos (
    id varchar(36) COLLATE pg_catalog."default" NOT NULL,
    title varchar(255) COLLATE pg_catalog."default" NOT NULL,
    completed boolean NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT todos_pkey PRIMARY KEY (id)
    );

CREATE INDEX IF NOT EXISTS idx_todos_deleted_at
    ON todos USING btree
    (deleted_at ASC NULLS LAST)
    TABLESPACE pg_default;
