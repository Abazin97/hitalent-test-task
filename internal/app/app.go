package app

import (
	"context"
	"errors"
	"hitalent-test-task/internal/config"
	v1 "hitalent-test-task/internal/http"
	"hitalent-test-task/internal/repository"
	"hitalent-test-task/internal/services"
	"hitalent-test-task/pkg"
	log "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	cfg := config.Load()

	log.Info("initializing postgres")

	db, err := pkg.NewGormPostgres(cfg.DBURL)
	if err != nil {
		log.Error("failed to connect postgres", "error: ", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error(
			"failed to get sql db",
			"error",
			err,
		)
		return
	}

	defer func() {
		err := sqlDB.Close()
		if err != nil {
			log.Error("failed to close db", "error: ", err)
		}
	}()

	log.Info("initializing repositories")

	departmentRepo := repository.NewDepartmentRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)

	log.Info("initializing services")

	departmentService := services.NewDepartmentService(
		departmentRepo,
		employeeRepo,
	)

	log.Info("initializing handlers")

	h := v1.New(departmentService)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("starting http server", "port", cfg.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server failed", "error", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	log.Info("shutting down server")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error(
			"server shutdown failed",
			"error",
			err,
		)
		return
	}

	log.Info("server stopped")
}
