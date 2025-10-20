package models

type Zipcode struct {
	ID            uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code          string   `json:"code" gorm:"type:varchar(5)"`
	SubDistrictID uint     `json:"sub_district_id" gorm:"type:integer;not null"`
	SubDistrict   District `json:"sub_district" gorm:"foreignKey:SubDistrictID;references:ID"`
}
