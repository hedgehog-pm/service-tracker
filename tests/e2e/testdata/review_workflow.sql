BEGIN;

INSERT INTO statuses (id, name, category)
OVERRIDING SYSTEM VALUE
VALUES (100, 'Review', 'in_progress');

INSERT INTO workflows (id, name, initial_status_id)
OVERRIDING SYSTEM VALUE
VALUES (100, 'Workflow with review', 1);

INSERT INTO workflow_statuses (workflow_id, status_id) VALUES
    (100, 1),
    (100, 2),
    (100, 100),
    (100, 3);

INSERT INTO workflow_transitions (workflow_id, from_status_id, to_status_id, name) VALUES
    (100, 1, 2, 'Start progress'),
    (100, 2, 100, 'Submit for review'),
    (100, 100, 3, 'Complete');

UPDATE workflow_schemes SET is_default = FALSE WHERE is_default;

INSERT INTO workflow_schemes (id, name, is_default)
OVERRIDING SYSTEM VALUE
VALUES (100, 'Workflow scheme with review', TRUE);

INSERT INTO workflow_scheme_rules (scheme_id, task_type_id, workflow_id)
SELECT 100, id, 100 FROM task_types;

INSERT INTO board_schemes (id, workflow_scheme_id, name, is_default, is_shared)
OVERRIDING SYSTEM VALUE
VALUES (100, 100, 'Доска с проверкой', TRUE, TRUE);

INSERT INTO board_scheme_columns (id, scheme_id, name, position)
OVERRIDING SYSTEM VALUE
VALUES
    (100, 100, 'К работе', 0),
    (101, 100, 'В работе', 1),
    (102, 100, 'На проверке', 2),
    (103, 100, 'Готово', 3);

INSERT INTO board_scheme_column_statuses (scheme_id, column_id, status_id) VALUES
    (100, 100, 1),
    (100, 101, 2),
    (100, 102, 100),
    (100, 103, 3);

COMMIT;
