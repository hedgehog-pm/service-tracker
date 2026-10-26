-- InsertProject
INSERT INTO
  projects (
    key,
    name,
    next_number,
    workflow_scheme_id,
    priority_scheme_id
  )
VALUES
  ($1, $2, $3, $4, $5)/* rows */
RETURNING
  id,
  key,
  name,
  next_number,
  workflow_scheme_id,
  priority_scheme_id;

-- GetProject
SELECT
  id,
  key,
  name,
  next_number,
  workflow_scheme_id,
  priority_scheme_id
FROM
  projects
WHERE
  id = $1;

-- ListProjects
SELECT
  id,
  key,
  name,
  next_number,
  workflow_scheme_id,
  priority_scheme_id
FROM
  projects
ORDER BY
  id ASC
OFFSET
  $1/* offset */
LIMIT
  $2/* limit */;

-- UpdateProject
UPDATE
  projects
SET
  name = $1
WHERE
  id = $2
RETURNING
  id,
  key,
  name,
  next_number,
  workflow_scheme_id,
  priority_scheme_id;

-- DeleteProject
DELETE FROM
  projects
WHERE
  id = $1;

-- ReserveTaskNumber
UPDATE
  projects
SET
  next_number = next_number + 1
WHERE
  id = $1
RETURNING
  next_number - 1 AS number,
  key;
