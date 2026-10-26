package api

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createproject"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/project"
)

type ProjectHandler struct {
	projects      *project.Usecase
	createProject *createproject.Usecase
}

func NewProjectHandler(
	projects *project.Usecase,
	createProject *createproject.Usecase,
) *ProjectHandler {
	return &ProjectHandler{
		projects:      projects,
		createProject: createProject,
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

	items := make([]generated.Project, 0, len(values))
	for _, value := range values {
		items = append(items, projectToAPI(value))
	}

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
