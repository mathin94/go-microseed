package server

import (
	"context"
	"errors"
	"net/http"

	"microseed/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type HTTP struct {
	Srv *http.Server
}

func NewHTTP(cfg *config.Config, r *gin.Engine, logger *zap.Logger) *HTTP {
	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: r,
	}
	logger.Info("http server created", zap.String("addr", cfg.HTTPAddr))
	return &HTTP{Srv: server}
}

func RegisterHooks(lc fx.Lifecycle, h *HTTP, cfg *config.Config, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := h.Srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("http server failed", zap.Error(err))
				}
			}()
			logger.Info("http server started", zap.String("addr", h.Srv.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(ctx, cfg.GracefulTimeout)
			defer cancel()
			logger.Info("shutting down http server", zap.Duration("timeout", cfg.GracefulTimeout))
			return h.Srv.Shutdown(shutdownCtx)
		},
	})
}
