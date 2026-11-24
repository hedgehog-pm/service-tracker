package api

import (
	"context"
	"slices"

	"github.com/samber/lo"

	"github.com/hedgehog-pm/service-tracker/internal/domain"
	generated "github.com/hedgehog-pm/service-tracker/internal/generated/api/boards"
	taskapi "github.com/hedgehog-pm/service-tracker/internal/generated/api/tasks"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/board"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createboardcolumn"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/deleteboardcolumn"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listboardcolumntasks"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listboards"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/moveboardtask"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/updateboardcolumn"
)

type BoardHandler struct {
	boards          *board.Usecase
	listBoards      *listboards.Usecase
	createColumn    *createboardcolumn.Usecase
	updateColumn    *updateboardcolumn.Usecase
	deleteColumn    *deleteboardcolumn.Usecase
	listColumnTasks *listboardcolumntasks.Usecase
	moveTask        *moveboardtask.Usecase
}

func NewBoardHandler(
	boards *board.Usecase,
	listBoards *listboards.Usecase,
	createColumn *createboardcolumn.Usecase,
	updateColumn *updateboardcolumn.Usecase,
	deleteColumn *deleteboardcolumn.Usecase,
	listColumnTasks *listboardcolumntasks.Usecase,
	moveTask *moveboardtask.Usecase,
) *BoardHandler {
	return &BoardHandler{
		boards:          boards,
		listBoards:      listBoards,
		createColumn:    createColumn,
		updateColumn:    updateColumn,
		deleteColumn:    deleteColumn,
		listColumnTasks: listColumnTasks,
		moveTask:        moveTask,
	}
}

func (h *BoardHandler) CreateBoard(
	ctx context.Context,
	request generated.CreateBoardRequestObject,
) (generated.CreateBoardResponseObject, error) {
	value, err := h.boards.Create(ctx, request.ProjectID, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return generated.CreateBoard201JSONResponse(boardToAPI(value)), nil
}

func (h *BoardHandler) GetBoard(
	ctx context.Context,
	request generated.GetBoardRequestObject,
) (generated.GetBoardResponseObject, error) {
	value, err := h.boards.Get(ctx, request.ProjectID, request.BoardID)
	if err != nil {
		return nil, err
	}

	return generated.GetBoard200JSONResponse(boardToAPI(value)), nil
}

func (h *BoardHandler) UpdateBoard(
	ctx context.Context,
	request generated.UpdateBoardRequestObject,
) (generated.UpdateBoardResponseObject, error) {
	value, err := h.boards.Update(ctx, request.ProjectID, request.BoardID, request.Body.Name)
	if err != nil {
		return nil, err
	}

	return generated.UpdateBoard200JSONResponse(boardToAPI(value)), nil
}

func (h *BoardHandler) DeleteBoard(
	ctx context.Context,
	request generated.DeleteBoardRequestObject,
) (generated.DeleteBoardResponseObject, error) {
	if err := h.boards.Delete(ctx, request.ProjectID, request.BoardID); err != nil {
		return nil, err
	}

	return generated.DeleteBoard204Response{}, nil
}

func (h *BoardHandler) ListBoards(
	ctx context.Context,
	request generated.ListBoardsRequestObject,
) (generated.ListBoardsResponseObject, error) {
	limit, offset := *request.Params.Limit, *request.Params.Offset

	values, err := h.listBoards.Do(ctx, request.ProjectID, limit, offset)
	if err != nil {
		return nil, err
	}

	items := lo.Map(values, func(value domain.Board, _ int) generated.BoardSummary {
		return generated.BoardSummary{
			ID:        value.ID,
			ProjectID: value.ProjectID,
			Name:      value.Name,
		}
	})

	return generated.ListBoards200JSONResponse(generated.BoardList{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}), nil
}

func (h *BoardHandler) CreateBoardColumn(
	ctx context.Context,
	request generated.CreateBoardColumnRequestObject,
) (generated.CreateBoardColumnResponseObject, error) {
	body := request.Body

	value, err := h.createColumn.Do(
		ctx,
		createboardcolumn.Params{
			ProjectID: request.ProjectID,
			BoardID:   request.BoardID,
			Name:      body.Name,
			Position:  body.Position,
			StatusIDs: body.StatusIds,
		},
	)
	if err != nil {
		return nil, err
	}

	return generated.CreateBoardColumn201JSONResponse(columnToAPI(request.BoardID, value)), nil
}

func (h *BoardHandler) UpdateBoardColumn(
	ctx context.Context,
	request generated.UpdateBoardColumnRequestObject,
) (generated.UpdateBoardColumnResponseObject, error) {
	body := request.Body

	value, err := h.updateColumn.Do(
		ctx,
		updateboardcolumn.Params{
			ProjectID: request.ProjectID,
			BoardID:   request.BoardID,
			ColumnID:  request.ColumnID,
			Name:      body.Name,
			Position:  body.Position,
			StatusIDs: body.StatusIds,
		},
	)
	if err != nil {
		return nil, err
	}

	return generated.UpdateBoardColumn200JSONResponse(columnToAPI(request.BoardID, value)), nil
}

func (h *BoardHandler) DeleteBoardColumn(
	ctx context.Context,
	request generated.DeleteBoardColumnRequestObject,
) (generated.DeleteBoardColumnResponseObject, error) {
	if err := h.deleteColumn.Do(
		ctx,
		request.ProjectID,
		request.BoardID,
		request.ColumnID,
	); err != nil {
		return nil, err
	}

	return generated.DeleteBoardColumn204Response{}, nil
}

func (h *BoardHandler) ListBoardColumnTasks(
	ctx context.Context,
	request generated.ListBoardColumnTasksRequestObject,
) (generated.ListBoardColumnTasksResponseObject, error) {
	query := listboardcolumntasks.Query{
		ProjectID: request.ProjectID,
		BoardID:   request.BoardID,
		ColumnID:  request.ColumnID,
		Limit:     *request.Params.Limit,
		Offset:    *request.Params.Offset,
		Filters: listboardcolumntasks.Filters{
			Query:      lo.FromPtr(request.Params.Q),
			StatusID:   request.Params.Status,
			ModuleID:   request.Params.Module,
			PriorityID: request.Params.Priority,
			TaskTypeID: request.Params.TaskType,
		},
	}

	tasks, err := h.listColumnTasks.Do(ctx, query)
	if err != nil {
		return nil, err
	}

	items := lo.Map(tasks, func(task domain.Task, _ int) taskapi.Task {
		return taskToAPI(task)
	})

	return generated.ListBoardColumnTasks200JSONResponse(generated.BoardColumnTaskList{
		Items:  items,
		Limit:  query.Limit,
		Offset: query.Offset,
	}), nil
}

func (h *BoardHandler) MoveBoardTask(
	ctx context.Context,
	request generated.MoveBoardTaskRequestObject,
) (generated.MoveBoardTaskResponseObject, error) {
	body := request.Body

	task, err := h.moveTask.Do(
		ctx,
		moveboardtask.Params{
			ProjectID:    request.ProjectID,
			BoardID:      request.BoardID,
			TaskID:       request.TaskID,
			ColumnID:     body.ColumnID,
			BeforeTaskID: body.BeforeTaskID,
			TransitionID: body.TransitionID,
		},
	)
	if err != nil {
		return nil, err
	}

	return generated.MoveBoardTask200JSONResponse(taskToAPI(task)), nil
}

func boardToAPI(value domain.Board) generated.Board {
	columns := lo.Map(value.Columns, func(column domain.BoardColumn, _ int) generated.BoardColumn {
		return columnToAPI(value.ID, column)
	})

	return generated.Board{
		ID:        value.ID,
		ProjectID: value.ProjectID,
		Name:      value.Name,
		Columns:   columns,
	}
}

func columnToAPI(boardID int, value domain.BoardColumn) generated.BoardColumn {
	statusIDs := slices.Clone(value.StatusIDs)
	slices.Sort(statusIDs)

	return generated.BoardColumn{
		ID:        value.ID,
		BoardID:   boardID,
		Name:      value.Name,
		Position:  value.Position,
		StatusIds: statusIDs,
	}
}
