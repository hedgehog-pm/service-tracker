-- +goose Up
-- Один статус можно использовать в разных workflow.
CREATE TABLE statuses (
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     TEXT NOT NULL UNIQUE,
    category TEXT NOT NULL CHECK (category IN ('todo', 'in_progress', 'done'))
);

CREATE TABLE workflows (
    id                INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name              TEXT NOT NULL UNIQUE,
    initial_status_id INT NOT NULL
);

CREATE TABLE workflow_statuses (
    workflow_id INT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    status_id   INT NOT NULL REFERENCES statuses(id),
    PRIMARY KEY (workflow_id, status_id)
);

-- Статусы добавляем после workflow, поэтому начальный статус проверяем в конце транзакции.
ALTER TABLE workflows
    ADD FOREIGN KEY (id, initial_status_id)
    REFERENCES workflow_statuses(workflow_id, status_id)
    DEFERRABLE INITIALLY DEFERRED;

-- Переход связывает разные статусы одного workflow.
CREATE TABLE workflow_transitions (
    id             INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    workflow_id    INT NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    from_status_id INT NOT NULL,
    to_status_id   INT NOT NULL,
    name           TEXT NOT NULL,
    CHECK (from_status_id <> to_status_id),
    FOREIGN KEY (workflow_id, from_status_id) REFERENCES workflow_statuses(workflow_id, status_id),
    FOREIGN KEY (workflow_id, to_status_id) REFERENCES workflow_statuses(workflow_id, status_id)
);

CREATE TABLE workflow_schemes (
    id         INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    is_default BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX workflow_schemes_one_default
    ON workflow_schemes (is_default) WHERE is_default;

-- В схеме каждый тип задачи связан с одним workflow.
CREATE TABLE workflow_scheme_rules (
    scheme_id     INT NOT NULL REFERENCES workflow_schemes(id) ON DELETE CASCADE,
    issue_type_id INT NOT NULL REFERENCES issue_types(id),
    workflow_id   INT NOT NULL REFERENCES workflows(id),
    PRIMARY KEY (scheme_id, issue_type_id),
    UNIQUE (scheme_id, issue_type_id, workflow_id)
);

-- Задаём id вручную, чтобы явно перечислить связи начальных записей.
-- OVERRIDING SYSTEM VALUE разрешает явные id для GENERATED ALWAYS AS IDENTITY.
INSERT INTO statuses (id, name, category)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 'Open', 'todo'),
    (2, 'In Progress', 'in_progress'),
    (3, 'Done', 'done');

INSERT INTO workflows (id, name, initial_status_id)
OVERRIDING SYSTEM VALUE
VALUES (1, 'Default workflow', 1);

INSERT INTO workflow_statuses (workflow_id, status_id) VALUES
    (1, 1),
    (1, 2),
    (1, 3);

INSERT INTO workflow_transitions (workflow_id, from_status_id, to_status_id, name) VALUES
    (1, 1, 2, 'Start progress'),
    (1, 2, 3, 'Complete'),
    (1, 2, 1, 'Return to open'),
    (1, 3, 1, 'Reopen');

INSERT INTO workflow_schemes (id, name, is_default)
OVERRIDING SYSTEM VALUE
VALUES (1, 'Default workflow scheme', TRUE);

INSERT INTO workflow_scheme_rules (scheme_id, issue_type_id, workflow_id) VALUES
    (1, 1, 1),
    (1, 2, 1),
    (1, 3, 1);

-- При вставке явных id последовательности не меняются.
-- Сдвигаем их, чтобы новые записи не получили занятые id.
ALTER SEQUENCE statuses_id_seq RESTART WITH 4;
ALTER SEQUENCE workflows_id_seq RESTART WITH 2;
ALTER SEQUENCE workflow_schemes_id_seq RESTART WITH 2;
