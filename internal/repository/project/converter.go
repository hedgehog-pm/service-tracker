package project

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/project"
)

func projectToDomain(row bobgen.GetProjectRow) domain.Project {
	return domain.Project{ID: row.ID, Key: row.Key, Name: row.Name, NextNumber: row.NextNumber, WorkflowSchemeID: row.WorkflowSchemeID, PrioritySchemeID: row.PrioritySchemeID}
}
