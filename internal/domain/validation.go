package domain

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func notBlank(value any) error {
	if strings.TrimSpace(value.(string)) == "" {
		return validation.NewError("required", "must not be blank")
	}

	return nil
}
