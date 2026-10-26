-- +goose Up
CREATE TABLE projects (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key         TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    next_number INT NOT NULL
);

CREATE TABLE modules (
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    UNIQUE (id, project_id)
);
