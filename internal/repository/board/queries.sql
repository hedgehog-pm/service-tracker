-- GetBoard
SELECT
  b.id,
  b.project_id,
  b.name,
  c.id AS column_id,
  c.name AS column_name,
  COALESCE(c.position, 0) AS position,
  membership.status_id
FROM
  boards AS b
  LEFT JOIN board_columns AS c ON c.board_id = b.id
  LEFT JOIN board_column_statuses AS membership ON membership.column_id = c.id AND membership.board_id = b.id
WHERE
  b.project_id = $1
  AND b.id = $2
ORDER BY
  c.position,
  membership.status_id;

-- ListBoards
WITH page AS (
  SELECT
    id,
    project_id,
    name
  FROM
    boards
  WHERE
    project_id = $1
  ORDER BY
    id
  LIMIT $2
  OFFSET $3
)
SELECT
  b.id,
  b.project_id,
  b.name,
  c.id AS column_id,
  c.name AS column_name,
  COALESCE(c.position, 0) AS position,
  membership.status_id
FROM
  page AS b
  LEFT JOIN board_columns AS c ON c.board_id = b.id
  LEFT JOIN board_column_statuses AS membership ON membership.column_id = c.id AND membership.board_id = b.id
ORDER BY
  b.id,
  c.position,
  membership.status_id;

-- InsertBoard
INSERT INTO boards (project_id, name)
SELECT
  id,
  $1/* name */
FROM
  projects
WHERE
  id = $2/* project_id */
RETURNING id, project_id, name;

-- UpdateBoard
WITH updated AS (
  UPDATE boards
  SET name = $1
  WHERE project_id = $2
    AND id = $3
  RETURNING id, project_id, name
)
SELECT
  b.id,
  b.project_id,
  b.name,
  c.id AS column_id,
  c.name AS column_name,
  COALESCE(c.position, 0) AS position,
  membership.status_id
FROM
  updated AS b
  LEFT JOIN board_columns AS c ON c.board_id = b.id
  LEFT JOIN board_column_statuses AS membership ON membership.column_id = c.id AND membership.board_id = b.id
ORDER BY
  c.position,
  membership.status_id;

-- LockBoard
SELECT
  id
FROM
  boards
WHERE
  project_id = $1
  AND id = $2
FOR NO KEY UPDATE;

-- DeleteBoard
DELETE FROM
  boards
WHERE
  project_id = $1
  AND id = $2;

-- DeleteBoardStatuses
DELETE FROM board_column_statuses WHERE board_id = $1;

-- DeleteRemovedColumns
DELETE FROM
  board_columns
WHERE
  board_id = $1
  AND NOT (id = ANY($2::INT[]));

-- InsertColumn
INSERT INTO board_columns (board_id, name, position)
VALUES ($1, $2, $3)/* rows */
RETURNING id;

-- UpdateColumn
UPDATE
  board_columns
SET
  name = $1,
  position = $2
WHERE
  board_id = $3
  AND id = $4;

-- InsertColumnStatus
INSERT INTO board_column_statuses (board_id, column_id, status_id) VALUES ($1, $2, $3)/* rows */;
