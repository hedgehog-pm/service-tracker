package board

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	boards boardRepository
}

func New(boards boardRepository) *Usecase { return &Usecase{boards: boards} }

func (u *Usecase) Get(ctx context.Context, projectID, boardID int) (domain.Board, error) {
	return u.boards.GetBoard(ctx, projectID, boardID)
}

func (u *Usecase) Create(ctx context.Context, projectID int, name string) (domain.Board, error) {
	value := domain.Board{ProjectID: projectID, Name: name}
	if err := value.Validate(); err != nil {
		return domain.Board{}, err
	}

	return u.boards.InsertBoard(ctx, value)
}

func (u *Usecase) Update(
	ctx context.Context,
	projectID, boardID int,
	name string,
) (domain.Board, error) {
	value := domain.Board{ID: boardID, ProjectID: projectID, Name: name}
	if err := value.Validate(); err != nil {
		return domain.Board{}, err
	}

	return u.boards.UpdateBoard(ctx, value)
}

func (u *Usecase) Delete(ctx context.Context, projectID, boardID int) error {
	return u.boards.DeleteBoard(ctx, projectID, boardID)
}
