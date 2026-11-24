-- GetProjectKey
SELECT
  key
FROM
  projects
WHERE
  id = $1;

-- GetTask
SELECT
  t.id,
  t.project_id,
  p.key AS project_key,
  t.number,
  t.rank,
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
  t.estimate_seconds,
  t.due_date,
  t.created_at,
  t.updated_at
FROM
  tasks AS t
  JOIN projects AS p ON p.id = t.project_id
WHERE
  t.project_id = $1
  AND t.id = $2;

-- GetTaskForUpdate
SELECT
  t.id,
  t.project_id,
  p.key AS project_key,
  t.number,
  t.rank,
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
  t.estimate_seconds,
  t.due_date,
  t.created_at,
  t.updated_at
FROM
  tasks AS t
  JOIN projects AS p ON p.id = t.project_id
WHERE
  t.project_id = $1
  AND t.id = $2
FOR UPDATE OF t;

-- ListTasks
SELECT
  t.id,
  t.project_id,
  p.key AS project_key,
  t.number,
  t.rank,
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
  t.estimate_seconds,
  t.due_date,
  t.created_at,
  t.updated_at
FROM
  tasks AS t
  JOIN projects AS p ON p.id = t.project_id
WHERE
  t.project_id = $1
ORDER BY
  t.number ASC;

-- InsertTask
INSERT INTO
  tasks (
    project_id,
    workflow_scheme_id,
    priority_scheme_id,
    number,
    rank,
    title,
    description,
    task_type_id,
    workflow_id,
    status_id,
    priority_id,
    module_id,
    estimate_seconds,
    due_date,
    created_at,
    updated_at
  )
SELECT
  p.id,
  p.workflow_scheme_id,
  p.priority_scheme_id,
  $1/* number */,
  $2/* rank */,
  $3/* title */,
  $4/* description */,
  $5/* task_type_id */,
  $6/* workflow_id */,
  $7/* status_id */,
  $8/* priority_id */,
  $9/* module_id::null */,
  $10/* estimate_seconds::null */,
  $11/* due_date::null */,
  $12/* created_at */,
  $13/* updated_at */
FROM
  projects AS p
WHERE
  p.id = $14/* project_id */
RETURNING
  id,
  created_at,
  updated_at;

-- UpdateTask
UPDATE
  tasks
SET
  title = $1,
  description = $2,
  priority_id = $3,
  module_id = $4/* module_id::null */,
  estimate_seconds = $5/* estimate_seconds::null */,
  due_date = $6/* due_date::null */,
  updated_at = $7
WHERE
  project_id = $8
  AND id = $9
RETURNING
  created_at,
  updated_at;

-- UpdateStatus
UPDATE
  tasks
SET
  status_id = $1,
  updated_at = $2
WHERE
  project_id = $3
  AND id = $4
RETURNING
  created_at,
  updated_at;

-- DeleteTask
DELETE FROM
  tasks
WHERE
  project_id = $1
  AND id = $2;

-- ListTaskOrder
SELECT
  t.id,
  t.project_id,
  p.key AS project_key,
  t.number,
  t.rank,
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
  t.estimate_seconds,
  t.due_date,
  t.created_at,
  t.updated_at
FROM
  tasks AS t
  JOIN projects AS p ON p.id = t.project_id
WHERE
  t.project_id = $1
ORDER BY
  t.rank,
  t.number;


-- GetLastRank
SELECT
  rank
FROM
  tasks
WHERE
  project_id = $1
ORDER BY
  rank DESC
LIMIT 1;

-- GetPreviousRank
SELECT
  rank
FROM
  tasks
WHERE
  project_id = $1
  AND id <> $2
  AND rank < $3
ORDER BY
  rank DESC
LIMIT 1;

-- GetNextRank
SELECT
  rank
FROM
  tasks
WHERE
  project_id = $1
  AND id <> $2
  AND rank > $3
ORDER BY
  rank
LIMIT 1;

-- GetLastColumnRank
SELECT
  t.rank
FROM
  tasks AS t
  JOIN boards AS b ON b.project_id = t.project_id
  JOIN board_scheme_column_statuses AS membership ON membership.scheme_id = b.board_scheme_id
    AND membership.status_id = t.status_id
WHERE
  t.project_id = $1
  AND b.id = $2
  AND membership.column_id = $3
  AND t.id <> $4
ORDER BY
  t.rank DESC
LIMIT 1;

-- UpdateRank
UPDATE
  tasks
SET
  rank = $1
WHERE
  project_id = $2
  AND id = $3;

-- UpdatePosition
UPDATE
  tasks
SET
  rank = $1,
  status_id = $2,
  updated_at = $3
WHERE
  project_id = $4
  AND id = $5
RETURNING
  created_at,
  updated_at;

-- ListBoardColumnTasks
SELECT
  t.id,
  t.project_id,
  p.key AS project_key,
  t.number,
  t.rank,
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
  t.estimate_seconds,
  t.due_date,
  t.created_at,
  t.updated_at
FROM
  tasks AS t
  JOIN projects AS p ON p.id = t.project_id
WHERE
  t.project_id = $1
  AND EXISTS (
    SELECT 1
    FROM boards AS b
    JOIN board_scheme_column_statuses AS membership ON membership.scheme_id = b.board_scheme_id
    WHERE b.id = $2
      AND b.project_id = t.project_id
      AND membership.column_id = $3
      AND membership.status_id = t.status_id
  )
  AND ($4::BIGINT IS NULL OR (t.rank, t.number) > ($4/* after_rank::null */, $5/* after_number::null */))
ORDER BY
  t.rank,
  t.number
LIMIT $6;
