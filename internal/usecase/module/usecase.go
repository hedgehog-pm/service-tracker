package module

import (
	"context"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Usecase struct {
	modules moduleRepository
}

func New(modules moduleRepository) *Usecase {
	return &Usecase{modules: modules}
}

func (u *Usecase) Get(ctx context.Context, projectID, id int) (domain.Module, error) {
	return u.modules.GetModule(ctx, projectID, id)
}

func (u *Usecase) Create(ctx context.Context, projectID int, name string) (domain.Module, error) {
	value := domain.Module{
		ProjectID: projectID,
		Name:      name,
	}
	if err := value.Validate(); err != nil {
		return domain.Module{}, err
	}

	return u.modules.InsertModule(ctx, value)
}

func (u *Usecase) Update(
	ctx context.Context,
	projectID, id int,
	name string,
) (domain.Module, error) {
	value := domain.Module{
		ID:        id,
		ProjectID: projectID,
		Name:      name,
	}
	if err := value.Validate(); err != nil {
		return domain.Module{}, err
	}

	return u.modules.UpdateModule(ctx, value)
}

func (u *Usecase) Delete(ctx context.Context, projectID, id int) error {
	return u.modules.DeleteModule(ctx, projectID, id)
}
