package deleteboardcolumn

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	boardRepository interface {
		GetBoardForUpdate(ctx context.Context, projectID, boardID int) (domain.Board, error)
		SaveColumns(ctx context.Context, board domain.Board) (domain.Board, error)
	}
	transactor interface {
		InTx(ctx context.Context, fn func(ctx context.Context) error) error
	}
)
