package project

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	repository repository
}

func New(repository repository) *Usecase {
	return &Usecase{repository: repository}
}

func (u *Usecase) Get(ctx context.Context, id int) (domain.Project, error) {
	return u.repository.GetProject(ctx, id)
}

func (u *Usecase) List(ctx context.Context, limit int, offset int) ([]domain.Project, error) {
	return u.repository.ListProjects(ctx, limit, offset)
}

func (u *Usecase) Update(ctx context.Context, id int, name string) (domain.Project, error) {
	value := domain.Project{ID: id, Name: name}
	if err := value.Validate(); err != nil {
		return domain.Project{}, err
	}

	return u.repository.UpdateProject(ctx, value)
}

func (u *Usecase) Delete(ctx context.Context, id int) error {
	return u.repository.DeleteProject(ctx, id)
}
