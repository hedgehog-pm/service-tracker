package e2e

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/netip"
	"os"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	dockernetwork "github.com/moby/moby/api/types/network"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

type Postgres struct {
	container *postgres.PostgresContainer
	dsn       string
	network   *testcontainers.DockerNetwork
}

type MigratedPostgres struct {
	db    Postgres
	image string
}

type Service struct {
	url    string
	client *http.Client
}

func StartPostgres(t testing.TB) Postgres {
	t.Helper()

	if os.Getenv("TRACKER_E2E") != "1" {
		t.Skip("run e2e tests with make test-e2e")
	}

	ctx := t.Context()
	nw, err := network.New(ctx)
	testcontainers.CleanupNetwork(t, nw)
	require.NoError(t, err, "create Docker network")

	const name = "tracker"

	db, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase(name),
		postgres.WithUsername(name),
		postgres.WithPassword(name),
		postgres.BasicWaitStrategies(),
		bindLoopback("5432/tcp"),
		network.WithNetwork([]string{"postgres"}, nw),
	)
	testcontainers.CleanupContainer(t, db)
	require.NoError(t, err, "start PostgreSQL")

	return Postgres{
		container: db,
		dsn:       fmt.Sprintf("postgres://%s:%s@postgres:5432/%s?sslmode=disable", name, name, name),
		network:   nw,
	}
}

func (db Postgres) Migrate(t testing.TB) MigratedPostgres {
	t.Helper()

	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Minute)
	defer cancel()

	migrate, err := testcontainers.Run(ctx, "",
		testcontainers.WithDockerfile(testcontainers.FromDockerfile{
			Context:        "../..",
			BuildLogWriter: t.Output(),
		}),
		testcontainers.WithCmd("/app/migrate"),
		testcontainers.WithEnv(map[string]string{
			"DATABASE_URL": db.dsn,
		}),
		network.WithNetwork(nil, db.network),
		testcontainers.WithWaitStrategy(wait.ForExit()),
	)
	cleanupContainer(t, migrate)
	require.NoError(t, err, "build image and run migrations")

	state, err := migrate.State(ctx)
	require.NoError(t, err, "inspect migration container")
	require.Zero(t, state.ExitCode, "migration container exit code")

	return MigratedPostgres{
		db:    db,
		image: migrate.Image,
	}
}

func (m MigratedPostgres) ApplyFixture(t testing.TB, path string) MigratedPostgres {
	t.Helper()

	sql, err := os.ReadFile(path)
	require.NoError(t, err, "read SQL fixture %s", path)

	args := []string{
		"psql",
		"--no-psqlrc",
		"--set", "ON_ERROR_STOP=1",
		"--dbname", m.db.dsn,
		"--command", string(sql),
	}

	code, output, err := m.db.container.Exec(t.Context(), args, exec.Multiplexed())
	require.NoError(t, err, "apply SQL fixture %s", path)

	if code != 0 {
		data, err := io.ReadAll(output)
		require.NoError(t, err, "read output for fixture %s after psql exited with code %d", path, code)
		require.FailNow(t, "psql failed", "SQL fixture %s exited with code %d\n%s", path, code, data)
	}

	return m
}

func (m MigratedPostgres) StartService(t testing.TB) Service {
	t.Helper()

	ctx := t.Context()
	service, err := testcontainers.Run(ctx, m.image,
		testcontainers.WithEnv(map[string]string{
			"DATABASE_URL": m.db.dsn,
		}),
		testcontainers.WithExposedPorts("8080/tcp"),
		bindLoopback("8080/tcp"),
		network.WithNetwork(nil, m.db.network),
		testcontainers.WithWaitStrategy(wait.ForHTTP("/projects").WithPort("8080/tcp")),
	)
	cleanupContainer(t, service)
	require.NoError(t, err, "start service")

	url, err := service.PortEndpoint(ctx, "8080/tcp", "http")
	require.NoError(t, err, "get service URL")

	return Service{
		url:    url,
		client: newHTTPClient(t),
	}
}

func (s Service) API() *Client {
	return &Client{
		url:  s.url,
		http: s.client,
	}
}

func bindLoopback(port string) testcontainers.CustomizeRequestOption {
	return testcontainers.WithHostConfigModifier(func(config *container.HostConfig) {
		config.PortBindings = dockernetwork.PortMap{
			dockernetwork.MustParsePort(port): {
				{
					HostIP: netip.MustParseAddr("127.0.0.1"),
				},
			},
		}
	})
}

func cleanupContainer(t testing.TB, container *testcontainers.DockerContainer) {
	t.Helper()
	testcontainers.CleanupContainer(t, container)

	if container == nil {
		return
	}

	t.Cleanup(func() {
		if !t.Failed() {
			return
		}

		ctx, cancel := context.WithTimeout(context.WithoutCancel(t.Context()), 5*time.Second)
		defer cancel()

		logs, err := container.Logs(ctx)
		if err != nil {
			t.Logf("read container logs: %v", err)
			return
		}
		defer func() { _ = logs.Close() }()

		if _, err := io.Copy(t.Output(), logs); err != nil {
			t.Logf("copy container logs: %v", err)
		}
	})
}
