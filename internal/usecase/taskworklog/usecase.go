package taskworklog

import (
	"context"
	"time"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
)

// Usecase ведёт журнал выполненной работы по задачам.
// Каждое добавление трудозатрат создаёт отдельную запись в журнале.
type Usecase struct {
	worklogs taskWorklogRepository
}

func New(worklogs taskWorklogRepository) *Usecase {
	return &Usecase{worklogs: worklogs}
}

// Get возвращает запись, только если она относится к указанной задаче проекта.
func (u *Usecase) Get(ctx context.Context, projectID, taskID, worklogID int) (domain.TaskWorklog, error) {
	return u.worklogs.GetTaskWorklog(ctx, projectID, taskID, worklogID)
}

// List возвращает страницу записей и сумму времени по всей задаче.
// Сумма учитывает и записи за пределами страницы.
func (u *Usecase) List(ctx context.Context, projectID, taskID, limit, offset int) ([]domain.TaskWorklog, int64, error) {
	return u.worklogs.ListTaskWorklogs(ctx, projectID, taskID, limit, offset)
}

// Create добавляет запись с серверным временем создания.
// Дата выполненной работы берётся из fields.StartedAt и может отличаться от времени создания.
func (u *Usecase) Create(ctx context.Context, projectID, taskID int, fields domain.TaskWorklogFields) (domain.TaskWorklog, error) {
	now := time.Now().UTC()

	return u.worklogs.InsertTaskWorklog(ctx, projectID, domain.TaskWorklog{
		TaskID:            taskID,
		TaskWorklogFields: fields,
		CreatedAt:         now,
		UpdatedAt:         now,
	})
}

// Update исправляет дату работы, длительность и комментарий существующей записи.
// Задача и время создания сохраняются, а время изменения задаёт сервер.
func (u *Usecase) Update(ctx context.Context, projectID, taskID, worklogID int, fields domain.TaskWorklogFields) (domain.TaskWorklog, error) {
	return u.worklogs.UpdateTaskWorklog(ctx, projectID, domain.TaskWorklog{
		ID:                worklogID,
		TaskID:            taskID,
		TaskWorklogFields: fields,
		UpdatedAt:         time.Now().UTC(),
	})
}

// Delete удаляет указанную запись, сохраняя остальные записи журнала.
func (u *Usecase) Delete(ctx context.Context, projectID, taskID, worklogID int) error {
	return u.worklogs.DeleteTaskWorklog(ctx, projectID, taskID, worklogID)
}
