package models

type SubDistrict struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name       string    `json:"name" gorm:"type:varchar(255)"`
	DistrictID uint      `json:"district_id" gorm:"type:integer;not null"`
	District   District  `json:"district" gorm:"foreignKey:DistrictID;references:ID"`
	Zipcodes   []Zipcode `json:"zipcodes" gorm:"foreignKey:SubDistrictID;references:ID"`
}
