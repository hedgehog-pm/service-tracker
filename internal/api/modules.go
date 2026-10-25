package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listmodules"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/module"
)

type ModuleHandler struct {
	modules    *module.Usecase
	listModule *listmodules.Usecase
}

func NewModuleHandler(modules *module.Usecase, listModule *listmodules.Usecase) *ModuleHandler {
	return &ModuleHandler{modules: modules, listModule: listModule}
}

func (h *ModuleHandler) CreateModule(ctx *echo.Context, projectID int) error {
	var body generated.CreateModule
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.modules.Create(ctx.Request().Context(), projectID, body.Name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, moduleToAPI(value))
}

func (h *ModuleHandler) GetModule(ctx *echo.Context, projectID int, id int) error {
	value, err := h.modules.Get(ctx.Request().Context(), projectID, id)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, moduleToAPI(value))
}

func (h *ModuleHandler) UpdateModule(ctx *echo.Context, projectID int, id int) error {
	var body generated.UpdateModule
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	value, err := h.modules.Update(ctx.Request().Context(), projectID, id, body.Name)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, moduleToAPI(value))
}

func (h *ModuleHandler) ListModules(ctx *echo.Context, projectID int, params generated.ListModulesParams) error {
	limit, offset := pagination(params.Limit, params.Offset)

	values, err := h.listModule.Do(ctx.Request().Context(), projectID, limit, offset)
	if err != nil {
		return err
	}

	items := make([]generated.Module, 0, len(values))
	for _, value := range values {
		items = append(items, moduleToAPI(value))
	}

	return ctx.JSON(http.StatusOK, generated.ModuleList{Items: items, Limit: limit, Offset: offset})
}

func (h *ModuleHandler) DeleteModule(ctx *echo.Context, projectID, id int) error {
	if err := h.modules.Delete(ctx.Request().Context(), projectID, id); err != nil {
		return err
	}

	return ctx.NoContent(http.StatusNoContent)
}

func moduleToAPI(value domain.Module) generated.Module {
	return generated.Module{ID: value.ID, ProjectID: value.ProjectID, Name: value.Name}
}
