package api

import (
	"context"

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
	return &ModuleHandler{
		modules:    modules,
		listModule: listModule,
	}
}

func (h *ModuleHandler) CreateModule(
	ctx context.Context,
	request generated.CreateModuleRequestObject,
) (generated.CreateModuleResponseObject, error) {
	value, err := h.modules.Create(ctx, request.ProjectID, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return generated.CreateModule201JSONResponse(moduleToAPI(value)), nil
}

func (h *ModuleHandler) GetModule(
	ctx context.Context,
	request generated.GetModuleRequestObject,
) (generated.GetModuleResponseObject, error) {
	value, err := h.modules.Get(ctx, request.ProjectID, request.ModuleID)
	if err != nil {
		return nil, err
	}

	return generated.GetModule200JSONResponse(moduleToAPI(value)), nil
}

func (h *ModuleHandler) UpdateModule(
	ctx context.Context,
	request generated.UpdateModuleRequestObject,
) (generated.UpdateModuleResponseObject, error) {
	value, err := h.modules.Update(ctx, request.ProjectID, request.ModuleID, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return generated.UpdateModule200JSONResponse(moduleToAPI(value)), nil
}

func (h *ModuleHandler) ListModules(
	ctx context.Context,
	request generated.ListModulesRequestObject,
) (generated.ListModulesResponseObject, error) {
	limit, offset := *request.Params.Limit, *request.Params.Offset

	values, err := h.listModule.Do(ctx, request.ProjectID, limit, offset)
	if err != nil {
		return nil, err
	}

	items := make([]generated.Module, 0, len(values))
	for _, value := range values {
		items = append(items, moduleToAPI(value))
	}

	return generated.ListModules200JSONResponse(generated.ModuleList{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}), nil
}

func (h *ModuleHandler) DeleteModule(
	ctx context.Context,
	request generated.DeleteModuleRequestObject,
) (generated.DeleteModuleResponseObject, error) {
	if err := h.modules.Delete(ctx, request.ProjectID, request.ModuleID); err != nil {
		return nil, err
	}

	return generated.DeleteModule204Response{}, nil
}

func moduleToAPI(value domain.Module) generated.Module {
	return generated.Module{
		ID:        value.ID,
		ProjectID: value.ProjectID,
		Name:      value.Name,
	}
}
