package domain

import (
	"errors"
	"testing"
	"time"
)

func TestTaskTransition(t *testing.T) {
	before := time.Unix(1, 0).UTC()
	now := time.Unix(2, 0).UTC()
	original := Task{WorkflowID: 10, StatusID: 1, UpdatedAt: before}

	for _, test := range []struct {
		name         string
		workflowID   int
		fromStatusID int
		wantErr      error
	}{
		{name: "available", workflowID: 10, fromStatusID: 1},
		{name: "another workflow", workflowID: 20, fromStatusID: 1, wantErr: ErrTransitionUnavailable},
		{name: "another status", workflowID: 10, fromStatusID: 3, wantErr: ErrTransitionUnavailable},
	} {
		t.Run(test.name, func(t *testing.T) {
			task := original
			transition := Transition{
				WorkflowID: test.workflowID, FromStatusID: test.fromStatusID,
				ToStatus: Status{ID: 2},
			}

			err := task.Transition(transition, now)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Transition() error = %v, want %v", err, test.wantErr)
			}

			want := original
			if test.wantErr == nil {
				want.StatusID = 2
				want.UpdatedAt = now
			}

			if task != want {
				t.Fatalf("Transition() task = %+v, want %+v", task, want)
			}
		})
	}
}
