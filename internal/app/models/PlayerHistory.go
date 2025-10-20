package models

type HistoryType string

const (
	EnumHistoryTypeTransaction   HistoryType = "TRANSACTION"
	EnumHistoryTypeDeleteAccount HistoryType = "DELETE_ACCOUNT"
)

type PlayerHistory struct {
	ID       string       `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type     HistoryType  `json:"type" gorm:"type:varchar(255)"`
	Platform EnumPlatform `json:"platform" gorm:"type:varchar(100);not null"`
	ActionBy string       `json:"action_by" gorm:"type:varchar(255)"`
	Rawdata  string       `json:"rawdata" gorm:"type:varchar(255)"`
	PlayerID string       `json:"player_id" gorm:"type:uuid"`
	Player   Player       `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
}
