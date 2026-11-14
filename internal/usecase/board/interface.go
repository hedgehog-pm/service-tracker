package board

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type boardRepository interface {
	GetBoard(ctx context.Context, projectID, boardID int) (domain.Board, error)
	InsertBoard(ctx context.Context, board domain.Board) (domain.Board, error)
	UpdateBoard(ctx context.Context, board domain.Board) (domain.Board, error)
	DeleteBoard(ctx context.Context, projectID, boardID int) error
}
