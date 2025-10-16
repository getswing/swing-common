package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Player struct {
	ID                    uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName             string         `gorm:"type:varchar(255)"`
	LastName              string         `gorm:"type:varchar(255)"`
	Dob                   *time.Time     `gorm:"type:timestamptz"`
	ReferralCode          string         `gorm:"type:varchar(255)"`
	Password              string         `gorm:"type:varchar(255)"`
	Email                 string         `gorm:"type:varchar(255)"`
	PhoneCountryCode      string         `gorm:"type:varchar(255)"`
	PhoneNumber           string         `gorm:"type:varchar(255)"`
	Nationality           string         `gorm:"type:varchar(255)"`
	CreatedAt             *time.Time     `gorm:"type:timestamptz"`
	UpdatedAt             *time.Time     `gorm:"type:timestamptz"`
	DeletedAt             gorm.DeletedAt `gorm:"type:timestamptz;index"`
	Experienced           string         `gorm:"type:varchar(255)"`
	ImageFileID           *uuid.UUID     `gorm:"type:uuid"`
	IsBetaTester          *bool
	IsOnboarding          *bool
	EmailVerificationCode string `gorm:"type:varchar(255)"`
	Gender                string `gorm:"type:varchar(255);default:NO_SPECIFIED"`
	EmailTemporary        string `gorm:"type:varchar(255)"`
	DeleteReason          string `gorm:"type:varchar(255)"`
	IsFirstBuy            *bool
	Username              string     `gorm:"type:varchar(255)"`
	IsDefaultUsername     *bool      `gorm:"default:true"`
	IsFirstBuyGolfCourse  *bool      `gorm:"default:true"`
	OneSignalID           *uuid.UUID `gorm:"type:uuid"`
	SubscriptionID        string     `gorm:"type:text"`
	PartnerSubscriptions  string     `gorm:"type:text"`
}
