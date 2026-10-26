package module

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/module"
)

func moduleToDomain(row bobgen.GetModuleRow) domain.Module {
	return domain.Module{
		ID:        row.ID,
		ProjectID: row.ProjectID,
		Name:      row.Name,
	}
}
