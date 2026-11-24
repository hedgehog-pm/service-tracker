package listmodules

import "context"

type ModuleWithProgress struct {
	ID         int
	ProjectID  int
	Name       string
	TasksTotal int64
	TasksDone  int64
}

type Usecase struct {
	modules  moduleReader
	projects projectReader
}

func New(modules moduleReader, projects projectReader) *Usecase {
	return &Usecase{
		modules:  modules,
		projects: projects,
	}
}

func (u *Usecase) Do(
	ctx context.Context,
	projectID int,
	limit int,
	offset int,
) ([]ModuleWithProgress, error) {
	if _, err := u.projects.GetProject(ctx, projectID); err != nil {
		return nil, err
	}

	return u.modules.ListModules(ctx, projectID, limit, offset)
}
