package listboardcolumntasks

import (
	"context"
	"strings"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type (
	Params struct {
		ProjectID int
		BoardID   int
		ColumnID  int
		Filters   Filters
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
	params.Filters.Query = strings.TrimSpace(params.Filters.Query)

	board, err := u.boards.GetBoard(ctx, params.ProjectID, params.BoardID)
	if err != nil {
		return Page{}, err
	}

	if _, err := board.Column(params.ColumnID); err != nil {
		return Page{}, err
	}

	query := Query{
		ProjectID: params.ProjectID,
		BoardID:   params.BoardID,
		ColumnID:  params.ColumnID,
		Filters:   params.Filters,
		Limit:     params.Limit + 1,
	}

	if params.Cursor != nil {
		rank, number, err := decodeCursor(*params.Cursor, params)
		if err != nil {
			return Page{}, err
		}

		query.AfterRank, query.AfterNumber = &rank, &number
	}

	tasks, err := u.tasks.ListBoardColumnTasks(ctx, query)
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
