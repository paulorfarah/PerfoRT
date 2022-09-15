package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type Version struct {
	Model
	MeasurementID uint
	Measurement   Measurement
	Version       string
}

func (r *Version) TableName() string {
	return "versions"
}

func CreateVersion(db *gorm.DB, v *Version) (uint, error) {
	err := db.Create(v).Error
	if err != nil {
		return 0, err
	}
	return v.ID, nil
}
