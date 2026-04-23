package app

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/pressly/goose/v3"
	logkit "github.com/wahrwelt-kit/go-logkit"
	_ "modernc.org/sqlite"

	"github.com/TakuyaYagam1/iustitia/config"
	"github.com/TakuyaYagam1/iustitia/internal/wire"
)

func Run(cfg *config.Config, l logkit.Logger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := openDB(ctx, cfg.DBPath)
	if err != nil {
		l.WithError(err).Error("failed to open database")
		return
	}
	defer func() { _ = db.Close() }()

	if err := applyMigrations(db, cfg.MigrationsDir, l); err != nil {
		l.WithError(err).Error("failed to apply migrations")
		return
	}

	app, err := wire.InitializeApp(cfg, l, db)
	if err != nil {
		l.WithError(err).Error("failed to initialize app")
		return
	}

	serverErr := make(chan error, 1)
	go func() {
		l.Info("http server listening", logkit.Fields{"addr": cfg.HTTPAddr})
		if err := app.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	select {
	case err := <-serverErr:
		l.WithError(err).Error("http server stopped unexpectedly")
	case <-ctx.Done():
		l.Info("shutdown signal received", nil)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := app.Server.Shutdown(shutdownCtx); err != nil {
		l.WithError(err).Error("server forced to shutdown")
		_ = app.Server.Close()
	}

	l.Info("http server stopped", nil)
}

func openDB(ctx context.Context, path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func applyMigrations(db *sql.DB, dir string, l logkit.Logger) error {
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}
	goose.SetLogger(goose.NopLogger())

	before, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}
	if err := goose.Up(db, dir); err != nil {
		return err
	}
	after, err := goose.GetDBVersion(db)
	if err != nil {
		return err
	}

	l.Info("migrations applied", logkit.Fields{
		"before": before,
		"after":  after,
		"dir":    dir,
	})
	return nil
}
