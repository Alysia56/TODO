-- Filename: migrations/000003_add_todo_indexes.up.sql
CREATE INDEX IF NOT EXISTS todo_name_idx ON todo USING GIN(to_tsvector('simple', name));
CREATE INDEX IF NOT EXISTS todo_level_idx ON todo USING GIN(to_tsvector('simple', level));
CREATE INDEX IF NOT EXISTS todo_mode_idx ON todo USING GIN(mode);