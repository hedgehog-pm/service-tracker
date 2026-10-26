package api

import (
	"net/http"

	"github.com/labstack/echo/v5"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/issuetypes"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/issuetype"
)

type IssueTypeHandler struct{ issueTypes *issuetype.Usecase }

func NewIssueTypeHandler(issueTypes *issuetype.Usecase) *IssueTypeHandler {
	return &IssueTypeHandler{issueTypes: issueTypes}
}

func (h *IssueTypeHandler) ListIssueTypes(ctx *echo.Context, projectID int) error {
	values, err := h.issueTypes.List(ctx.Request().Context(), projectID)
	if err != nil {
		return err
	}

	items := make([]generated.IssueType, 0, len(values))
	for _, value := range values {
		items = append(items, generated.IssueType{ID: value.ID, Name: value.Name, Position: value.Position, InitialStatus: statusToAPI(value.InitialStatus)})
	}

	return ctx.JSON(http.StatusOK, generated.IssueTypeList{Items: items})
}
