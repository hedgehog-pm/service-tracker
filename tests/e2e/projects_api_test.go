package e2e

import (
	"context"
	"net/http"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/projects"
)

func (p Project) Summary(ctx context.Context) (Response[projects.ProjectSummary], error) {
	return request[projects.ProjectSummary](ctx, p.client, http.MethodGet, p.path+"/summary", nil)
}
