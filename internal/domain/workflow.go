package domain

type Transition struct {
	ID           int
	WorkflowID   int
	FromStatusID int
	Name         string
	ToStatus     Status
}

type InitialStatus struct {
	WorkflowID int
	Status     Status
}

type TaskType struct {
	ID       int
	Name     string
	Position int
}

type TaskTypeWithInitialStatus struct {
	TaskType
	InitialStatus Status
}
