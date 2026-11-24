package e2e

import (
	"context"
	"net/http"
	"testing"

	"github.com/hedgehog-pm/service-tracker/internal/generated/api/modules"
)

func (p Project) AddModule(t testing.TB, name string) modules.Module {
	t.Helper()

	response, err := request[modules.Module](
		t.Context(),
		p.client,
		http.MethodPost,
		p.path+"/modules",
		modules.CreateModule{Name: name},
	)

	return requireResponse(t, http.StatusCreated, response, err)
}

func (p Project) Modules(ctx context.Context) (Response[modules.ModuleList], error) {
	return request[modules.ModuleList](ctx, p.client, http.MethodGet, p.path+"/modules", nil)
}
