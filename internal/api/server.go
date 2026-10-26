package api

import (
	"fmt"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/oapi-codegen/echo-v5-middleware"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/issuetypes"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/priorities"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/transitions"
)

func NewServer(projectHandler projects.ServerInterface, statusHandler statuses.ServerInterface, moduleHandler modules.ServerInterface, taskHandler tasks.ServerInterface, transitionHandler transitions.ServerInterface, issueTypeHandler issuetypes.ServerInterface, priorityHandler priorities.ServerInterface) (*echo.Echo, error) {
	projectSpec, err := projects.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load projects OpenAPI contract: %w", err)
	}

	statusSpec, err := statuses.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load statuses OpenAPI contract: %w", err)
	}

	moduleSpec, err := modules.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load modules OpenAPI contract: %w", err)
	}

	taskSpec, err := tasks.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load tasks OpenAPI contract: %w", err)
	}

	transitionsSpec, err := transitions.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load transitions OpenAPI contract: %w", err)
	}

	issuetypesSpec, err := issuetypes.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load issuetypes OpenAPI contract: %w", err)
	}

	prioritiesSpec, err := priorities.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load priorities OpenAPI contract: %w", err)
	}

	server := echo.New()
	server.HTTPErrorHandler = HTTPErrorHandler
	projects.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(projectSpec)), projectHandler)
	statuses.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(statusSpec)), statusHandler)
	modules.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(moduleSpec)), moduleHandler)
	tasks.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(taskSpec)), taskHandler)

	transitions.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(transitionsSpec)), transitionHandler)
	issuetypes.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(issuetypesSpec)), issueTypeHandler)
	priorities.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(prioritiesSpec)), priorityHandler)

	return server, nil
}
