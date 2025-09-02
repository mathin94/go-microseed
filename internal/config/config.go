package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName         string
	HTTPAddr        string
	GracefulTimeout time.Duration

	// Postgres
	DBDSN             string
	DBMaxOpen         int
	DBMaxIdle         int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Logger
	LogLevel          string
	LogConsole        bool
	LogFilePath       string
	LogFileMaxSizeMB  int
	LogFileMaxBackups int
	LogFileMaxAgeDays int
	LogFileCompress   bool
	LogStackAt        string

	// OTel
	OTLPEndpoint string
	OTelService  string
	OTelEnv      string
}

func New() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(".env")
	v.SetConfigType("env")
	v.AutomaticEnv()

	// defaults
	v.SetDefault("APP_NAME", "microseed")
	v.SetDefault("HTTP_ADDR", ":8080")
	v.SetDefault("GRACEFUL_TIMEOUT", "10s")

	v.SetDefault("DB_DSN", "host=localhost user=postgres password=postgres dbname=microseed port=5432 sslmode=disable TimeZone=Asia/Jakarta")
	v.SetDefault("DB_MAX_OPEN", 30)
	v.SetDefault("DB_MAX_IDLE", 10)
	v.SetDefault("DB_CONN_MAX_LIFETIME", "60m")
	v.SetDefault("DB_CONN_MAX_IDLE_TIME", "10m")

	v.SetDefault("REDIS_ADDR", "localhost:6379")
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)

	// --- Logging defaults ---
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_CONSOLE", true)
	v.SetDefault("LOG_FILE_PATH", "") // kosong = nonaktif
	v.SetDefault("LOG_FILE_MAX_SIZE_MB", 50)
	v.SetDefault("LOG_FILE_MAX_BACKUPS", 5)
	v.SetDefault("LOG_FILE_MAX_AGE_DAYS", 30)
	v.SetDefault("LOG_FILE_COMPRESS", true)
	v.SetDefault("LOG_STACK_AT", "error")

	v.SetDefault("OTEL_EXPORTER_OTLP_ENDPOINT", "")
	v.SetDefault("OTEL_SERVICE_NAME", "microseed-api")
	v.SetDefault("OTEL_ENV", "dev")

	_ = v.ReadInConfig()

	timeout, _ := time.ParseDuration(v.GetString("GRACEFUL_TIMEOUT"))
	lifetime, _ := time.ParseDuration(v.GetString("DB_CONN_MAX_LIFETIME"))
	idleTime, _ := time.ParseDuration(v.GetString("DB_CONN_MAX_IDLE_TIME"))

	cfg := &Config{
		AppName:           v.GetString("APP_NAME"),
		HTTPAddr:          v.GetString("HTTP_ADDR"),
		GracefulTimeout:   defDur(timeout, 10*time.Second),
		DBDSN:             v.GetString("DB_DSN"),
		DBMaxOpen:         v.GetInt("DB_MAX_OPEN"),
		DBMaxIdle:         v.GetInt("DB_MAX_IDLE"),
		DBConnMaxLifetime: defDur(lifetime, 60*time.Minute),
		DBConnMaxIdleTime: defDur(idleTime, 10*time.Minute),
		RedisAddr:         v.GetString("REDIS_ADDR"),
		RedisPassword:     v.GetString("REDIS_PASSWORD"),
		RedisDB:           v.GetInt("REDIS_DB"),
		OTLPEndpoint:      v.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OTelService:       v.GetString("OTEL_SERVICE_NAME"),
		OTelEnv:           v.GetString("OTEL_ENV"),
		LogLevel:          v.GetString("LOG_LEVEL"),
		LogConsole:        v.GetBool("LOG_CONSOLE"),
		LogFilePath:       v.GetString("LOG_FILE_PATH"),
		LogFileMaxSizeMB:  v.GetInt("LOG_FILE_MAX_SIZE_MB"),
		LogFileMaxBackups: v.GetInt("LOG_FILE_MAX_BACKUPS"),
		LogFileMaxAgeDays: v.GetInt("LOG_FILE_MAX_AGE_DAYS"),
		LogFileCompress:   v.GetBool("LOG_FILE_COMPRESS"),
		LogStackAt:        v.GetString("LOG_STACK_AT"),
	}
	_ = os.Setenv("OTEL_SERVICE_NAME", cfg.OTelService)
	return cfg, nil
}

func defDur(d, def time.Duration) time.Duration {
	if d <= 0 {
		return def
	}
	return d
}
