package models

import (
	"time"

	"gorm.io/gorm"
)


type ReferralGroup struct {
	ID            string          `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PlayerID      string          `json:"player_id" gorm:"type:uuid"`
	CreatedAt     *time.Time      `json:"created_at" gorm:"type:timestamptz"`
	UpdatedAt     *time.Time      `json:"updated_at" gorm:"type:timestamptz"`
	DeletedAt     gorm.DeletedAt  `json:"deleted_at" gorm:"type:timestamptz;index"`
	Player        Player          `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
}
