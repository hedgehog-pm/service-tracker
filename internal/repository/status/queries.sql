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
