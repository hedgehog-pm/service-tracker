package e2e

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
	"github.com/hedgehog-pm/service-tracker/internal/generated/api/taskworklogs"
)

func TestProjectStatistics(t *testing.T) {
	api := StartPostgres(t).
		Migrate(t).
		StartService(t).
		API()

	t.Run("summary", func(t *testing.T) {
		checkProjectSummary(t, api)
	})
	t.Run("modules", func(t *testing.T) {
		checkModuleProgress(t, api)
	})
	t.Run("worklogs", func(t *testing.T) {
		checkProjectWorklogs(t, api)
	})
}

func checkProjectSummary(t *testing.T, api *Client) {
	project := api.CreateProject(t, projects.CreateProject{
		Key:  "STAT",
		Name: "Project summary",
	})

	response, err := project.Summary(t.Context())
	empty := requireResponse(t, http.StatusOK, response, err)
	assert.Equal(t, projects.ProjectSummary{
		UpcomingTasks: []projects.UpcomingTask{},
	}, empty)

	today := time.Now().UTC().Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)

	project.AddTask(t, TaskSpec{
		Title:   "Overdue",
		DueDate: today.AddDate(0, 0, -2),
	})
	outsideWeek := project.AddTask(t, TaskSpec{
		Title:   "Outside the week",
		DueDate: today.AddDate(0, 0, 7),
	})
	lastDay := project.AddTask(t, TaskSpec{
		Title:   "Last day of the week",
		DueDate: today.AddDate(0, 0, 6),
	})
	completed := project.AddTask(t, TaskSpec{
		Title:           "Already completed",
		DueDate:         today.AddDate(0, 0, 2),
		EstimateSeconds: new(900),
	})
	completed.Transition(t, "Start progress", "Complete")

	soon := project.AddTask(t, TaskSpec{
		Title:           "Due tomorrow",
		DueDate:         tomorrow,
		EstimateSeconds: new(600),
	})
	sameDate := project.AddTask(t, TaskSpec{
		Title:   "Also due tomorrow",
		DueDate: tomorrow,
	})
	later := project.AddTask(t, TaskSpec{
		Title:   "Later",
		DueDate: today.AddDate(0, 0, 8),
	})
	project.AddTask(t, TaskSpec{
		Title:   "Beyond the preview",
		DueDate: today.AddDate(0, 0, 9),
	})

	startedAt := time.Date(2026, 1, 2, 9, 0, 0, 0, time.UTC)
	soon.AddWorklog(t, startedAt, 6*time.Minute)
	soon.AddWorklog(t, startedAt.Add(time.Hour), 3*time.Minute)

	response, err = project.Summary(t.Context())
	summary := requireResponse(t, http.StatusOK, response, err)

	assert.Equal(t, int64(8), summary.TasksTotal)
	assert.Equal(t, int64(7), summary.TasksOpen)
	assert.Equal(t, int64(3), summary.TasksDueSoon)
	assert.Equal(t, int64(1), summary.TasksOverdue)
	assert.Equal(t, int64(540), summary.TotalTimeSpentSeconds)
	assert.Equal(t, new(int64(1500)), summary.TotalEstimateSeconds)

	upcomingIDs := lo.Map(summary.UpcomingTasks, func(task projects.UpcomingTask, _ int) int {
		return task.ID
	})
	assert.Equal(t, []int{
		soon.ID,
		sameDate.ID,
		lastDay.ID,
		outsideWeek.ID,
		later.ID,
	}, upcomingIDs)
}

func checkModuleProgress(t *testing.T, api *Client) {
	project := api.CreateProject(t, projects.CreateProject{
		Key:  "MOD",
		Name: "Module progress",
	})
	active := project.AddModule(t, "Active module")
	empty := project.AddModule(t, "Empty module")

	project.AddTask(t, TaskSpec{Title: "Open"}).AssignToModule(t, active)

	completed := project.AddTask(t, TaskSpec{Title: "Completed"})
	completed.AssignToModule(t, active)
	completed.Transition(t, "Start progress", "Complete")
	project.AddTask(t, TaskSpec{Title: "Without a module"})

	response, err := project.Modules(t.Context())
	list := requireResponse(t, http.StatusOK, response, err)

	assert.Equal(t, []modules.ModuleListItem{
		{
			ID:         active.ID,
			Name:       active.Name,
			ProjectID:  active.ProjectID,
			TasksTotal: 2,
			TasksDone:  1,
		},
		{
			ID:        empty.ID,
			Name:      empty.Name,
			ProjectID: empty.ProjectID,
		},
	}, list.Items)
}

func checkProjectWorklogs(t *testing.T, api *Client) {
	ctx := t.Context()
	project := api.CreateProject(t, projects.CreateProject{
		Key:  "TIME",
		Name: "Project worklogs",
	})
	importTask := project.AddTask(t, TaskSpec{
		Title:           "Import 100%",
		EstimateSeconds: new(1200),
	})
	retryTask := project.AddTask(t, TaskSpec{Title: "Import 1000"})

	from := time.Date(2026, 1, 2, 10, 0, 0, 0, time.UTC)
	to := from.Add(time.Hour)
	importTask.AddWorklog(t, from.Add(-time.Hour), 30*time.Second)
	importLog := importTask.AddWorklog(t, from, time.Minute)
	retryLog := retryTask.AddWorklog(t, from, 90*time.Second)
	retryTask.AddWorklog(t, to, 2*time.Minute)

	other := api.CreateProject(t, projects.CreateProject{
		Key:  "OTHER",
		Name: "Other project",
	})
	other.AddTask(t, TaskSpec{Title: importTask.Title}).AddWorklog(t, from, 2*time.Hour)

	journalResponse, err := project.Worklogs(ctx, WorklogQuery{
		Limit:  2,
		Offset: 1,
	})
	journal := requireResponse(t, http.StatusOK, journalResponse, err)

	journalIDs := lo.Map(journal.Items, func(entry taskworklogs.ProjectWorklog, _ int) int {
		return entry.ID
	})
	assert.Equal(t, []int{retryLog.ID, importLog.ID}, journalIDs)
	assert.Equal(t, 2, journal.Limit)
	assert.Equal(t, 1, journal.Offset)
	assert.Equal(t, int64(300), journal.TotalTimeSpentSeconds)
	assert.Equal(t, new(int64(1200)), journal.TotalEstimateSeconds)

	rankedResponse, err := project.WorklogsByTask(ctx, WorklogQuery{Limit: 1})
	ranked := requireResponse(t, http.StatusOK, rankedResponse, err)

	assert.Equal(t, int64(300), ranked.TotalTimeSpentSeconds)
	require.Len(t, ranked.Items, 1)
	assert.Equal(t, retryTask.ID, ranked.Items[0].TaskID)
	assert.Equal(t, int64(210), ranked.Items[0].TimeSpentSeconds)
	assert.Equal(t, int64(2), ranked.Items[0].WorklogsCount)

	groupedResponse, err := project.WorklogsByTask(ctx, WorklogQuery{
		From: from,
		To:   to,
	})
	grouped := requireResponse(t, http.StatusOK, groupedResponse, err)

	assert.Equal(t, []taskworklogs.ProjectWorklogByTask{
		{
			TaskID:           retryTask.ID,
			TaskKey:          retryTask.Key,
			TaskTitle:        retryTask.Title,
			TimeSpentSeconds: 90,
			WorklogsCount:    1,
		},
		{
			TaskID:           importTask.ID,
			TaskKey:          importTask.Key,
			TaskTitle:        importTask.Title,
			TimeSpentSeconds: 60,
			WorklogsCount:    1,
		},
	}, grouped.Items)
	assert.Equal(t, int64(150), grouped.TotalTimeSpentSeconds)
	assert.Equal(t, journal.TotalEstimateSeconds, grouped.TotalEstimateSeconds)

	searchedResponse, err := project.Worklogs(ctx, WorklogQuery{
		From:   from,
		To:     to,
		Search: "IMPORT 100%",
	})
	searched := requireResponse(t, http.StatusOK, searchedResponse, err)

	assert.Equal(t, int64(60), searched.TotalTimeSpentSeconds)
	require.Len(t, searched.Items, 1)
	assert.Equal(t, importLog.ID, searched.Items[0].ID)
	assert.Equal(t, importTask.Key, searched.Items[0].TaskKey)

	exhaustedResponse, err := project.WorklogsByTask(ctx, WorklogQuery{
		From:   from,
		To:     to,
		Search: strings.ToLower(retryTask.Key),
		Offset: 1,
	})
	exhausted := requireResponse(t, http.StatusOK, exhaustedResponse, err)

	assert.Equal(t, []taskworklogs.ProjectWorklogByTask{}, exhausted.Items)
	assert.Equal(t, int64(90), exhausted.TotalTimeSpentSeconds)
	assert.Equal(t, journal.TotalEstimateSeconds, exhausted.TotalEstimateSeconds)
}
