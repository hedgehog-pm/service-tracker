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

func pointerToNull[T any](value *T) sql.Null[T] {
	if value == nil {
		return sql.Null[T]{}
	}

	return sql.Null[T]{
		V:     *value,
		Valid: true,
	}
}

func taskToDomain(row bobgen.GetTaskRow) domain.Task {
	return domain.Task{
		ID:          row.ID,
		ProjectID:   row.ProjectID,
		ProjectKey:  row.ProjectKey,
		Number:      row.Number,
		TaskTypeID:  row.TaskTypeID,
		WorkflowID:  row.WorkflowID,
		StatusID:    row.StatusID,
		Title:       row.Title,
		Description: row.Description,
		PriorityID:  row.PriorityID,
		ModuleID:    nullToPointer(row.ModuleID),
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func nullToPointer[T any](value sql.Null[T]) *T {
	if !value.Valid {
		return nil
	}

	return &value.V
}
