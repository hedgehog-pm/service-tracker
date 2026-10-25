package status

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

func (u *Usecase) Get(ctx context.Context, projectID, id int) (domain.Status, error) {
	return u.repository.GetStatus(ctx, projectID, id)
}

func (u *Usecase) Create(ctx context.Context, projectID int, name, category string) (domain.Status, error) {
	value := domain.Status{ProjectID: projectID, Name: name, Category: category}
	if err := value.Validate(); err != nil {
		return domain.Status{}, err
	}

	return u.repository.InsertStatus(ctx, value)
}

func (u *Usecase) Update(ctx context.Context, projectID, id int, name, category string) (domain.Status, error) {
	value := domain.Status{ID: id, ProjectID: projectID, Name: name, Category: category}
	if err := value.Validate(); err != nil {
		return domain.Status{}, err
	}

	return u.repository.UpdateStatus(ctx, value)
}

func (u *Usecase) Delete(ctx context.Context, projectID, id int) error {
	return u.repository.DeleteStatus(ctx, projectID, id)
}
