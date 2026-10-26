-- +goose Up
-- Сначала назначаем схемы существующим проектам, затем запрещаем пустые значения.
ALTER TABLE projects
    ADD COLUMN workflow_scheme_id INT REFERENCES workflow_schemes(id),
    ADD COLUMN priority_scheme_id INT REFERENCES priority_schemes(id);

UPDATE projects
SET workflow_scheme_id = (SELECT id FROM workflow_schemes WHERE is_default),
    priority_scheme_id = (SELECT id FROM priority_schemes WHERE is_default);

-- Связь задачи спроектом должна учитывать обе его схемы.
ALTER TABLE projects
    ALTER COLUMN workflow_scheme_id SET NOT NULL,
    ALTER COLUMN priority_scheme_id SET NOT NULL,
    ADD UNIQUE (id, workflow_scheme_id, priority_scheme_id);
