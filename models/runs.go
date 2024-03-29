package models

import (
	"time"

	"gorm.io/gorm"
)

type Run struct {
	Model
	MeasurementID uint
	Measurement   Measurement
	TestCaseID    uint
	TestCase      TestCase
	Type          string
	Number        int
	TestCaseTime  time.Duration
}

func (r *Run) TableName() string {
	return "runs"
}

func CreateRun(db *gorm.DB, mr *Run) (uint, error) {
	err := db.Create(mr).Error
	if err != nil {
		return 0, err

	}
	return mr.ID, nil
}

func SaveRun(db *gorm.DB, mr *Run) error {
	return db.Save(mr).Error

}
