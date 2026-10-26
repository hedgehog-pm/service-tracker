package domain

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Project struct {
	ID               int
	Key              string
	Name             string
	NextNumber       int
	WorkflowSchemeID int
	PrioritySchemeID int
}

func NewProject(key, name string) Project {
	return Project{
		Key:        key,
		Name:       name,
		NextNumber: 1,
	}
}

func (p Project) Validate() error {
	if err := validation.ValidateStruct(
		&p,
		validation.Field(&p.Name, validation.By(notBlank)),
	); err != nil {
		return NewError(CodeInvalidInput, err.Error())
	}

	return nil
}
