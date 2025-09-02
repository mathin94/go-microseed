package db

import (
	"context"
	"time"

	"microseed/internal/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewGorm(cfg *config.Config, log *zap.Logger) (*gorm.DB, error) {
	gcfg := &gorm.Config{
		Logger: logger.New(
			zap.NewStdLog(log),
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logger.Warn,
				Colorful:      false,
			},
		),
	}
	db, err := gorm.Open(postgres.Open(cfg.DBDSN), gcfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpen)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	return db, nil
}

func RegisterHooks(lc fx.Lifecycle, gdb *gorm.DB, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			sqlDB, _ := gdb.DB()
			return sqlDB.Ping()
		},
		OnStop: func(_ context.Context) error {
			sqlDB, _ := gdb.DB()
			log.Info("closing db")
			return sqlDB.Close()
		},
	})
}

func Close(gdb *gorm.DB) error {
	sqlDB, err := gdb.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
