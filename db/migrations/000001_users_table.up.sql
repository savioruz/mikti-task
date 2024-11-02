-- Table: public.users

CREATE TABLE IF NOT EXISTS users (
    id varchar(36) NOT NULL,
    email varchar(100) NOT NULL,
    password varchar(255) NOT NULL,
    role varchar(5) NOT NULL,
    status boolean NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
    );

ALTER TABLE users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE users
    ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'user'));

CREATE INDEX idx_users_deleted_at
    ON users USING btree
    (deleted_at ASC NULLS LAST);
