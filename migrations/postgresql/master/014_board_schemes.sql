-- +goose Up
CREATE TABLE board_schemes (
    id                 INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    workflow_scheme_id INT NOT NULL REFERENCES workflow_schemes(id),
    name               TEXT NOT NULL,
    is_default         BOOLEAN NOT NULL DEFAULT FALSE,
    is_shared          BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE (id, workflow_scheme_id),
    CHECK (NOT is_default OR is_shared)
);

CREATE UNIQUE INDEX board_schemes_one_default_per_workflow_scheme
    ON board_schemes (workflow_scheme_id) WHERE is_default;

ALTER TABLE projects ADD UNIQUE (id, workflow_scheme_id);

ALTER TABLE boards
    ADD COLUMN board_scheme_id INT,
    ADD COLUMN workflow_scheme_id INT;

-- Каждая существующая доска получает свою схему, чтобы сохранить настройки колонок.
INSERT INTO board_schemes (id, workflow_scheme_id, name)
OVERRIDING SYSTEM VALUE
SELECT b.id, p.workflow_scheme_id, b.name
FROM boards AS b
JOIN projects AS p ON p.id = b.project_id;

SELECT setval(
    'board_schemes_id_seq',
    COALESCE(MAX(id), 1),
    COUNT(*) > 0
)
FROM board_schemes;

UPDATE boards AS b
SET board_scheme_id = s.id, workflow_scheme_id = s.workflow_scheme_id
FROM board_schemes AS s
WHERE s.id = b.id;

-- Обе ссылки должны вести к одной workflow-схеме.
ALTER TABLE boards
    DROP CONSTRAINT boards_project_id_fkey,
    ALTER COLUMN board_scheme_id SET NOT NULL,
    ALTER COLUMN workflow_scheme_id SET NOT NULL,
    ADD FOREIGN KEY (project_id, workflow_scheme_id)
        REFERENCES projects(id, workflow_scheme_id) ON DELETE CASCADE,
    ADD FOREIGN KEY (board_scheme_id, workflow_scheme_id)
        REFERENCES board_schemes(id, workflow_scheme_id);

-- Меняем владельца колонок на схему, сохраняя идентификаторы и связи статусов.
ALTER TABLE board_columns DROP CONSTRAINT board_columns_board_id_fkey;
ALTER TABLE board_columns RENAME TO board_scheme_columns;
ALTER TABLE board_scheme_columns RENAME COLUMN board_id TO scheme_id;
ALTER TABLE board_scheme_columns ADD FOREIGN KEY (scheme_id)
    REFERENCES board_schemes(id) ON DELETE CASCADE;
ALTER SEQUENCE board_columns_id_seq RENAME TO board_scheme_columns_id_seq;

ALTER TABLE board_column_statuses RENAME TO board_scheme_column_statuses;
ALTER TABLE board_scheme_column_statuses RENAME COLUMN board_id TO scheme_id;
