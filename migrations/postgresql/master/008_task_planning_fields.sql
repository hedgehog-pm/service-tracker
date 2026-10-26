-- +goose Up
ALTER TABLE tasks
    ADD COLUMN parent_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    ADD COLUMN estimate_seconds INT,
    ADD COLUMN due_date DATE;
