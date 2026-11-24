package listboardcolumntasks

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

func encodeCursor(params Params, task domain.Task) string {
	value := fmt.Sprintf(
		"%d:%d:%d:%d:%d",
		params.ProjectID,
		params.BoardID,
		params.ColumnID,
		task.Rank,
		task.Number,
	)

	return base64.RawURLEncoding.EncodeToString([]byte(value))
}

func decodeCursor(value string, params Params) (int, int, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "invalid cursor")
	}

	parts := strings.Split(string(raw), ":")
	if len(parts) != 5 {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "invalid cursor")
	}

	var values [5]int
	for i, part := range parts {
		values[i], err = strconv.Atoi(part)
		if err != nil {
			return 0, 0, domain.NewError(domain.CodeInvalidInput, "invalid cursor")
		}
	}

	if values[0] != params.ProjectID || values[1] != params.BoardID || values[2] != params.ColumnID || values[4] <= 0 {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "cursor does not belong to this column")
	}

	return values[3], values[4], nil
}
