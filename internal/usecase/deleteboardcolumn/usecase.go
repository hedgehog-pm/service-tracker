package deleteboardcolumn

import (
	"context"
)

type Usecase struct {
	boards     boardRepository
	transactor transactor
}

func New(boards boardRepository, transactor transactor) *Usecase {
	return &Usecase{boards: boards, transactor: transactor}
}

func (u *Usecase) Do(ctx context.Context, projectID, boardID, columnID int) error {
	return u.transactor.InTx(ctx, func(ctx context.Context) error {
		board, err := u.boards.GetBoardForUpdate(ctx, projectID, boardID)
		if err != nil {
			return err
		}

		if err := board.RemoveColumn(columnID); err != nil {
			return err
		}

		_, err = u.boards.SaveColumns(ctx, board)

		return err
	})
}
