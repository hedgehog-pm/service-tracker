package status

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/status"
)

func statusToDomain(row bobgen.ListStatusesRow) domain.Status {
	return domain.Status{ID: row.ID, Name: row.Name, Category: row.Category}
}
