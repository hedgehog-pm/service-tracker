package e2e

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/boards"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/statuses"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasktypes"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/transitions"
)

type Project struct {
	client *Client
	path   string
}

type Board struct {
	project Project
	path    string
}

type Column struct {
	board Board
	model boards.BoardColumn
	path  string
}

type ColumnSpec struct {
	Name     string
	Statuses []string
}

type MoveOptions struct {
	Before     *tasks.Task
	Transition string
}

func (c *Client) CreateProject(t testing.TB, body projects.CreateProject) Project {
	t.Helper()

	response, err := request[projects.Project](t.Context(), c, http.MethodPost, "/projects", body)
	created := requireResponse(t, http.StatusCreated, response, err)

	return c.Project(created.ID)
}

func (c *Client) Project(id int) Project {
	return Project{
		client: c,
		path:   "/projects/" + strconv.Itoa(id),
	}
}

func (p Project) Boards(ctx context.Context) (Response[boards.BoardList], error) {
	return request[boards.BoardList](ctx, p.client, http.MethodGet, p.path+"/boards", nil)
}

func (p Project) ExpectOnlyBoard(t testing.TB) Board {
	t.Helper()

	response, err := p.Boards(t.Context())
	list := requireResponse(t, http.StatusOK, response, err)

	require.Len(t, list.Items, 1, "%s", response)

	return p.Board(list.Items[0].ID)
}

func (p Project) CreateBoard(ctx context.Context, body boards.CreateBoard) (Response[boards.Board], error) {
	return request[boards.Board](ctx, p.client, http.MethodPost, p.path+"/boards", body)
}

func (p Project) Board(id int) Board {
	return Board{
		project: p,
		path:    p.path + "/boards/" + strconv.Itoa(id),
	}
}

func (p Project) Statuses(ctx context.Context) (Response[statuses.StatusList], error) {
	return request[statuses.StatusList](ctx, p.client, http.MethodGet, p.path+"/statuses", nil)
}

func (p Project) TaskTypes(ctx context.Context) (Response[tasktypes.TaskTypeList], error) {
	return request[tasktypes.TaskTypeList](ctx, p.client, http.MethodGet, p.path+"/task-types", nil)
}

func (p Project) TaskType(t testing.TB, name string) tasktypes.TaskType {
	t.Helper()

	response, err := p.TaskTypes(t.Context())
	list := requireResponse(t, http.StatusOK, response, err)

	index := slices.IndexFunc(list.Items, func(taskType tasktypes.TaskType) bool {
		return taskType.Name == name
	})
	require.GreaterOrEqual(t, index, 0, "task type %q not found\n%s", name, response)

	return list.Items[index]
}

func (p Project) CreateTask(t testing.TB, body tasks.CreateTask) tasks.Task {
	t.Helper()

	response, err := request[tasks.Task](t.Context(), p.client, http.MethodPost, p.path+"/tasks", body)

	return requireResponse(t, http.StatusCreated, response, err)
}

func (p Project) Transitions(ctx context.Context, taskID int) (Response[transitions.TransitionList], error) {
	path := p.path + "/tasks/" + strconv.Itoa(taskID) + "/transitions"

	return request[transitions.TransitionList](ctx, p.client, http.MethodGet, path, nil)
}

func (p Project) Transition(t testing.TB, taskID int, name string) transitions.Transition {
	t.Helper()

	response, err := p.Transitions(t.Context(), taskID)
	list := requireResponse(t, http.StatusOK, response, err)

	index := slices.IndexFunc(list.Items, func(transition transitions.Transition) bool {
		return transition.Name == name
	})
	require.GreaterOrEqual(t, index, 0, "transition %q unavailable for task %d\n%s", name, taskID, response)

	return list.Items[index]
}

func (b Board) Get(ctx context.Context) (Response[boards.Board], error) {
	return request[boards.Board](ctx, b.project.client, http.MethodGet, b.path, nil)
}

func (b Board) ExpectLayout(t testing.TB, name string, columns ...ColumnSpec) map[string]Column {
	t.Helper()

	ctx := t.Context()
	response, err := b.Get(ctx)
	board := requireResponse(t, http.StatusOK, response, err)

	require.Equal(t, name, board.Name, "board name")
	require.Len(t, board.Columns, len(columns), "columns of board %q", name)

	statusResponse, err := b.project.Statuses(ctx)
	statuses := requireResponse(t, http.StatusOK, statusResponse, err)
	statusNames := make(map[int]string, len(statuses.Items))

	for _, status := range statuses.Items {
		statusNames[status.ID] = status.Name
	}

	result := make(map[string]Column, len(columns))

	for position, spec := range columns {
		column := board.Columns[position]
		require.Equal(t, spec.Name, column.Name, "column at index %d", position)
		require.Equal(t, position, column.Position, "position of column %q", column.Name)

		names := make([]string, len(column.StatusIds))

		for index, statusID := range column.StatusIds {
			names[index] = statusNames[statusID]
		}

		require.Equal(t, spec.Statuses, names, "statuses of column %q", column.Name)

		result[spec.Name] = b.Column(column)
	}

	return result
}

func (b Board) Delete(ctx context.Context) (Response[struct{}], error) {
	return request[struct{}](ctx, b.project.client, http.MethodDelete, b.path, nil)
}

func (b Board) CreateColumn(ctx context.Context, body boards.CreateBoardColumn) (Response[boards.BoardColumn], error) {
	return request[boards.BoardColumn](ctx, b.project.client, http.MethodPost, b.path+"/columns", body)
}

func (b Board) Column(model boards.BoardColumn) Column {
	return Column{
		board: b,
		model: model,
		path:  b.path + "/columns/" + strconv.Itoa(model.ID),
	}
}

func (b Board) Move(ctx context.Context, taskID int, body boards.MoveBoardTask) (Response[tasks.Task], error) {
	path := b.path + "/tasks/" + strconv.Itoa(taskID) + "/move"

	return request[tasks.Task](ctx, b.project.client, http.MethodPost, path, body)
}

func (c Column) Move(t testing.TB, task tasks.Task, options MoveOptions) {
	t.Helper()

	body := boards.MoveBoardTask{
		ColumnID: c.model.ID,
	}

	if options.Before != nil {
		body.BeforeTaskID = &options.Before.ID
	}

	if options.Transition != "" {
		transition := c.board.project.Transition(t, task.ID, options.Transition)
		body.TransitionID = &transition.ID
	}

	response, err := c.board.Move(t.Context(), task.ID, body)
	requireResponse(t, http.StatusOK, response, err)
}

func (c Column) ExpectTasks(t testing.TB, want ...tasks.Task) {
	t.Helper()

	response, err := c.Tasks(t.Context(), boards.ListBoardColumnTasksParams{})
	list := requireResponse(t, http.StatusOK, response, err)

	require.Len(t, list.Items, len(want), "tasks in column %q\n%s", c.model.Name, response)

	for index, task := range list.Items {
		require.Equal(t, want[index].ID, task.ID, "column %q at index %d: got %s, want %s",
			c.model.Name, index, task.Key, want[index].Key)
		require.Contains(t, c.model.StatusIds, task.StatusID, "status of task %s in column %q",
			task.Key, c.model.Name)
	}
}

func (c Column) Update(ctx context.Context, body boards.UpdateBoardColumn) (Response[boards.BoardColumn], error) {
	return request[boards.BoardColumn](ctx, c.board.project.client, http.MethodPut, c.path, body)
}

func (c Column) Delete(ctx context.Context) (Response[struct{}], error) {
	return request[struct{}](ctx, c.board.project.client, http.MethodDelete, c.path, nil)
}

func (c Column) Tasks(
	ctx context.Context, params boards.ListBoardColumnTasksParams,
) (Response[boards.BoardColumnTaskList], error) {
	query := url.Values{}

	if params.Limit != nil {
		query.Set("limit", strconv.Itoa(*params.Limit))
	}

	if params.Cursor != nil {
		query.Set("cursor", *params.Cursor)
	}

	path := c.path + "/tasks"
	if len(query) != 0 {
		path += "?" + query.Encode()
	}

	return request[boards.BoardColumnTaskList](ctx, c.board.project.client, http.MethodGet, path, nil)
}
