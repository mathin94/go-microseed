package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Entity struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (Entity) TableName() string { return "users" }

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Entity, error)
}

type serviceImpl struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &serviceImpl{DB: db}
}

func (s *serviceImpl) GetByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
	var u Entity
	if err := s.DB.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
