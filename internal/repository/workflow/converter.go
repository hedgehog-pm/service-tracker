package workflow

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/workflow"
)

func transitionToDomain(row bobgen.GetTransitionRow) domain.Transition {
	return domain.Transition{
		ID:           row.ID,
		WorkflowID:   row.WorkflowID,
		FromStatusID: row.FromStatusID,
		Name:         row.Name,
		ToStatus: domain.Status{
			ID:       row.ToStatusID,
			Name:     row.StatusName,
			Category: row.StatusCategory,
		},
	}
}

func initialStatusToDomain(row bobgen.GetInitialStatusRow) domain.InitialStatus {
	return domain.InitialStatus{
		WorkflowID: row.WorkflowID,
		Status: domain.Status{
			ID:       row.StatusID,
			Name:     row.StatusName,
			Category: row.StatusCategory,
		},
	}
}

func issueTypeToDomain(row bobgen.ListIssueTypesRow) domain.IssueTypeWithInitialStatus {
	return domain.IssueTypeWithInitialStatus{
		ID:       row.ID,
		Name:     row.Name,
		Position: row.Position,
		InitialStatus: domain.Status{
			ID:       row.StatusID,
			Name:     row.StatusName,
			Category: row.StatusCategory,
		},
	}
}
