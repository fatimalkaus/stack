package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatimalkaus/stack/internal/config"
	"github.com/fatimalkaus/stack/internal/db"
	"github.com/fatimalkaus/stack/internal/handler"
	"github.com/fatimalkaus/stack/internal/stack"
	"github.com/gorilla/mux"
)

// Start starts application.
func Start(version string) {
	cfgPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		die("failed to load config: " + err.Error())
	}

	db, err := db.InitPostgres(cfg.Postgres)
	if err != nil {
		die(err.Error())
	}

	stack, err := stack.NewPostgresStack(db)
	if err != nil {
		die(err.Error())
	}

	router := mux.NewRouter()
	handler.InitRoutes(router, stack)

	server := http.Server{
		Addr:              net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		slog.Info("server started", "version", version, "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
		close(done)
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "err", err)
	}
	<-done
	slog.Info("server stopped")
}

func die(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
