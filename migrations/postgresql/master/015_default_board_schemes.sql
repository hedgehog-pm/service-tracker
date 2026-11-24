-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    new_scheme_id INT;
    open_column_id INT;
    in_progress_column_id INT;
    done_column_id INT;
BEGIN
    -- ID штатной workflow-схемы и её статусов заданы в 003_workflows.sql.
    -- ID схемы доски и колонок получаем из последовательностей, поскольку в БД уже могут быть доски.
    INSERT INTO board_schemes (workflow_scheme_id, name, is_default, is_shared)
    VALUES (1, 'Основная', TRUE, TRUE)
    RETURNING id INTO new_scheme_id;

    INSERT INTO board_scheme_columns (scheme_id, name, position)
    VALUES (new_scheme_id, 'Открытые', 0)
    RETURNING id INTO open_column_id;

    INSERT INTO board_scheme_columns (scheme_id, name, position)
    VALUES (new_scheme_id, 'В работе', 1)
    RETURNING id INTO in_progress_column_id;

    INSERT INTO board_scheme_columns (scheme_id, name, position)
    VALUES (new_scheme_id, 'Готово', 2)
    RETURNING id INTO done_column_id;

    INSERT INTO board_scheme_column_statuses (scheme_id, column_id, status_id) VALUES
        (new_scheme_id, open_column_id, 1),
        (new_scheme_id, in_progress_column_id, 2),
        (new_scheme_id, done_column_id, 3);
END;
$$;
-- +goose StatementEnd
