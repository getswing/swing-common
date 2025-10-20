package models

import "time"

type NotificatinType string

const (
	EnumFlightArrangementUpdates   NotificatinType = "FLIGHT_ARRANGEMENT_UPDATES"
	EnumPartnerDashboardInvitation NotificatinType = "PARTNER_DASHBOARD_INVITATION"
	EnumReferralRemainder          NotificatinType = "REFERRAL_REMAINDER"
	EnumNotification               NotificatinType = "NOTIFICATION"
	EnumMpProductShared            NotificatinType = "MP_PRODUCT_SHARED"
	EnumMpShopShared               NotificatinType = "MP_SHOP_SHARED"
	EnumSwingPass                  NotificatinType = "SWING_PASS"
	EnumRemainderTournamentReview  NotificatinType = "REMAINDER_TOURNAMENT_REVIEW"
	EnumProductShared              NotificatinType = "PRODUCT_SHARED"
	EnumAnnouncement               NotificatinType = "ANNOUNCEMENT"
	EnumPromotion                  NotificatinType = "PROMOTION"
)

type Notification struct {
	ID                   string          `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title                string          `json:"title" gorm:"type:varchar(255)"`
	Type                 NotificatinType `json:"type" gorm:"type:varchar(200);not null"`
	Description          string          `json:"description" gorm:"type:varchar(255)"`
	Module               string          `json:"module" gorm:"type:varchar(255)"`
	TargetID             string          `json:"target_id" gorm:"type:varchar(255)"`
	Status               string          `json:"status" gorm:"type:varchar(255)"`
	Icon                 string          `json:"icon" gorm:"type:varchar(255)"`
	ImagePath            string          `json:"image_path" gorm:"type:varchar(255)"`
	CustomNotificationID string          `json:"custom_notification_id" gorm:"type:varchar(255)"`
	ProductID            string          `json:"product_id" gorm:"type:varchar(255)"`
	SentBy               string          `json:"sent_by" gorm:"type:varchar(255)"`
	ReferenceID          string          `json:"reference_id" gorm:"type:varchar(255)"`
	PlayerID             string          `json:"player_id" gorm:"type:uuid"`
	CreatedAt            *time.Time      `json:"created_at" gorm:"type:timestamptz"`
	UpdatedAt            *time.Time      `json:"updated_at" gorm:"type:timestamptz"`
	Player               Player          `json:"player" gorm:"foreignKey:PlayerID;references:ID;"`
}
