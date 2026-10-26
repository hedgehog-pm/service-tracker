package api

import (
	"time"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
)

func taskToAPI(task domain.Task) generated.Task {
	return generated.Task{
		ID:              task.ID,
		ProjectID:       task.ProjectID,
		Number:          task.Number,
		Key:             task.Key(),
		Title:           task.Title,
		Description:     task.Description,
		StatusID:        task.StatusID,
		TaskTypeID:      task.TaskTypeID,
		PriorityID:      task.PriorityID,
		ModuleID:        task.ModuleID,
		ParentID:        task.ParentID,
		EstimateSeconds: task.EstimateSeconds,
		DueDate:         dateToAPI(task.DueDate),
		CreatedAt:       task.CreatedAt,
		UpdatedAt:       task.UpdatedAt,
	}
}

func dateToDomain(date *openapi_types.Date) *time.Time {
	if date == nil {
		return nil
	}

	return &date.Time
}

func dateToAPI(date *time.Time) *openapi_types.Date {
	if date == nil {
		return nil
	}

	return &openapi_types.Date{Time: *date}
}
