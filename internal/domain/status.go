package domain

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Status struct {
	ID        int
	ProjectID int
	Name      string
	Category  string
}

func (e Status) Validate() error {
	if err := validation.ValidateStruct(&e, validation.Field(&e.Name, validation.By(notBlank))); err != nil {
		return NewError(CodeInvalidInput, err.Error())
	}

	return nil
}
