package projectworklogs

import (
	"context"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

type Query struct {
	ProjectID int
	From      *time.Time
	To        *time.Time
	Search    string
	Limit     int
	Offset    int
}

type Page[T any] struct {
	Items                 []T
	TotalTimeSpentSeconds int64
	TotalEstimateSeconds  *int64
}

type Entry struct {
	ID               int
	TaskID           int
	TaskKey          string
	TaskTitle        string
	StartedAt        time.Time
	TimeSpentSeconds int
	Comment          *string
}

type TaskTotal struct {
	TaskID           int
	TaskKey          string
	TaskTitle        string
	WorklogsCount    int64
	TimeSpentSeconds int64
}

type Usecase struct {
	worklogs worklogReader
}

func New(worklogs worklogReader) *Usecase {
	return &Usecase{worklogs: worklogs}
}

func (u *Usecase) List(ctx context.Context, query Query) (Page[Entry], error) {
	if query.From != nil &&
		query.To != nil &&
		!query.From.Before(lo.FromPtr(query.To)) {
		return Page[Entry]{}, domain.NewError(domain.CodeInvalidInput, "from must be before to")
	}

	query.Search = strings.TrimSpace(query.Search)

	return u.worklogs.ListProjectWorklogs(ctx, query)
}

func (u *Usecase) ListByTask(ctx context.Context, query Query) (Page[TaskTotal], error) {
	if query.From != nil &&
		query.To != nil &&
		!query.From.Before(lo.FromPtr(query.To)) {
		return Page[TaskTotal]{}, domain.NewError(domain.CodeInvalidInput, "from must be before to")
	}

	query.Search = strings.TrimSpace(query.Search)

	return u.worklogs.ListProjectWorklogsByTask(ctx, query)
}
