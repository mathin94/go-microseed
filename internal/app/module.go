package app

import (
	"context"
	"go.uber.org/zap"
	"microseed/internal/cache"
	"microseed/internal/config"
	"microseed/internal/db"
	"microseed/internal/domain/health"
	"microseed/internal/domain/user"
	"microseed/internal/httpx"
	applog "microseed/internal/log"
	"microseed/internal/obs"
	"microseed/internal/server"

	"go.uber.org/fx"
)

func provideLogger(cfg *config.Config) (*zap.Logger, error) {
	lg, err := applog.New(applog.Options{
		Level:          cfg.LogLevel,
		ConsoleEnabled: cfg.LogConsole,
		FilePath:       cfg.LogFilePath,
		MaxSizeMB:      cfg.LogFileMaxSizeMB,
		MaxBackups:     cfg.LogFileMaxBackups,
		MaxAgeDays:     cfg.LogFileMaxAgeDays,
		Compress:       cfg.LogFileCompress,
		StacktraceAt:   cfg.LogStackAt,
	})
	if err != nil {
		return nil, err
	}
	return lg.With(
		zap.String("service", cfg.AppName),
		zap.String("env", cfg.OTelEnv),
	), nil
}

func loggerHook(lc fx.Lifecycle, lg *zap.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			_ = lg.Sync()
			return nil
		},
	})
}

var Module = fx.Options(
	// Infra
	fx.Provide(
		config.New,
		provideLogger,
		obs.New,
		db.NewGorm,
		cache.NewRedis,
		httpx.NewRouter,
		server.NewHTTP,
	),
	fx.Invoke(
		server.RegisterHooks,
		cache.RegisterHooks,
		db.RegisterHooks,
		loggerHook,
	),

	// Routes auto-register
	httpx.RoutesModule,

	// Feature modules (bounded contexts)
	health.Module,
	user.Module,
)
