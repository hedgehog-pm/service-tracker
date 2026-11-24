package listboardcolumntasks

import (
	"context"
	"strings"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	boards boardReader
	tasks  taskReader
}

func New(
	boards boardReader,
	tasks taskReader,
) *Usecase {
	return &Usecase{
		boards: boards,
		tasks:  tasks,
	}
}

func (u *Usecase) Do(ctx context.Context, query Query) ([]domain.Task, error) {
	board, err := u.boards.GetBoard(ctx, query.ProjectID, query.BoardID)
	if err != nil {
		return nil, err
	}

	if _, err := board.Column(query.ColumnID); err != nil {
		return nil, err
	}

	query.Filters.Query = strings.TrimSpace(query.Filters.Query)

	return u.tasks.ListBoardColumnTasks(ctx, query)
}
