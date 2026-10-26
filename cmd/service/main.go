package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hedgehog-pm/service-tracker/internal/api"
	modulerepository "github.com/hedgehog-pm/service-tracker/internal/repository/module"
	priorityrepository "github.com/hedgehog-pm/service-tracker/internal/repository/priority"
	projectrepository "github.com/hedgehog-pm/service-tracker/internal/repository/project"
	statusrepository "github.com/hedgehog-pm/service-tracker/internal/repository/status"
	taskrepository "github.com/hedgehog-pm/service-tracker/internal/repository/task"
	workflowrepository "github.com/hedgehog-pm/service-tracker/internal/repository/workflow"
	"github.com/hedgehog-pm/service-tracker/internal/storage/postgres"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createproject"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/createtask"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/issuetype"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listmodules"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listtasks"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/listtasktransitions"
	moduleusecase "github.com/hedgehog-pm/service-tracker/internal/usecase/module"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/priority"
	projectusecase "github.com/hedgehog-pm/service-tracker/internal/usecase/project"
	statususecase "github.com/hedgehog-pm/service-tracker/internal/usecase/status"
	taskusecase "github.com/hedgehog-pm/service-tracker/internal/usecase/task"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/transitiontask"
	"github.com/hedgehog-pm/service-tracker/internal/usecase/updatetask"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		slog.ErrorContext(ctx, "service failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return errors.New("DATABASE_URL is not set")
	}

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	tasks := taskrepository.NewPostgres(db)
	projects := projectrepository.NewPostgres(db)
	transactor := postgres.NewTransactor(db)
	workflows := workflowrepository.NewPostgres(db)
	priorities := priorityrepository.NewPostgres(db)
	taskHandler := api.NewTaskHandler(
		taskusecase.New(tasks), createtask.New(tasks, projects, workflows, priorities, transactor),
		listtasks.New(tasks), updatetask.New(tasks, transactor),
	)

	statuses := statusrepository.NewPostgres(db)
	modules := modulerepository.NewPostgres(db)
	projectHandler := api.NewProjectHandler(projectusecase.New(projects), createproject.New(projects, workflows, priorities))
	statusHandler := api.NewStatusHandler(statususecase.New(statuses, projects))
	moduleHandler := api.NewModuleHandler(moduleusecase.New(modules), listmodules.New(modules, projects))

	transitionHandler := api.NewTransitionHandler(listtasktransitions.New(tasks, workflows), transitiontask.New(tasks, workflows, transactor))
	issueTypeHandler := api.NewIssueTypeHandler(issuetype.New(workflows, projects))
	priorityHandler := api.NewPriorityHandler(priority.New(priorities, projects))

	server, err := api.NewServer(projectHandler, statusHandler, moduleHandler, taskHandler, transitionHandler, issueTypeHandler, priorityHandler)
	if err != nil {
		return err
	}

	address := os.Getenv("HTTP_ADDR")
	if address == "" {
		address = ":8080"
	}

	httpServer := &http.Server{Addr: address, Handler: server, ReadHeaderTimeout: 5 * time.Second}

	return serve(ctx, httpServer)
}

func serve(ctx context.Context, server *http.Server) error {
	listener, err := (&net.ListenConfig{}).Listen(ctx, "tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("listen http: %w", err)
	}
	defer func() { _ = server.Close() }()

	serveErr := make(chan error, 1)
	go func() { serveErr <- server.Serve(listener) }()

	slog.InfoContext(ctx, "http server started", "address", listener.Addr().String())

	select {
	case <-ctx.Done():
	case err := <-serveErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve http: %w", err)
		}

		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := <-serveErr; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}

	return nil
}
