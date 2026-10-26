package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createproject"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/project"
)

type ProjectHandler struct {
	projects      *project.Usecase
	createProject *createproject.Usecase
}

func NewProjectHandler(projects *project.Usecase, createProject *createproject.Usecase) *ProjectHandler {
	return &ProjectHandler{projects: projects, createProject: createProject}
}

func (h *ProjectHandler) CreateProject(ctx *echo.Context) error {
	var body generated.CreateProject
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.createProject.Do(ctx.Request().Context(), body.Key, body.Name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, projectToAPI(value))
}

func (h *ProjectHandler) GetProject(ctx *echo.Context, projectID int) error {
	value, err := h.projects.Get(ctx.Request().Context(), projectID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, projectToAPI(value))
}

func (h *ProjectHandler) UpdateProject(ctx *echo.Context, projectID int) error {
	var body generated.UpdateProject
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.projects.Update(ctx.Request().Context(), projectID, body.Name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, projectToAPI(value))
}

func (h *ProjectHandler) ListProjects(ctx *echo.Context, params generated.ListProjectsParams) error {
	limit, offset := *params.Limit, *params.Offset

	values, err := h.projects.List(ctx.Request().Context(), limit, offset)
	if err != nil {
		return err
	}

	items := make([]generated.Project, 0, len(values))
	for _, value := range values {
		items = append(items, projectToAPI(value))
	}

	return ctx.JSON(http.StatusOK, generated.ProjectList{Items: items, Limit: limit, Offset: offset})
}

func (h *ProjectHandler) DeleteProject(ctx *echo.Context, projectID int) error {
	if err := h.projects.Delete(ctx.Request().Context(), projectID); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func projectToAPI(value domain.Project) generated.Project {
	return generated.Project{ID: value.ID, Key: value.Key, Name: value.Name}
}
