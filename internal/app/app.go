package app

import (
	"context"
	"fmt"
	"log/slog"
	"main/internal/bootstrap"
	"main/internal/config"
	"main/internal/repositories/peoplerepository/peoplesqlx"
	"main/internal/repositories/taskrepository/tasksqlx"
	"main/internal/services/timetrackerservice"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
)

func Run(cfg config.Config) error {
	db, err := bootstrap.InitSqlxDB(cfg)
	if err != nil {
		slog.Error("Failed to initialize database connection.", "Error message", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = goose.UpContext(ctx, db.DB, cfg.DBMigrationsDir)
	if err != nil {
		slog.Error("Failed to apply database migrations.", "Error message", err)
		os.Exit(1)
	}

	autoService := timetrackerservice.New(peoplesqlx.New(db), tasksqlx.New(db))

	server := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", cfg.ServerHost, cfg.ServerPort),
		Handler: autoService.GetHandler(),
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			slog.Warn("", "msg", err)
		}
	}()

	gracefullyShutdown(ctx, cancel, server)
	return nil
}

func gracefullyShutdown(ctx context.Context, cancel context.CancelFunc, server *http.Server) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Warn("", "msg", err)
	}
	cancel()
}
