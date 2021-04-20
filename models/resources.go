package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MeasurementResources struct {
	Model
	MeasurementID uint //   `gorm:"not null"`
	Measurement   Measurement
	Type          string
	Cpu           float64
	Mem           float32
	ReadCount     uint64 `json:"readCount"`
	WriteCount    uint64 `json:"writeCount"`
	ReadBytes     uint64 `json:"readBytes"`
	WriteBytes    uint64 `json:"writeBytes"`
}

func (r *MeasurementResources) TableName() string {
	return "measurementresources"
}

func CreateMeasurementResources(db *gorm.DB, mr *MeasurementResources) (uint, error) {
	err := db.Create(mr).Error
	if err != nil {
		return 0, err
	}
	return mr.ID, nil
}

type TestResources struct {
	Model
	TestID     uint
	Test       Test
	Type       string
	Cpu        float64
	Mem        float32
	ReadCount  uint64 `json:"readCount"`
	WriteCount uint64 `json:"writeCount"`
	ReadBytes  uint64 `json:"readBytes"`
	WriteBytes uint64 `json:"writeBytes"`
}

func (r *TestResources) TableName() string {
	return "testresources"
}

func CreateTestResources(db *gorm.DB, tr *TestResources) (uint, error) {
	err := db.Create(tr).Error
	if err != nil {
		return 0, err
	}
	return tr.ID, nil
}
