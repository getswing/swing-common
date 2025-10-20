package repository

import (
	"context"
	"errors"
	"fmt"

	"getswing.app/player-service/internal/app/models"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(ctx context.Context, u *models.Player) error
	FindByEmail(ctx context.Context, email string) (*models.Player, error)
	FindByID(ctx context.Context, id string) (*models.Player, error)
	FindByPhone(ctx context.Context, phoneNumber, phoneContryCode string) (*models.Player, error)
}

type playerRepo struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepo{db: db}
}

func (r *playerRepo) Create(ctx context.Context, u *models.Player) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("user create: %w", err)
	}
	return nil
}

func (r *playerRepo) FindByEmail(ctx context.Context, email string) (*models.Player, error) {
	var u models.Player
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("user find by email: %w", err)
	}
	return &u, nil
}

func (r *playerRepo) FindByPhone(ctx context.Context, phoneNumber, phoneCountryCode string) (*models.Player, error) {
	var u models.Player
	if err := r.db.WithContext(ctx).Where("phone_number = ? and phone_country_code = ?", phoneNumber, phoneCountryCode).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("user find by phone: %w", err)
	}
	return &u, nil
}

func (r *playerRepo) FindByID(ctx context.Context, id string) (*models.Player, error) {
	var u models.Player
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
