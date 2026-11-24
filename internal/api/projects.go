package api

import (
	"context"

	openapiTypes "github.com/oapi-codegen/runtime/types"
	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createproject"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/project"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/projectsummary"
)

type ProjectHandler struct {
	projects      *project.Usecase
	createProject *createproject.Usecase
	summary       *projectsummary.Usecase
}

func NewProjectHandler(
	projects *project.Usecase,
	createProject *createproject.Usecase,
	summary *projectsummary.Usecase,
) *ProjectHandler {
	return &ProjectHandler{
		projects:      projects,
		createProject: createProject,
		summary:       summary,
	}
}

func (h *ProjectHandler) CreateProject(
	ctx context.Context,
	request generated.CreateProjectRequestObject,
) (generated.CreateProjectResponseObject, error) {
	body := request.Body

	value, err := h.createProject.Do(ctx, body.Key, body.Name)
	if err != nil {
		return nil, err
	}

	return generated.CreateProject201JSONResponse(projectToAPI(value)), nil
}

func (h *ProjectHandler) GetProject(
	ctx context.Context,
	request generated.GetProjectRequestObject,
) (generated.GetProjectResponseObject, error) {
	value, err := h.projects.Get(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}

	return generated.GetProject200JSONResponse(projectToAPI(value)), nil
}

func (h *ProjectHandler) GetProjectSummary(
	ctx context.Context,
	request generated.GetProjectSummaryRequestObject,
) (generated.GetProjectSummaryResponseObject, error) {
	value, err := h.summary.Do(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}

	upcoming := lo.Map(
		value.UpcomingTasks,
		func(task projectsummary.UpcomingTask, _ int) generated.UpcomingTask {
			return generated.UpcomingTask{
				ID:         task.ID,
				Key:        task.Key,
				Title:      task.Title,
				StatusID:   task.StatusID,
				PriorityID: task.PriorityID,
				DueDate:    openapiTypes.Date{Time: task.DueDate},
			}
		},
	)

	return generated.GetProjectSummary200JSONResponse(generated.ProjectSummary{
		TasksTotal:            value.TasksTotal,
		TasksOpen:             value.TasksOpen,
		TasksDueSoon:          value.TasksDueSoon,
		TasksOverdue:          value.TasksOverdue,
		TotalTimeSpentSeconds: value.TotalTimeSpentSeconds,
		TotalEstimateSeconds:  value.TotalEstimateSeconds,
		UpcomingTasks:         upcoming,
	}), nil
}

func (h *ProjectHandler) UpdateProject(
	ctx context.Context,
	request generated.UpdateProjectRequestObject,
) (generated.UpdateProjectResponseObject, error) {
	value, err := h.projects.Update(ctx, request.ProjectID, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return generated.UpdateProject200JSONResponse(projectToAPI(value)), nil
}

func (h *ProjectHandler) ListProjects(
	ctx context.Context,
	request generated.ListProjectsRequestObject,
) (generated.ListProjectsResponseObject, error) {
	limit, offset := *request.Params.Limit, *request.Params.Offset

	values, err := h.projects.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	items := lo.Map(values, func(value domain.Project, _ int) generated.Project {
		return projectToAPI(value)
	})

	return generated.ListProjects200JSONResponse(generated.ProjectList{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}), nil
}

func (h *ProjectHandler) DeleteProject(
	ctx context.Context,
	request generated.DeleteProjectRequestObject,
) (generated.DeleteProjectResponseObject, error) {
	if err := h.projects.Delete(ctx, request.ProjectID); err != nil {
		return nil, err
	}

	return generated.DeleteProject204Response{}, nil
}

func projectToAPI(value domain.Project) generated.Project {
	return generated.Project{
		ID:   value.ID,
		Key:  value.Key,
		Name: value.Name,
	}
}
