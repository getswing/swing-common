package models

type District struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name" gorm:"type:varchar(255)"`
	CityID uint   `json:"city_id" gorm:"type:integer;not null"`
	City   City   `json:"city" gorm:"foreignKey:CityID;references:ID"`
}
