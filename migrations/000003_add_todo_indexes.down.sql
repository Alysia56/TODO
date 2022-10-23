--Filename: migrations/000003_add_Todo_indexes.down.sql
DROP INDEX IF EXISTS todo_name_idx;
DROP INDEX IF EXISTS todo_level_idx;
DROP INDEX IF EXISTS todo_mode_idx;