-- +goose Up
CREATE TABLE task_worklogs (
    id                 INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_id            INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    started_at         TIMESTAMPTZ NOT NULL,
    time_spent_seconds INT NOT NULL CHECK (time_spent_seconds > 0),
    comment            TEXT,
    created_at         TIMESTAMPTZ NOT NULL,
    updated_at         TIMESTAMPTZ NOT NULL
);
