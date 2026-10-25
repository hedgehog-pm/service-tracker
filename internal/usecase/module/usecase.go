package module

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

func (u *Usecase) Get(ctx context.Context, projectID, id int) (domain.Module, error) {
	return u.repository.GetModule(ctx, projectID, id)
}

func (u *Usecase) Create(ctx context.Context, projectID int, name string) (domain.Module, error) {
	value := domain.Module{ProjectID: projectID, Name: name}
	if err := value.Validate(); err != nil {
		return domain.Module{}, err
	}

	return u.repository.InsertModule(ctx, value)
}

func (u *Usecase) Update(ctx context.Context, projectID, id int, name string) (domain.Module, error) {
	value := domain.Module{ID: id, ProjectID: projectID, Name: name}
	if err := value.Validate(); err != nil {
		return domain.Module{}, err
	}

	return u.repository.UpdateModule(ctx, value)
}

func (u *Usecase) Delete(ctx context.Context, projectID, id int) error {
	return u.repository.DeleteModule(ctx, projectID, id)
}
