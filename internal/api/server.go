package api

import (
	"fmt"

	"github.com/labstack/echo/v5"
	echomiddleware "github.com/oapi-codegen/echo-v5-middleware"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
)

func NewServer(projectHandler projects.ServerInterface, statusHandler statuses.ServerInterface, moduleHandler modules.ServerInterface, taskHandler tasks.ServerInterface) (*echo.Echo, error) {
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

	server := echo.New()
	server.HTTPErrorHandler = HTTPErrorHandler
	projects.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(projectSpec)), projectHandler)
	statuses.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(statusSpec)), statusHandler)
	modules.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(moduleSpec)), moduleHandler)
	tasks.RegisterHandlers(server.Group("", echomiddleware.OapiRequestValidator(taskSpec)), taskHandler)

	return server, nil
}
