package models

import (
	"time"

	"gorm.io/gorm"
)

type ReferralMember struct {
	ID              string         `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	PlayerID        string         `json:"player_id" gorm:"type:uuid"`
	ReferralGroupID string         `json:"referral_group_id" gorm:"type:uuid"`
	CreatedAt       *time.Time     `json:"created_at" gorm:"type:timestamptz"`
	UpdatedAt       *time.Time     `json:"updated_at" gorm:"type:timestamptz"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamptz;index"`
	Player          Player         `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
	ReferralGroup   ReferralGroup  `json:"referral_group" gorm:"foreignKey:ReferralGroupID;references:ID;"`
}
