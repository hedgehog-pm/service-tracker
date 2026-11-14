-- +goose Up
CREATE TABLE board_columns (
    id       INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    board_id INT NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name     TEXT NOT NULL,
    position INT NOT NULL CHECK (position >= 0),
    UNIQUE (id, board_id),
    -- При перестановке колонок проверяем итоговый порядок в конце транзакции.
    UNIQUE (board_id, position) DEFERRABLE INITIALLY DEFERRED
);
