-- InsertModule
INSERT INTO
  modules (
    project_id,
    name
  )
SELECT
  id,
  $1/* name */
FROM
  projects
WHERE
  id = $2/* project_id */
RETURNING
  id,
  project_id,
  name;

-- GetModule
SELECT
  id,
  project_id,
  name
FROM
  modules
WHERE
  project_id = $1
  AND id = $2;

-- ListModules
WITH module_page AS (
  SELECT id, project_id, name
  FROM modules
  WHERE project_id = $1
  ORDER BY id
  OFFSET $2/* offset */
  LIMIT $3/* limit */
)
SELECT
  m.id,
  m.project_id,
  m.name,
  COUNT(t.id) AS tasks_total,
  COUNT(t.id) FILTER (WHERE s.category = 'done') AS tasks_done
FROM module_page AS m
LEFT JOIN tasks AS t ON t.project_id = m.project_id AND t.module_id = m.id
LEFT JOIN statuses AS s ON s.id = t.status_id
GROUP BY m.id, m.project_id, m.name
ORDER BY m.id;

-- UpdateModule
UPDATE
  modules
SET
  name = $1
WHERE
  project_id = $2
  AND id = $3
RETURNING
  id,
  project_id,
  name;

-- DeleteModule
DELETE FROM
  modules
WHERE
  project_id = $1
  AND id = $2;
