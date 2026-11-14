package listboards

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	boards   boardReader
	projects projectReader
}

func New(boards boardReader, projects projectReader) *Usecase {
	return &Usecase{boards: boards, projects: projects}
}

func (u *Usecase) Do(ctx context.Context, projectID, limit, offset int) ([]domain.Board, error) {
	if _, err := u.projects.GetProject(ctx, projectID); err != nil {
		return nil, err
	}

	return u.boards.ListBoards(ctx, projectID, limit, offset)
}
