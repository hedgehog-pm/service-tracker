package listboardcolumntasks

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/samber/lo"

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

	if fingerprint := filterFingerprint(params.Filters); fingerprint != "" {
		value += ":" + fingerprint
	}

	return base64.RawURLEncoding.EncodeToString([]byte(value))
}

func decodeCursor(value string, params Params) (int, int, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "invalid cursor")
	}

	parts := strings.Split(string(raw), ":")
	if len(parts) < 5 || len(parts) > 6 {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "invalid cursor")
	}

	var values [5]int
	for i, part := range parts[:5] {
		values[i], err = strconv.Atoi(part)
		if err != nil {
			return 0, 0, domain.NewError(domain.CodeInvalidInput, "invalid cursor")
		}
	}

	if values[0] != params.ProjectID || values[1] != params.BoardID || values[2] != params.ColumnID || values[4] <= 0 {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "cursor does not belong to this column")
	}

	var fingerprint string
	if len(parts) == 6 {
		fingerprint = parts[5]
	}

	if fingerprint != filterFingerprint(params.Filters) {
		return 0, 0, domain.NewError(domain.CodeInvalidInput, "cursor does not match the filters")
	}

	return values[3], values[4], nil
}

func filterFingerprint(filters Filters) string {
	if filters == (Filters{}) {
		return ""
	}

	value := fmt.Sprintf(
		"%d:%d:%d:%d:%s",
		lo.FromPtr(filters.StatusID),
		lo.FromPtr(filters.ModuleID),
		lo.FromPtr(filters.PriorityID),
		lo.FromPtr(filters.TaskTypeID),
		filters.Query,
	)
	sum := sha256.Sum256([]byte(value))

	return hex.EncodeToString(sum[:])
}
