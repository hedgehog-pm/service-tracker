-- InsertStatus
INSERT INTO statuses (project_id, name, category)
SELECT id, $1 /* name */, $2 /* category */
FROM projects
WHERE id = $3 /* project_id */
RETURNING id, project_id, name, category;

-- GetStatus
SELECT id, project_id, name, category
FROM statuses
WHERE project_id = $1 AND id = $2;

-- ListStatuses
SELECT id, project_id, name, category
FROM statuses
WHERE project_id = $1
ORDER BY id ASC
OFFSET $2 /* offset */ LIMIT $3 /* limit */;

-- UpdateStatus
UPDATE statuses
SET name = $1, category = $2
WHERE project_id = $3 AND id = $4
RETURNING id, project_id, name, category;

-- DeleteStatus
DELETE FROM statuses WHERE project_id = $1 AND id = $2;
