-- GetTaskWorklog
SELECT
  w.id,
  w.task_id,
  w.started_at,
  w.time_spent_seconds,
  w.comment,
  w.created_at,
  w.updated_at
FROM
  task_worklogs AS w
  JOIN tasks AS t ON t.id = w.task_id
WHERE
  t.project_id = $1
  AND t.id = $2
  AND w.id = $3;

-- ListTaskWorklogs
WITH task AS (
  SELECT id
  FROM tasks
  WHERE project_id = $1/* project_id */ AND id = $2/* task_id */
)
SELECT
  w.id/* ::null */,
  w.task_id/* ::null */,
  w.started_at/* ::null */,
  w.time_spent_seconds/* ::null */,
  w.comment/* ::null */,
  w.created_at/* ::null */,
  w.updated_at/* ::null */,
  (
    SELECT COALESCE(SUM(all_worklogs.time_spent_seconds), 0)::BIGINT
    FROM task_worklogs AS all_worklogs
    WHERE all_worklogs.task_id = t.id
  ) AS total_time_spent_seconds
FROM
  task AS t
  LEFT JOIN (
    SELECT
      id,
      task_id,
      started_at,
      time_spent_seconds,
      comment,
      created_at,
      updated_at
    FROM task_worklogs
    WHERE task_id = $2
    ORDER BY started_at DESC, id DESC
    OFFSET $3/* offset */
    LIMIT $4/* limit */
  ) AS w ON w.task_id = t.id
ORDER BY
  w.started_at DESC,
  w.id DESC;

-- InsertTaskWorklog
INSERT INTO task_worklogs (
  task_id,
  started_at,
  time_spent_seconds,
  comment,
  created_at,
  updated_at
)
SELECT
  t.id,
  $1/* started_at */,
  $2/* time_spent_seconds */,
  $3/* comment::null */,
  $4/* created_at */,
  $5/* updated_at */
FROM tasks AS t
WHERE
  t.project_id = $6/* project_id */
  AND t.id = $7/* task_id */
RETURNING id, task_id, started_at, time_spent_seconds, comment, created_at, updated_at;

-- UpdateTaskWorklog
UPDATE task_worklogs AS w
SET
  started_at = $1,
  time_spent_seconds = $2,
  comment = $3/* comment::null */,
  updated_at = $4
WHERE
  w.task_id = $6/* task_id */
  AND w.id = $7
  AND EXISTS (
    SELECT 1 FROM tasks AS t
    WHERE t.id = w.task_id AND t.project_id = $5/* project_id */
  )
RETURNING w.id, w.task_id, w.started_at, w.time_spent_seconds, w.comment, w.created_at, w.updated_at;

-- DeleteTaskWorklog
DELETE FROM task_worklogs AS w
WHERE
  w.task_id = $2/* task_id */
  AND w.id = $3
  AND EXISTS (
    SELECT 1 FROM tasks AS t
    WHERE t.id = w.task_id AND t.project_id = $1/* project_id */
  );

-- ListProjectWorklogs
WITH filtered AS (
  SELECT
    w.id, w.task_id, w.started_at, w.time_spent_seconds, w.comment,
    p.key || '-' || t.number::TEXT AS task_key,
    t.title AS task_title
  FROM task_worklogs AS w
  JOIN tasks AS t ON t.id = w.task_id
  JOIN projects AS p ON p.id = t.project_id
  WHERE t.project_id = $1/* project_id */
    AND ($2/* from_time:*time.Time:notnull */::TIMESTAMPTZ IS NULL OR w.started_at >= $2)
    AND ($3/* to_time:*time.Time:notnull */::TIMESTAMPTZ IS NULL OR w.started_at < $3)
    AND (
      $4/* search */::TEXT = ''
      OR strpos(lower(t.title), lower($4)) > 0
      OR strpos(lower(p.key || '-' || t.number::TEXT), lower($4)) > 0
    )
),
time_total AS (
  SELECT COALESCE(SUM(time_spent_seconds), 0)::BIGINT AS total_time_spent_seconds
  FROM filtered
),
page AS (
  SELECT id, task_id, task_key, task_title, started_at, time_spent_seconds, comment
  FROM filtered
  ORDER BY started_at DESC, id DESC
  LIMIT $5/* limit */
  OFFSET $6/* offset */
)
SELECT
  time_total.total_time_spent_seconds,
  (SELECT SUM(estimate_seconds) FROM tasks WHERE project_id = p.id)
    AS total_estimate_seconds/* :*int64:notnull */,
  page.id/* ::null */,
  page.task_id/* ::null */,
  page.task_key/* ::null */,
  page.task_title/* ::null */,
  page.started_at/* ::null */,
  page.time_spent_seconds/* ::null */,
  page.comment/* ::null */
FROM projects AS p
CROSS JOIN time_total
LEFT JOIN page ON TRUE
WHERE p.id = $1
ORDER BY page.started_at DESC, page.id DESC;

-- ListProjectWorklogsByTask
WITH filtered AS (
  SELECT
    w.task_id, w.time_spent_seconds,
    p.key || '-' || t.number::TEXT AS task_key,
    t.title AS task_title
  FROM task_worklogs AS w
  JOIN tasks AS t ON t.id = w.task_id
  JOIN projects AS p ON p.id = t.project_id
  WHERE t.project_id = $1/* project_id */
    AND ($2/* from_time:*time.Time:notnull */::TIMESTAMPTZ IS NULL OR w.started_at >= $2)
    AND ($3/* to_time:*time.Time:notnull */::TIMESTAMPTZ IS NULL OR w.started_at < $3)
    AND (
      $4/* search */::TEXT = ''
      OR strpos(lower(t.title), lower($4)) > 0
      OR strpos(lower(p.key || '-' || t.number::TEXT), lower($4)) > 0
    )
),
task_totals AS (
  SELECT task_id, task_key, task_title,
    COUNT(*) AS worklogs_count, SUM(time_spent_seconds) AS time_spent_seconds
  FROM filtered
  GROUP BY task_id, task_key, task_title
),
time_total AS (
  SELECT COALESCE(SUM(time_spent_seconds), 0)::BIGINT AS total_time_spent_seconds
  FROM task_totals
),
page AS (
  SELECT task_id, task_key, task_title, worklogs_count, time_spent_seconds
  FROM task_totals
  ORDER BY time_spent_seconds DESC, task_id
  LIMIT $5/* limit */
  OFFSET $6/* offset */
)
SELECT
  time_total.total_time_spent_seconds,
  (SELECT SUM(estimate_seconds) FROM tasks WHERE project_id = p.id)
    AS total_estimate_seconds/* :*int64:notnull */,
  page.task_id/* ::null */,
  page.task_key/* ::null */,
  page.task_title/* ::null */,
  page.worklogs_count/* ::null */,
  page.time_spent_seconds/* ::null */
FROM projects AS p
CROSS JOIN time_total
LEFT JOIN page ON TRUE
WHERE p.id = $1
ORDER BY page.time_spent_seconds DESC, page.task_id;
