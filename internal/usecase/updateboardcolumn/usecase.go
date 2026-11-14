package updateboardcolumn

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Params struct {
		ProjectID int
		BoardID   int
		ColumnID  int
		Name      string
		Position  int
		StatusIDs []int
	}
	Usecase struct {
		boards     boardRepository
		statuses   statusReader
		transactor transactor
	}
)

func New(boards boardRepository, statuses statusReader, transactor transactor) *Usecase {
	return &Usecase{boards: boards, statuses: statuses, transactor: transactor}
}

func (u *Usecase) Do(ctx context.Context, params Params) (domain.BoardColumn, error) {
	var column domain.BoardColumn

	err := u.transactor.InTx(ctx, func(ctx context.Context) error {
		board, err := u.boards.GetBoardForUpdate(ctx, params.ProjectID, params.BoardID)
		if err != nil {
			return err
		}

		available, err := u.statuses.ListProjectStatuses(ctx, params.ProjectID)
		if err != nil {
			return err
		}

		column = domain.BoardColumn{
			ID:        params.ColumnID,
			Name:      params.Name,
			Position:  params.Position,
			StatusIDs: params.StatusIDs,
		}
		if err := board.UpdateColumn(column, available); err != nil {
			return err
		}

		board, err = u.boards.SaveColumns(ctx, board)
		if err != nil {
			return err
		}

		column = board.Columns[params.Position]

		return nil
	})
	if err != nil {
		return domain.BoardColumn{}, err
	}

	return column, nil
}
