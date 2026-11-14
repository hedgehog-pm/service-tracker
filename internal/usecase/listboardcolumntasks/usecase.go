package listboardcolumntasks

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Params struct {
		ProjectID int
		BoardID   int
		ColumnID  int
		Limit     int
		Cursor    *string
	}
	Page struct {
		Tasks      []domain.Task
		NextCursor *string
	}
	Usecase struct {
		boards boardReader
		tasks  taskReader
	}
)

func New(
	boards boardReader,
	tasks taskReader,
) *Usecase {
	return &Usecase{
		boards: boards,
		tasks:  tasks,
	}
}

func (u *Usecase) Do(ctx context.Context, params Params) (Page, error) {
	board, err := u.boards.GetBoard(ctx, params.ProjectID, params.BoardID)
	if err != nil {
		return Page{}, err
	}

	if _, err := board.Column(params.ColumnID); err != nil {
		return Page{}, err
	}

	var afterRank, afterNumber *int

	if params.Cursor != nil {
		rank, number, err := decodeCursor(*params.Cursor, params)
		if err != nil {
			return Page{}, err
		}

		afterRank, afterNumber = &rank, &number
	}

	tasks, err := u.tasks.ListBoardColumnTasks(
		ctx,
		params.ProjectID,
		params.BoardID,
		params.ColumnID,
		afterRank,
		afterNumber,
		params.Limit+1,
	)
	if err != nil {
		return Page{}, err
	}

	page := Page{Tasks: tasks}
	if len(tasks) > params.Limit {
		page.Tasks = tasks[:params.Limit]
		last := page.Tasks[len(page.Tasks)-1]
		page.NextCursor = new(encodeCursor(params, last))
	}

	return page, nil
}
