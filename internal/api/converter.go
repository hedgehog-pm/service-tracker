package api

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
)

func taskToAPI(task domain.Task) generated.Task {
	return generated.Task{
		ID: task.ID, ProjectID: task.ProjectID, Number: task.Number, Key: task.Key(),
		Title: task.Title, Description: task.Description, StatusID: task.StatusID, ModuleID: task.ModuleID,
		CreatedAt: task.CreatedAt, UpdatedAt: task.UpdatedAt,
	}
}
