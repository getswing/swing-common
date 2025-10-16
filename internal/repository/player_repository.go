package repository

import (
	"context"
	"errors"
	"fmt"

	"getswing.app/player-service/internal/models"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(ctx context.Context, u *models.Player) error
	FindByEmail(ctx context.Context, email string) (*models.Player, error)
	FindByID(ctx context.Context, id uint) (*models.Player, error)
	FindByPhone(ctx context.Context, phoneNumber, phoneContryCode string) (*models.Player, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) PlayerRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *models.Player) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		return fmt.Errorf("user create: %w", err)
	}
	return nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*models.Player, error) {
	var u models.Player
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("user find by email: %w", err)
	}
	return &u, nil
}

func (r *userRepo) FindByPhone(ctx context.Context, phoneNumber, phoneCountryCode string) (*models.Player, error) {
	var u models.Player
	if err := r.db.WithContext(ctx).Where("phone_number = ? and phone_country_code", phoneNumber, phoneCountryCode).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("user find by phone: %w", err)
	}
	return &u, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uint) (*models.Player, error) {
	var u models.Player
	if err := r.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
