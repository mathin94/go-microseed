package httpx

import (
	"microseed/internal/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRouter(cfg *config.Config, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	for _, m := range Middlewares(logger) {
		r.Use(m)
	}
	_ = cfg // reserved for future toggles/flags
	return r
}
