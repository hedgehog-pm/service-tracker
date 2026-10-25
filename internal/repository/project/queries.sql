-- InsertProject
INSERT INTO projects (key, name, next_number)
VALUES ($1, $2, $3) /* rows */
RETURNING id, key, name, next_number;

-- GetProject
SELECT id, key, name, next_number
FROM projects
WHERE id = $1;

-- ListProjects
SELECT id, key, name, next_number
FROM projects
ORDER BY id ASC
OFFSET $1 /* offset */ LIMIT $2 /* limit */;

-- UpdateProject
UPDATE projects
SET name = $1
WHERE id = $2
RETURNING id, key, name, next_number;

-- DeleteProject
DELETE FROM projects WHERE id = $1;

-- ReserveTaskNumber
UPDATE projects
SET next_number = next_number + 1
WHERE id = $1
RETURNING next_number - 1 AS number, key;
