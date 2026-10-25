-- InsertModule
INSERT INTO modules (project_id, name)
SELECT id, $1 /* name */
FROM projects
WHERE id = $2 /* project_id */
RETURNING id, project_id, name;

-- GetModule
SELECT id, project_id, name
FROM modules
WHERE project_id = $1 AND id = $2;

-- ListModules
SELECT id, project_id, name
FROM modules
WHERE project_id = $1
ORDER BY id ASC
OFFSET $2 /* offset */ LIMIT $3 /* limit */;

-- UpdateModule
UPDATE modules
SET name = $1
WHERE project_id = $2 AND id = $3
RETURNING id, project_id, name;

-- DeleteModule
DELETE FROM modules WHERE project_id = $1 AND id = $2;
