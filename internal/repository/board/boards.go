package board

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/scan/pgxscan"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	bobgen "github.com/hedgehog-pm/service-tracker/internal/generated/bob/board"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
)

func (r *Postgres) GetBoard(ctx context.Context, projectID, boardID int) (domain.Board, error) {
	q := bobgen.GetBoard(projectID, boardID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("build get board query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("get board: %w", err))
	}

	if len(rows) == 0 {
		return domain.Board{}, domain.ErrBoardNotFound
	}

	return boardsToDomain(rows)[0], nil
}

func (r *Postgres) GetBoardForUpdate(
	ctx context.Context,
	projectID, boardID int,
) (domain.Board, error) {
	q := bobgen.LockBoard(projectID, boardID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("build lock board query: %w", err))
	}

	_, err = pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Board{}, domain.ErrBoardNotFound
	}

	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("lock board: %w", err))
	}

	return r.GetBoard(ctx, projectID, boardID)
}

func (r *Postgres) ListBoards(
	ctx context.Context,
	projectID, limit, offset int,
) ([]domain.Board, error) {
	q := bobgen.ListBoards(projectID, limit, offset)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("build list boards query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return nil, postgres.MapError(ctx, fmt.Errorf("list boards: %w", err))
	}

	return boardsToDomain(rows), nil
}

func (r *Postgres) InsertBoard(ctx context.Context, board domain.Board) (domain.Board, error) {
	q := bobgen.InsertBoard(board.ProjectID, board.Name)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("build insert board query: %w", err))
	}

	row, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Board{}, domain.ErrProjectNotFound
	}

	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("insert board: %w", err))
	}

	board.ID = row.ID
	board.SchemeID = row.BoardSchemeID

	return board, nil
}

func (r *Postgres) InsertDefaultBoard(ctx context.Context, projectID int) error {
	q := bobgen.InsertDefaultBoard(projectID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build insert default board query: %w", err))
	}

	if _, err := pgxscan.One(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...); err != nil {
		return postgres.MapError(ctx, fmt.Errorf("insert default board for project %d: %w", projectID, err))
	}

	return nil
}

func (r *Postgres) UpdateBoard(ctx context.Context, board domain.Board) (domain.Board, error) {
	q := bobgen.UpdateBoard(board.Name, board.ProjectID, board.ID)

	query, args, err := bob.Build(ctx, q)
	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("build update board query: %w", err))
	}

	rows, err := pgxscan.All(ctx, postgres.Executor(ctx, r.db), q.Scanner, query, args...)
	if err != nil {
		return domain.Board{}, postgres.MapError(ctx, fmt.Errorf("update board: %w", err))
	}

	if len(rows) == 0 {
		return domain.Board{}, domain.ErrBoardNotFound
	}

	return boardsToDomain(rows)[0], nil
}

func (r *Postgres) DeleteBoard(ctx context.Context, projectID, boardID int) error {
	query, args, err := bob.Build(ctx, bobgen.DeleteBoard(projectID, boardID))
	if err != nil {
		return postgres.MapError(ctx, fmt.Errorf("build delete board query: %w", err))
	}

	result, err := postgres.Executor(ctx, r.db).Exec(ctx, query, args...)
	if err != nil {
		return postgres.MapDeleteError(ctx, fmt.Errorf("delete board: %w", err))
	}

	if result.RowsAffected() == 0 {
		return domain.ErrBoardNotFound
	}

	return nil
}
