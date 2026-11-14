-- +goose Up
ALTER TABLE tasks ADD COLUMN rank BIGINT;

-- Сохраняем порядок существующих задач и оставляем промежутки для перестановок.
UPDATE tasks SET rank = number::BIGINT * 1024;

ALTER TABLE tasks
    ALTER COLUMN rank SET NOT NULL,
    -- Ранг общий для всех досок проекта. Перенумерация проверяется после всех изменений.
    ADD CONSTRAINT tasks_project_rank_key
        UNIQUE (project_id, rank) DEFERRABLE INITIALLY DEFERRED;
