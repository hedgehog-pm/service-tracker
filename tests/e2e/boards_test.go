package e2e

import (
	"testing"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
)

func TestBoardTaskFlow(t *testing.T) {
	api := StartPostgres(t).
		Migrate(t).
		StartService(t).
		API()

	project := api.CreateProject(t, projects.CreateProject{
		Key:  "FLOW",
		Name: "Board task flow",
	})

	columns := project.ExpectOnlyBoard(t).ExpectLayout(t, "Основная",
		ColumnSpec{
			Name:     "Открытые",
			Statuses: []string{"Open"},
		},
		ColumnSpec{
			Name:     "В работе",
			Statuses: []string{"In Progress"},
		},
		ColumnSpec{
			Name:     "Готово",
			Statuses: []string{"Done"},
		},
	)
	open := columns["Открытые"]
	inProgress := columns["В работе"]
	done := columns["Готово"]

	taskType := project.TaskType(t, "Task")
	a := project.CreateTask(t, tasks.CreateTask{
		Title:      "A",
		TaskTypeID: taskType.ID,
	})
	b := project.CreateTask(t, tasks.CreateTask{
		Title:      "B",
		TaskTypeID: taskType.ID,
	})
	open.ExpectTasks(t, a, b)

	open.Move(t, b, MoveOptions{
		Before: &a,
	})
	open.ExpectTasks(t, b, a)

	inProgress.Move(t, a, MoveOptions{
		Transition: "Start progress",
	})
	open.ExpectTasks(t, b)
	inProgress.ExpectTasks(t, a)

	inProgress.Move(t, b, MoveOptions{
		Transition: "Start progress",
	})
	open.ExpectTasks(t)
	inProgress.ExpectTasks(t, a, b)

	done.Move(t, a, MoveOptions{
		Transition: "Complete",
	})
	inProgress.ExpectTasks(t, b)
	done.ExpectTasks(t, a)
}

func TestProjectGetsBoardForWorkflowScheme(t *testing.T) {
	api := StartPostgres(t).
		Migrate(t).
		ApplyFixture(t, "testdata/review_workflow.sql").
		StartService(t).
		API()

	project := api.CreateProject(t, projects.CreateProject{
		Key:  "REVIEW",
		Name: "Project with review",
	})
	project.ExpectOnlyBoard(t).ExpectLayout(t, "Доска с проверкой",
		ColumnSpec{
			Name:     "К работе",
			Statuses: []string{"Open"},
		},
		ColumnSpec{
			Name:     "В работе",
			Statuses: []string{"In Progress"},
		},
		ColumnSpec{
			Name:     "На проверке",
			Statuses: []string{"Review"},
		},
		ColumnSpec{
			Name:     "Готово",
			Statuses: []string{"Done"},
		},
	)
}
