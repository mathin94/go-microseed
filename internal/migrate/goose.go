package migrate

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"microseed/internal/config"

	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx stdlib driver for goose
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func openDB(cfg *config.Config) (*sql.DB, error) {
	// pgx stdlib menerima DSN key=val atau URL
	return sql.Open("pgx", cfg.DBDSN)
}

func prepare() error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(goose.NopLogger())
	return nil
}

func Up(ctx context.Context, cfg *config.Config, log *zap.Logger) error {
	if err := prepare(); err != nil {
		return err
	}
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	log.Info("migrations up applied")
	return nil
}

func Down(ctx context.Context, cfg *config.Config, log *zap.Logger, steps int) error {
	if err := prepare(); err != nil {
		return err
	}
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if steps <= 0 {
		steps = 1
	}
	if err := goose.DownToContext(ctx, db, "migrations", int64(steps)); err != nil {
		return fmt.Errorf("goose down: %w", err)
	}
	log.Info("migrations down applied", zap.Int("steps", steps))
	return nil
}

func Reset(ctx context.Context, cfg *config.Config, log *zap.Logger) error {
	if err := prepare(); err != nil {
		return err
	}
	db, err := openDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.ResetContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("goose reset: %w", err)
	}
	log.Info("migrations reset to version 0")
	return nil
}
