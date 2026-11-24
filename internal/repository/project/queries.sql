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

-- LockProject
SELECT
  id
FROM
  projects
WHERE
  id = $1
FOR NO KEY UPDATE;

-- GetProjectSummary
WITH project_tasks AS (
  SELECT
    t.id, t.number, t.title, t.status_id, t.priority_id,
    t.due_date, t.estimate_seconds, s.category
  FROM tasks AS t
  JOIN statuses AS s ON s.id = t.status_id
  WHERE t.project_id = $1/* project_id */
),
task_totals AS (
  SELECT
    COUNT(*) AS tasks_total,
    COUNT(*) FILTER (WHERE category <> 'done') AS tasks_open,
    COUNT(*) FILTER (
      WHERE category <> 'done' AND due_date >= $2/* today */::DATE AND due_date < $2::DATE + 7
    ) AS tasks_due_soon,
    COUNT(*) FILTER (WHERE category <> 'done' AND due_date < $2::DATE) AS tasks_overdue,
    SUM(estimate_seconds) AS total_estimate_seconds
  FROM project_tasks
),
time_total AS (
  SELECT COALESCE(SUM(w.time_spent_seconds), 0)::BIGINT AS total_time_spent_seconds
  FROM task_worklogs AS w
  JOIN project_tasks AS t ON t.id = w.task_id
),
deadlines AS (
  SELECT id, number, title, status_id, priority_id, due_date
  FROM project_tasks
  WHERE category <> 'done' AND due_date >= $2::DATE
  ORDER BY due_date, number
  LIMIT 5
)
SELECT
  totals.tasks_total,
  totals.tasks_open,
  totals.tasks_due_soon,
  totals.tasks_overdue,
  totals.total_estimate_seconds/* :*int64:notnull */,
  time_total.total_time_spent_seconds,
  p.key AS project_key,
  d.id/* ::null */,
  d.number/* ::null */,
  d.title/* ::null */,
  d.status_id/* ::null */,
  d.priority_id/* ::null */,
  d.due_date/* ::null */
FROM projects AS p
CROSS JOIN task_totals AS totals
CROSS JOIN time_total
LEFT JOIN deadlines AS d ON TRUE
WHERE p.id = $1
ORDER BY d.due_date, d.number;
