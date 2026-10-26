-- +goose Up
ALTER TABLE issue_types RENAME TO task_types;
ALTER SEQUENCE issue_types_id_seq RENAME TO task_types_id_seq;

ALTER TABLE task_types RENAME CONSTRAINT issue_types_pkey TO task_types_pkey;
ALTER TABLE task_types RENAME CONSTRAINT issue_types_name_key TO task_types_name_key;
ALTER TABLE task_types RENAME CONSTRAINT issue_types_position_check TO task_types_position_check;

ALTER TABLE workflow_scheme_rules RENAME COLUMN issue_type_id TO task_type_id;
ALTER TABLE workflow_scheme_rules
    RENAME CONSTRAINT workflow_scheme_rules_issue_type_id_fkey
    TO workflow_scheme_rules_task_type_id_fkey;
ALTER TABLE workflow_scheme_rules
    RENAME CONSTRAINT workflow_scheme_rules_scheme_id_issue_type_id_workflow_id_key
    TO workflow_scheme_rules_scheme_id_task_type_id_workflow_id_key;

ALTER TABLE tasks RENAME COLUMN issue_type_id TO task_type_id;
ALTER TABLE tasks
    RENAME CONSTRAINT tasks_workflow_scheme_id_issue_type_id_workflow_id_fkey
    TO tasks_workflow_scheme_id_task_type_id_workflow_id_fkey;
