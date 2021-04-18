package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MavenResources struct {
	Model
	MeasurementID uint //   `gorm:"not null"`
	Measurement   Measurement
	Cpu           float64
	Mem           float32
	ReadCount     uint64 `json:"readCount"`
	WriteCount    uint64 `json:"writeCount"`
	ReadBytes     uint64 `json:"readBytes"`
	WriteBytes    uint64 `json:"writeBytes"`
}

func (r *MavenResources) TableName() string {
	return "mavenresources"
}

func CreateMavenResources(db *gorm.DB, mr *MavenResources) (uint, error) {
	err := db.Create(mr).Error
	if err != nil {
		return 0, err
	}
	return mr.ID, nil
}
