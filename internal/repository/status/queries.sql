-- ListStatuses
SELECT DISTINCT
  status.id,
  status.name,
  status.category
FROM
  workflow_scheme_rules AS rule
  JOIN workflow_statuses AS membership ON membership.workflow_id = rule.workflow_id
  JOIN statuses AS status ON status.id = membership.status_id
WHERE
  rule.scheme_id = $1
ORDER BY
  status.id ASC;

-- ListProjectStatuses
SELECT DISTINCT
  status.id,
  status.name,
  status.category
FROM
  projects AS project
  JOIN workflow_scheme_rules AS rule ON rule.scheme_id = project.workflow_scheme_id
  JOIN workflow_statuses AS membership ON membership.workflow_id = rule.workflow_id
  JOIN statuses AS status ON status.id = membership.status_id
WHERE
  project.id = $1
ORDER BY
  status.id ASC;
