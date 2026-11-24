package board

import (
	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/board"
)

func boardsToDomain(rows []bobgen.GetBoardRow) []domain.Board {
	boards := make([]domain.Board, 0)
	for _, row := range rows {
		if len(boards) == 0 || boards[len(boards)-1].ID != row.ID {
			boards = append(boards, domain.Board{
				ID:           row.ID,
				ProjectID:    row.ProjectID,
				Name:         row.Name,
				SchemeID:     row.BoardSchemeID,
				SharedScheme: row.IsShared,
			})
		}

		if !row.ColumnID.Valid {
			continue
		}

		board := &boards[len(boards)-1]
		if len(board.Columns) == 0 || board.Columns[len(board.Columns)-1].ID != row.ColumnID.V {
			board.Columns = append(board.Columns, domain.BoardColumn{
				ID: row.ColumnID.V, Name: row.ColumnName.V, Position: row.Position,
			})
		}

		if row.StatusID.Valid {
			column := &board.Columns[len(board.Columns)-1]
			column.StatusIDs = append(column.StatusIDs, row.StatusID.V)
		}
	}

	return boards
}
