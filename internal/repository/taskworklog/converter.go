package taskworklog

import (
	"database/sql"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/taskworklog"
)

func worklogToDomain(row bobgen.GetTaskWorklogRow) domain.TaskWorklog {
	return domain.TaskWorklog{
		ID:               row.ID,
		TaskID:           row.TaskID,
		StartedAt:        row.StartedAt,
		TimeSpentSeconds: row.TimeSpentSeconds,
		Comment:          commentToPointer(row.Comment),
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func commentToNull(comment *string) sql.Null[string] {
	return sql.Null[string]{V: lo.FromPtr(comment), Valid: comment != nil}
}

func commentToPointer(comment sql.Null[string]) *string {
	if !comment.Valid {
		return nil
	}

	return &comment.V
}
