package e2e

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	openapiTypes "github.com/oapi-codegen/runtime/types"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/transitions"
)

type Task struct {
	tasks.Task
	project Project
}

type TaskSpec struct {
	Title           string
	TaskTypeID      int
	DueDate         time.Time
	EstimateSeconds *int
}

func (p Project) AddTask(t testing.TB, spec TaskSpec) *Task {
	t.Helper()

	taskTypeID := spec.TaskTypeID
	if taskTypeID == 0 {
		taskTypeID = p.TaskType(t, "Task").ID
	}

	body := tasks.CreateTask{
		Title:           spec.Title,
		TaskTypeID:      taskTypeID,
		EstimateSeconds: spec.EstimateSeconds,
	}
	if !spec.DueDate.IsZero() {
		body.DueDate = &openapiTypes.Date{Time: spec.DueDate}
	}

	return &Task{Task: p.CreateTask(t, body), project: p}
}

func (task *Task) AssignToModule(t testing.TB, module modules.Module) {
	t.Helper()

	body := tasks.UpdateTask{
		Title:           task.Title,
		Description:     task.Description,
		ModuleID:        &module.ID,
		DueDate:         task.DueDate,
		PriorityID:      task.PriorityID,
		EstimateSeconds: task.EstimateSeconds,
	}
	response, err := request[tasks.Task](
		t.Context(),
		task.project.client,
		http.MethodPut,
		task.path(),
		body,
	)
	task.Task = requireResponse(t, http.StatusOK, response, err)
}

func (task *Task) Transition(t testing.TB, names ...string) {
	t.Helper()

	path := task.path() + "/transitions"

	for _, name := range names {
		transition := task.project.Transition(t, task.ID, name)
		response, err := request[tasks.Task](
			t.Context(),
			task.project.client,
			http.MethodPost,
			path,
			transitions.TransitionTask{TransitionID: transition.ID},
		)
		task.Task = requireResponse(t, http.StatusOK, response, err)
	}
}

func (task *Task) path() string {
	return task.project.path + "/tasks/" + strconv.Itoa(task.ID)
}
