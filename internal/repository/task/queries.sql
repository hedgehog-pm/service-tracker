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
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
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
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
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
  t.title,
  t.description,
  t.task_type_id,
  t.workflow_id,
  t.status_id,
  t.priority_id,
  t.module_id,
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
    title,
    description,
    task_type_id,
    workflow_id,
    status_id,
    priority_id,
    module_id,
    created_at,
    updated_at
  )
SELECT
  p.id,
  p.workflow_scheme_id,
  p.priority_scheme_id,
  $1/* number */,
  $2/* title */,
  $3/* description */,
  $4/* task_type_id */,
  $5/* workflow_id */,
  $6/* status_id */,
  $7/* priority_id */,
  $8/* module_id::null */,
  $9/* created_at */,
  $10/* updated_at */
FROM
  projects AS p
WHERE
  p.id = $11/* project_id */
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
  updated_at = $5
WHERE
  project_id = $6
  AND id = $7
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
