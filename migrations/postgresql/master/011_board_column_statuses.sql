-- +goose Up
CREATE TABLE board_column_statuses (
    board_id  INT NOT NULL,
    column_id INT NOT NULL,
    status_id INT NOT NULL REFERENCES statuses(id),
    -- В пределах доски статус отображается только в одной колонке.
    PRIMARY KEY (board_id, status_id),
    -- Колонка должна принадлежать той же доске.
    FOREIGN KEY (column_id, board_id)
        REFERENCES board_columns(id, board_id) ON DELETE CASCADE
);
