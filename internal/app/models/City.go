package models

type City struct {
	ID         uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string   `json:"name" gorm:"type:varchar(255)"`
	ProvinceID uint     `json:"province_id" gorm:"type:integer;not null"`
	Province   Province `json:"province" gorm:"foreignKey:ProvinceID;references:ID"`
}
