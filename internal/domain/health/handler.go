package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	DB  *gorm.DB
	RDB *redis.Client
	Log *zap.Logger
}

func NewHandler(db *gorm.DB, rdb *redis.Client, log *zap.Logger) *Handler {
	return &Handler{DB: db, RDB: rdb, Log: log}
}

func (h *Handler) Register(r *gin.Engine) {
	r.GET("/healthz", h.liveness)
	r.GET("/readyz", h.readiness)
}

func (h *Handler) liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "type": "liveness"})
}

func (h *Handler) readiness(c *gin.Context) {
	// DB check
	sqlDB, _ := h.DB.DB()
	dbOK := (sqlDB.Ping() == nil)

	// Redis check
	ctx, cancel := context.WithTimeout(c.Request.Context(), 200*time.Millisecond)
	defer cancel()
	_, rerr := h.RDB.Ping(ctx).Result()
	redisOK := rerr == nil

	status := http.StatusOK
	if !dbOK || !redisOK {
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, gin.H{
		"db":    dbOK,
		"redis": redisOK,
		"ts":    time.Now().Format(time.RFC3339),
	})
}
