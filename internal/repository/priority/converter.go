package priority

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/priority"
)

func priorityToDomain(row bobgen.ListPrioritiesRow) domain.Priority {
	return domain.Priority{ID: row.ID, Name: row.Name, Position: row.Position}
}
