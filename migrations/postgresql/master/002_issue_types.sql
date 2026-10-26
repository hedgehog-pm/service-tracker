-- +goose Up
-- Один тип задачи можно использовать вразных проектах.
CREATE TABLE issue_types (
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     TEXT NOT NULL UNIQUE,
    position INT NOT NULL CHECK (position >= 0)
);

-- Начальные правила workflow ссылаются на эти id.
-- OVERRIDING SYSTEM VALUE разрешает явные id для GENERATED ALWAYS AS IDENTITY.
INSERT INTO issue_types (id, name, position)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 'Task', 1),
    (2, 'Bug', 2),
    (3, 'Story', 3),
    (4, 'Epic', 4);

-- При вставке явных id последовательности не меняются.
-- Сдвигаем их, чтобы новые записи не получили занятые id.
ALTER SEQUENCE issue_types_id_seq RESTART WITH 5;
