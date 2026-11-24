package projectworklogs

import "context"

type worklogReader interface {
	ListProjectWorklogs(ctx context.Context, query Query) (Page[Entry], error)
	ListProjectWorklogsByTask(ctx context.Context, query Query) (Page[TaskTotal], error)
}
