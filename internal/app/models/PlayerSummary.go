package models

type PlayerSummary struct {
	ID                string `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	BookingCount      int32  `json:"boking_count" gorm:"type:integer"`
	NotificationReads int32  `json:"notification_reads" gorm:"type:integer"`
	NotificationCount int32  `json:"notification_count" gorm:"type:integer"`
	PlayerID          string `json:"player_id" gorm:"type:uuid"`
	Player            Player `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
}
