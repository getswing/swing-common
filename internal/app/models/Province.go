package models

type Province struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string  `json:"name" gorm:"type:varchar(255)"`
	CountryID uint    `json:"country_id" gorm:"type:integer;not null"`
	Country   Country `json:"country" gorm:"foreignKey:CountryID;references:ID"`
}
