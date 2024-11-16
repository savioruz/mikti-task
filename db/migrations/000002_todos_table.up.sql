-- Table: public.todos

CREATE TABLE IF NOT EXISTS todos (
    id varchar(36) COLLATE pg_catalog."default" NOT NULL,
    title varchar(255) COLLATE pg_catalog."default" NOT NULL,
    done boolean NOT NULL,
    user_id varchar(36) NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT todos_pkey PRIMARY KEY (id),
    CONSTRAINT fk_todos_user FOREIGN KEY (user_id)
        REFERENCES users (id) ON DELETE CASCADE
    );

CREATE INDEX IF NOT EXISTS idx_todos_deleted_at
    ON todos USING btree
    (deleted_at ASC NULLS LAST)
    TABLESPACE pg_default;
