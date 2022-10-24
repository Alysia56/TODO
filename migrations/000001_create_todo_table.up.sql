-- Filename: migrations/000001_create_todo_table.up.sq1

CREATE TABLE IF NOT EXISTS todo (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    title text NOT NULL, 
    label text NOT NULL,
    task text NOT NULL,
    status text NOT NULL,
    priority text NOT NULL,
    website text NOT NULL,
    address text NOT NULL,
    mode text[] NOT NULL,
    version integer NOT NULL DEFAULT 1
);