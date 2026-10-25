package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/common"
)

func HTTPErrorHandler(ctx *echo.Context, err error) {
	if response, unwrapErr := echo.UnwrapResponse(ctx.Response()); unwrapErr == nil && response.Committed {
		return
	}

	status := http.StatusInternalServerError

	body := common.Error{Code: "internal_error", Error: "internal error"}
	if problem, ok := errors.AsType[*domain.Error](err); ok {
		switch problem.Code {
		case domain.CodeInvalidInput:
			status = http.StatusBadRequest
		case domain.CodeNotFound:
			status = http.StatusNotFound
		case domain.CodeConflict:
			status = http.StatusConflict
		case domain.CodeInternalError:
			status = http.StatusInternalServerError
		}

		if status < http.StatusInternalServerError {
			body = common.Error{Code: string(problem.Code), Error: problem.Message}
		}
	} else if httpStatus := echo.StatusCode(err); httpStatus >= 400 && httpStatus < 500 {
		status = httpStatus

		body = common.Error{Code: "invalid_input", Error: http.StatusText(status)}
		if status == http.StatusNotFound {
			body.Code = "not_found"
		}
	}

	if status >= http.StatusInternalServerError && !errors.Is(err, domain.ErrInternalServerError) {
		slog.ErrorContext(ctx.Request().Context(), "request failed", "error", err)
	}

	if writeErr := ctx.JSON(status, body); writeErr != nil {
		slog.ErrorContext(ctx.Request().Context(), "write error response", "error", writeErr)
	}
}
