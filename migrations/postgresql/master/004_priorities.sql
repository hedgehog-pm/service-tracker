-- +goose Up
-- Укаждого проекта может быть свой набор приоритетов.
CREATE TABLE priorities (
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name     TEXT NOT NULL UNIQUE,
    position INT NOT NULL CHECK (position >= 0)
);

CREATE TABLE priority_schemes (
    id                  INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name                TEXT NOT NULL UNIQUE,
    is_default          BOOLEAN NOT NULL,
    default_priority_id INT NOT NULL
);

CREATE UNIQUE INDEX priority_schemes_one_default
    ON priority_schemes (is_default) WHERE is_default;

CREATE TABLE priority_scheme_priorities (
    scheme_id   INT NOT NULL REFERENCES priority_schemes(id) ON DELETE CASCADE,
    priority_id INT NOT NULL REFERENCES priorities(id),
    position    INT NOT NULL CHECK (position >= 0),
    PRIMARY KEY (scheme_id, priority_id)
);

-- Набор приоритетов добавляем после схемы.
-- Откладываем проверку default_priority_id до конца транзакции.
ALTER TABLE priority_schemes
    ADD FOREIGN KEY (id, default_priority_id)
    REFERENCES priority_scheme_priorities(scheme_id, priority_id)
    DEFERRABLE INITIALLY DEFERRED;

-- Задаём id приоритетов вручную для начальной схемы.
-- OVERRIDING SYSTEM VALUE разрешает явные id для GENERATED ALWAYS AS IDENTITY.
INSERT INTO priorities (id, name, position)
OVERRIDING SYSTEM VALUE
VALUES
    (1, 'Lowest', 1),
    (2, 'Low', 2),
    (3, 'Medium', 3),
    (4, 'High', 4),
    (5, 'Highest', 5);

INSERT INTO priority_schemes (id, name, is_default, default_priority_id)
OVERRIDING SYSTEM VALUE
VALUES (1, 'Default priority scheme', TRUE, 3);

INSERT INTO priority_scheme_priorities (scheme_id, priority_id, position) VALUES
    (1, 1, 1),
    (1, 2, 2),
    (1, 3, 3),
    (1, 4, 4),
    (1, 5, 5);

-- При вставке явных id последовательности не меняются.
-- Сдвигаем их, чтобы новые записи не получили занятые id.
ALTER SEQUENCE priorities_id_seq RESTART WITH 6;
ALTER SEQUENCE priority_schemes_id_seq RESTART WITH 2;
