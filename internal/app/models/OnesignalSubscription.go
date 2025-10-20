package models

import "time"

type OsType string

const (
	EnumTypeCustomer OsType = "CUSTOMER"
	EnumTypePartner  OsType = "PARTNER"
)

type OnesignalSubscription struct {
	ID        string     `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OsID      string     `json:"osid" gorm:"type:varchar(200)"`
	Type      OsType     `json:"type" gorm:"type:varchar(100);not null"`
	PlayerID  string     `json:"player_id" gorm:"type:uuid"`
	CreatedAt *time.Time `json:"created_at" gorm:"type:timestamptz"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:timestamptz"`
	Player    Player     `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
}
