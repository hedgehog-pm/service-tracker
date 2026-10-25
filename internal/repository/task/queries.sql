-- GetProjectKey
SELECT key
FROM projects
WHERE id = $1;

-- GetTask
SELECT t.id, t.project_id, p.key AS project_key, t.number,
       t.title, t.description, t.status_id, t.module_id,
       t.created_at, t.updated_at
FROM tasks AS t
JOIN projects AS p ON p.id = t.project_id
WHERE t.project_id = $1 AND t.id = $2;

-- ListTasks
SELECT t.id, t.project_id, p.key AS project_key, t.number,
       t.title, t.description, t.status_id, t.module_id,
       t.created_at, t.updated_at
FROM tasks AS t
JOIN projects AS p ON p.id = t.project_id
WHERE t.project_id = $1
ORDER BY t.number ASC;

-- InsertTask
INSERT INTO tasks (project_id, number, title, description, status_id, module_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) /* rows */
RETURNING id, created_at, updated_at;

-- UpdateTask
UPDATE tasks
SET title = $1, description = $2, status_id = $3,
    module_id = $4 /* module_id::null */, updated_at = $5
WHERE project_id = $6 AND id = $7
RETURNING created_at, updated_at;

-- DeleteTask
DELETE FROM tasks
WHERE project_id = $1 AND id = $2;
