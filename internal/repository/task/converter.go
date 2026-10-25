package task

import (
	"database/sql"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/task"
)

func insertedTaskToDomain(task domain.Task, row bobgen.InsertTaskRow) domain.Task {
	task.ID = row.ID
	task.CreatedAt = row.CreatedAt
	task.UpdatedAt = row.UpdatedAt

	return task
}

func updatedTaskToDomain(task domain.Task, row bobgen.UpdateTaskRow) domain.Task {
	task.CreatedAt = row.CreatedAt
	task.UpdatedAt = row.UpdatedAt

	return task
}

func taskToInsertParams(task domain.Task) []bobgen.InsertTask_Rows {
	return []bobgen.InsertTask_Rows{{
		ProjectID: task.ProjectID, Number: task.Number, Title: task.Title, Description: task.Description,
		StatusID: task.StatusID, ModuleID: pointerToNull(task.ModuleID), CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
	}}
}

func pointerToNull[T any](value *T) sql.Null[T] {
	if value == nil {
		return sql.Null[T]{}
	}

	return sql.Null[T]{V: *value, Valid: true}
}

func taskToDomain(row bobgen.GetTaskRow) domain.Task {
	return domain.Task{
		ID: row.ID, ProjectID: row.ProjectID, ProjectKey: row.ProjectKey, Number: row.Number,
		Title: row.Title, Description: row.Description, StatusID: row.StatusID, ModuleID: nullToPointer(row.ModuleID),
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}

func nullToPointer[T any](value sql.Null[T]) *T {
	if !value.Valid {
		return nil
	}

	return &value.V
}
