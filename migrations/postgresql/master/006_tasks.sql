-- +goose Up
-- Тип задачи и приоритет должны входить в схемы её проекта.
CREATE TABLE tasks (
    id                 INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    project_id         INT NOT NULL,
    number             INT NOT NULL,
    module_id          INT,
    issue_type_id      INT NOT NULL,
    workflow_scheme_id INT NOT NULL,
    workflow_id        INT NOT NULL,
    status_id          INT NOT NULL,
    priority_scheme_id INT NOT NULL,
    priority_id        INT NOT NULL,
    title              TEXT NOT NULL,
    description        TEXT NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL,
    updated_at         TIMESTAMPTZ NOT NULL,
    UNIQUE (project_id, number),
    FOREIGN KEY (project_id, workflow_scheme_id, priority_scheme_id)
        REFERENCES projects(id, workflow_scheme_id, priority_scheme_id) ON DELETE CASCADE,
    FOREIGN KEY (workflow_scheme_id, issue_type_id, workflow_id)
        REFERENCES workflow_scheme_rules(scheme_id, issue_type_id, workflow_id),
    FOREIGN KEY (workflow_id, status_id) REFERENCES workflow_statuses(workflow_id, status_id),
    FOREIGN KEY (priority_scheme_id, priority_id)
        REFERENCES priority_scheme_priorities(scheme_id, priority_id),
    FOREIGN KEY (module_id, project_id) REFERENCES modules(id, project_id)
);
