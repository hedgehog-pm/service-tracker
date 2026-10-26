-- +goose Up
ALTER TABLE tasks
    ADD COLUMN estimate_seconds INT,
    ADD COLUMN due_date DATE;
