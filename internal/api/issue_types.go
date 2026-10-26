package api

import (
	"context"

	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/issuetypes"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/issuetype"
)

type IssueTypeHandler struct{ issueTypes *issuetype.Usecase }

func NewIssueTypeHandler(issueTypes *issuetype.Usecase) *IssueTypeHandler {
	return &IssueTypeHandler{issueTypes: issueTypes}
}

func (h *IssueTypeHandler) ListIssueTypes(
	ctx context.Context,
	request generated.ListIssueTypesRequestObject,
) (generated.ListIssueTypesResponseObject, error) {
	values, err := h.issueTypes.List(ctx, request.ProjectID)
	if err != nil {
		return nil, err
	}

	items := make([]generated.IssueType, 0, len(values))
	for _, value := range values {
		items = append(items, generated.IssueType{
			ID:            value.ID,
			Name:          value.Name,
			Position:      value.Position,
			InitialStatus: statusToAPI(value.InitialStatus),
		})
	}

	return generated.ListIssueTypes200JSONResponse(generated.IssueTypeList{Items: items}), nil
}
