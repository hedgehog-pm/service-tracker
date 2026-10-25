-- +goose Up
CREATE TABLE projects (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    key         TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    next_number INT  NOT NULL
);

CREATE TABLE statuses (
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       TEXT   NOT NULL,
    category   TEXT   NOT NULL,
    UNIQUE (id, project_id)
);

CREATE TABLE modules (
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name       TEXT   NOT NULL,
    UNIQUE (id, project_id)
);

CREATE TABLE tasks (
    id          INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id  INT      NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    number      INT         NOT NULL,
    module_id   INT,
    status_id   INT      NOT NULL,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL,
    UNIQUE (project_id, number),
    FOREIGN KEY (status_id, project_id) REFERENCES statuses(id, project_id),
    FOREIGN KEY (module_id, project_id) REFERENCES modules(id, project_id)
);
