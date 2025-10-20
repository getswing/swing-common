package models

import (
	"time"

	"gorm.io/gorm"
)

type Gender string

const (
	EnumGenderMale        Gender = "MALE"
	EnumGenderFemale      Gender = "FEMALE"
	EnumGenderNoSpecified Gender = "NO_SPECIFIED"
)

type Player struct {
	ID                   string          `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	FirstName            string          `json:"first_name" gorm:"type:varchar(255)"`
	LastName             string          `json:"last_name" gorm:"type:varchar(255)"`
	Dob                  *time.Time      `json:"dob" gorm:"type:timestamptz"`
	Password             string          `json:"password" gorm:"type:varchar(255)"`
	Email                string          `json:"email" gorm:"type:varchar(255)"`
	DialCode             string          `json:"dialcode" gorm:"type:varchar(5)"`
	PhoneNumber          string          `json:"phone_number" gorm:"type:varchar(255)"`
	Nationality          string          `json:"nationality" gorm:"type:varchar(255)"`
	Experienced          string          `json:"experienced" gorm:"type:varchar(255)"`
	IsBetaTester         *bool           `json:"is_beta_tester"`
	IsOnboarding         *bool           `json:"is_onboarding"`
	Gender               Gender          `json:"gender" gorm:"type:varchar(255);default:NO_SPECIFIED"`
	ProfilePath          string          `json:"profile_path" gorm:"type:varchar(255)"`
	IsFirstBuy           *bool           `json:"is_first_buy"`
	Username             string          `json:"username" gorm:"type:varchar(255)"`
	IsDefaultUsername    *bool           `json:"is_default_username"`
	IsFirstBuyGolfCourse *bool           `json:"is_firstbuy_golf_course"`
	CreatedAt            *time.Time      `json:"created_at" gorm:"type:timestamptz"`
	UpdatedAt            *time.Time      `json:"updated_at" gorm:"type:timestamptz"`
	DeletedAt            gorm.DeletedAt  `json:"deleted_at" gorm:"type:timestamptz;index"`
	Histories            []PlayerHistory `json:"histories" gorm:"foreignKey:PlayerID;references:ID"`
}
