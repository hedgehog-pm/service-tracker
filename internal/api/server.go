package api

import (
	"fmt"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/oapi-codegen/echo-v5-middleware"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/priorities"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasktypes"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/transitions"
)

func NewServer(
	projectHandler projects.StrictServerInterface,
	statusHandler statuses.StrictServerInterface,
	moduleHandler modules.StrictServerInterface,
	taskHandler tasks.StrictServerInterface,
	transitionHandler transitions.StrictServerInterface,
	taskTypeHandler tasktypes.StrictServerInterface,
	priorityHandler priorities.StrictServerInterface,
) (*echo.Echo, error) {
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

	tasktypesSpec, err := tasktypes.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load tasktypes OpenAPI contract: %w", err)
	}

	prioritiesSpec, err := priorities.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("load priorities OpenAPI contract: %w", err)
	}

	server := echo.New()
	server.HTTPErrorHandler = HTTPErrorHandler
	projects.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(projectSpec)),
		projects.NewStrictHandler(projectHandler, nil),
	)
	statuses.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(statusSpec)),
		statuses.NewStrictHandler(statusHandler, nil),
	)
	modules.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(moduleSpec)),
		modules.NewStrictHandler(moduleHandler, nil),
	)
	tasks.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(taskSpec)),
		tasks.NewStrictHandler(taskHandler, nil),
	)

	transitions.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(transitionsSpec)),
		transitions.NewStrictHandler(transitionHandler, nil),
	)
	tasktypes.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(tasktypesSpec)),
		tasktypes.NewStrictHandler(taskTypeHandler, nil),
	)
	priorities.RegisterHandlers(
		server.Group("", echomiddleware.OapiRequestValidator(prioritiesSpec)),
		priorities.NewStrictHandler(priorityHandler, nil),
	)

	return server, nil
}
