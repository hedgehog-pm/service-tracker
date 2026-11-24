package e2e

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/taskworklogs"
)

type WorklogQuery struct {
	From   time.Time
	To     time.Time
	Search string
	Limit  int
	Offset int
}

func (task *Task) AddWorklog(
	t testing.TB,
	startedAt time.Time,
	duration time.Duration,
) taskworklogs.TaskWorklog {
	t.Helper()

	response, err := request[taskworklogs.TaskWorklog](
		t.Context(),
		task.project.client,
		http.MethodPost,
		task.path()+"/worklogs",
		taskworklogs.CreateTaskWorklog{
			StartedAt:        startedAt,
			TimeSpentSeconds: int(duration / time.Second),
		},
	)

	return requireResponse(t, http.StatusCreated, response, err)
}

func (p Project) Worklogs(
	ctx context.Context,
	query WorklogQuery,
) (Response[taskworklogs.ProjectWorklogList], error) {
	path := p.path + "/worklogs" + query.suffix()

	return request[taskworklogs.ProjectWorklogList](ctx, p.client, http.MethodGet, path, nil)
}

func (p Project) WorklogsByTask(
	ctx context.Context,
	query WorklogQuery,
) (Response[taskworklogs.ProjectWorklogByTaskList], error) {
	path := p.path + "/worklogs/by-task" + query.suffix()

	return request[taskworklogs.ProjectWorklogByTaskList](ctx, p.client, http.MethodGet, path, nil)
}

func (q WorklogQuery) suffix() string {
	params := url.Values{}

	if !q.From.IsZero() {
		params.Set("from", q.From.Format(time.RFC3339))
	}

	if !q.To.IsZero() {
		params.Set("to", q.To.Format(time.RFC3339))
	}

	if q.Search != "" {
		params.Set("q", q.Search)
	}

	if q.Limit != 0 {
		params.Set("limit", strconv.Itoa(q.Limit))
	}

	if q.Offset != 0 {
		params.Set("offset", strconv.Itoa(q.Offset))
	}

	if len(params) == 0 {
		return ""
	}

	return "?" + params.Encode()
}
