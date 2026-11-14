-- +goose Up
CREATE TABLE boards (
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       TEXT NOT NULL
);
