package models

import (
	"time"

	"gorm.io/gorm"
)

type EnumAddressType string

const (
	EnumTypePrimary EnumAddressType = "PRIMARY"
)

type PlayerAddress struct {
	ID            string          `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title         string          `json:"title" gorm:"type:varchar(255)"`
	Latitude      float64         `json:"latitude" gorm:"type:float64;not null"`
	Longitude     float64         `json:"longitude" gorm:"type:float64;not null"`
	Detail        string          `json:"detail" gorm:"type:varchar(255)"`
	Type          EnumAddressType `json:"type" gorm:"type:varchar(255)"`
	PhoneNumber   string          `json:"phone_number" gorm:"type:varchar(255)"`
	DialCode      string          `json:"dialcode" gorm:"type:varchar(5)"`
	PlayerID      string          `json:"player_id" gorm:"type:uuid"`
	SubdistrictID int32           `json:"subdistrict_id" gorm:"type:integer"`
	Zipcode       string          `json:"zipcode" gorm:"type:varchar(5)"`
	CreatedAt     *time.Time      `json:"created_at" gorm:"type:timestamptz"`
	UpdatedAt     *time.Time      `json:"updated_at" gorm:"type:timestamptz"`
	DeletedAt     gorm.DeletedAt  `json:"deleted_at" gorm:"type:timestamptz;index"`
	Player        Player          `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
	SubDistrict   SubDistrict     `json:"sub_district" gorm:"foreignKey:SubdistrictIDs;references:ID;"`
}
