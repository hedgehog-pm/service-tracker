package domain

import (
	"slices"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/samber/lo"
)

type Board struct {
	ID           int
	ProjectID    int
	Name         string
	SchemeID     int
	SharedScheme bool
	Columns      []BoardColumn
}

type BoardColumn struct {
	ID        int
	Name      string
	Position  int
	StatusIDs []int
}

func (b Board) Validate() error {
	if err := validation.Validate(b.Name, validation.By(notBlank)); err != nil {
		return NewError(CodeInvalidInput, err.Error())
	}

	return nil
}

func (c BoardColumn) ContainsStatus(statusID int) bool {
	return slices.Contains(c.StatusIDs, statusID)
}

func (b Board) Column(id int) (BoardColumn, error) {
	index := b.columnIndex(id)
	if index == -1 {
		return BoardColumn{}, ErrBoardColumnNotFound
	}

	return b.Columns[index], nil
}

func (b *Board) AddColumn(column BoardColumn, available []Status) error {
	if b.SharedScheme {
		return ErrBoardSchemeShared
	}

	if column.Position < 0 || column.Position > len(b.Columns) {
		return NewError(CodeInvalidInput, "column position is out of range")
	}

	if err := b.validateColumn(column, available); err != nil {
		return err
	}

	b.Columns = slices.Insert(b.Columns, column.Position, column)
	b.reindexColumns()

	return nil
}

func (b *Board) UpdateColumn(column BoardColumn, available []Status) error {
	if b.SharedScheme {
		return ErrBoardSchemeShared
	}

	index := b.columnIndex(column.ID)
	if index == -1 {
		return ErrBoardColumnNotFound
	}

	if column.Position < 0 || column.Position >= len(b.Columns) {
		return NewError(CodeInvalidInput, "column position is out of range")
	}

	if err := b.validateColumn(column, available); err != nil {
		return err
	}

	b.Columns = slices.Delete(b.Columns, index, index+1)
	b.Columns = slices.Insert(b.Columns, column.Position, column)
	b.reindexColumns()

	return nil
}

func (b *Board) RemoveColumn(id int) error {
	if b.SharedScheme {
		return ErrBoardSchemeShared
	}

	index := b.columnIndex(id)
	if index == -1 {
		return ErrBoardColumnNotFound
	}

	b.Columns = slices.Delete(b.Columns, index, index+1)
	b.reindexColumns()

	return nil
}

func (b Board) columnIndex(id int) int {
	return slices.IndexFunc(b.Columns, func(column BoardColumn) bool {
		return column.ID == id
	})
}

func (b *Board) reindexColumns() {
	for i := range b.Columns {
		b.Columns[i].Position = i
	}
}

func (b Board) validateColumn(column BoardColumn, available []Status) error {
	if err := validation.Validate(column.Name, validation.By(notBlank)); err != nil {
		return NewError(CodeInvalidInput, err.Error())
	}

	if len(column.StatusIDs) == 0 {
		return NewError(CodeInvalidInput, "column must contain at least one status")
	}

	allowed := lo.Associate(available, func(status Status) (int, bool) {
		return status.ID, true
	})

	seen := make(map[int]bool, len(column.StatusIDs))
	for _, id := range column.StatusIDs {
		if !allowed[id] {
			return ErrStatusUnavailable
		}

		if seen[id] {
			return NewError(CodeInvalidInput, "column contains duplicate statuses")
		}

		seen[id] = true
	}

	return b.checkStatusMapping(column)
}

func (b Board) checkStatusMapping(column BoardColumn) error {
	for _, other := range b.Columns {
		if column.ID != 0 && other.ID == column.ID {
			continue
		}

		if slices.ContainsFunc(column.StatusIDs, other.ContainsStatus) {
			return ErrStatusAlreadyMapped
		}
	}

	return nil
}
