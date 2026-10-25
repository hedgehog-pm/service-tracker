package domain

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// TaskFields объединяет входные поля создания и обновления.
type TaskFields struct {
	Title       string
	Description string
	StatusID    int
	ModuleID    *int
}

func (f TaskFields) Validate() error {
	err := validation.ValidateStruct(&f,
		validation.Field(&f.Title, validation.By(notBlank)),
	)
	if err != nil {
		return NewError(CodeInvalidInput, err.Error())
	}

	return nil
}

type Task struct {
	ID         int
	ProjectID  int
	ProjectKey string
	Number     int
	TaskFields
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t Task) Key() string {
	return fmt.Sprintf("%s-%d", t.ProjectKey, t.Number)
}
