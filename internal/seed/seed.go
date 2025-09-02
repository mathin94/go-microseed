package seed

import (
	"context"
	"time"

	"microseed/internal/domain/user"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedAll(ctx context.Context, db *gorm.DB, log *zap.Logger) error {
	// idempotent upsert by email
	users := []user.Entity{
		{ID: uuid.New(), Email: "user1@example.com", CreatedAt: time.Now()},
		{ID: uuid.New(), Email: "user2@example.com", CreatedAt: time.Now()},
	}
	for _, u := range users {
		if err := db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "email"}},
				DoNothing: true,
			}).Create(&u).Error; err != nil {
			return err
		}
	}
	log.Info("seed users applied", zap.Int("count", len(users)))
	return nil
}
