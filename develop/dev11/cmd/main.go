package main

import (
	"L2-task11/internal/config"
	"L2-task11/internal/http-server/handlers/createEvent"
	"L2-task11/internal/http-server/handlers/deleteEvent"
	"L2-task11/internal/http-server/handlers/eventForDay"
	"L2-task11/internal/http-server/handlers/eventForMonth"
	"L2-task11/internal/http-server/handlers/eventForWeek"
	"L2-task11/internal/http-server/handlers/updateEvent"
	"L2-task11/internal/http-server/middleware"
	"L2-task11/internal/storage/sqlite"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// конфиг
	cfg := config.MustLoad()

	// logger
	myLogger, err := setupLogger(cfg.Env)
	if err != nil {
		log.Fatal("Could not setup logger: ", err)
	}

	myLogger.Info("Initializing server.")
	myLogger.Debug("debug messages are enabled")
	myLogger.Debug("Config:", cfg) //todo remove this line

	// storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		myLogger.Error("could not connect to the storage", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// http-server server
	mux := http.NewServeMux()

	mux.HandleFunc("/create_event", middleware.LoggerWrapper(createEvent.New(myLogger, storage), myLogger))
	mux.HandleFunc("/update_event", middleware.LoggerWrapper(updateEvent.New(myLogger, storage), myLogger))
	mux.HandleFunc("/delete_event", middleware.LoggerWrapper(deleteEvent.New(myLogger, storage), myLogger))
	mux.HandleFunc("/events_for_day", middleware.LoggerWrapper(eventForDay.New(myLogger, storage), myLogger))
	mux.HandleFunc("/events_for_week", middleware.LoggerWrapper(eventForWeek.New(myLogger, storage), myLogger))
	mux.HandleFunc("/events_for_month", middleware.LoggerWrapper(eventForMonth.New(myLogger, storage), myLogger))

	srv := &http.Server{
		Addr:         cfg.Host + ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	myLogger.Info("starting server: ", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		myLogger.Error("failed to start server")
	}

	myLogger.Info("server stopped")

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, fmt.Errorf("uknown env mode: %s", env)
	}
	return logger, nil
}
