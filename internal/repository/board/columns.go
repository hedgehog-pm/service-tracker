package board

import (
	"context"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/board"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) SaveColumns(ctx context.Context, board domain.Board) (domain.Board, error) {
	query, args, err := bob.Build(ctx, bobgen.DeleteBoardStatuses(board.ID))
	if err != nil {
		return domain.Board{}, postgres.MapError(
			ctx,
			fmt.Errorf(
				"build delete board statuses query: %w",
				err,
			),
		)
	}

	if _, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...); err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("delete board statuses: %w", err))
	}

	ids := lo.FilterMap(board.Columns, func(column domain.BoardColumn, _ int) (string, bool) {
		return strconv.Itoa(column.ID), column.ID != 0
	})

	query, args, err = bob.Build(ctx, bobgen.DeleteRemovedColumns(board.ID, ids))
	if err != nil {
		return domain.Board{}, postgres.MapError(
			ctx,
			fmt.Errorf(
				"build delete removed columns query: %w",
				err,
			),
		)
	}

	if _, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...); err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("delete removed columns: %w", err))
	}

	for i := range board.Columns {
		if err := r.saveColumn(ctx, board.ID, &board.Columns[i]); err != nil {
			return domain.Board{}, err
		}
	}

	return board, nil
}

func (r *Postgres) saveColumn(ctx context.Context, boardID int, column *domain.BoardColumn) error {
	if column.ID == 0 {
		q := bobgen.InsertColumn([]bobgen.InsertColumn_Rows{{
			BoardID:  boardID,
			Name:     column.Name,
			Position: column.Position,
		}})

		query, args, err := bob.Build(ctx, q)
		if err != nil {
			return postgres.MapError(ctx, fmt.Errorf("build insert column query: %w", err))
		}

		row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
		if err != nil {
			return postgres.MapError(ctx, fmt.Errorf("insert column: %w", err))
		}

		column.ID = row
	} else {
		query, args, err := bob.Build(
			ctx,
			bobgen.UpdateColumn(
				column.Name,
				column.Position,
				boardID,
				column.ID,
			),
		)
		if err != nil {
			return postgres.MapError(ctx, fmt.Errorf("build update column query: %w", err))
		}

		if _, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...); err != nil {
			return postgres.MapError(ctx, fmt.Errorf("update column: %w", err))
		}
	}

	return r.saveColumnStatuses(ctx, boardID, *column)
}

func (r *Postgres) saveColumnStatuses(
	ctx context.Context,
	boardID int,
	column domain.BoardColumn,
) error {
	for _, statusID := range column.StatusIDs {
		query, args, err := bob.Build(
			ctx,
			bobgen.InsertColumnStatus([]bobgen.InsertColumnStatus_Rows{{
				BoardID:  boardID,
				ColumnID: column.ID,
				StatusID: statusID,
			}}),
		)
		if err != nil {
			return postgres.MapError(ctx, fmt.Errorf("build insert column status query: %w", err))
		}

		if _, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...); err != nil {
			return postgres.MapError(ctx, fmt.Errorf("insert column status: %w", err))
		}
	}

	return nil
}
