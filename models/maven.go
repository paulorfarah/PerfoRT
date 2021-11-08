package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type Maven struct {
	Model
	MeasurementID uint //   `gorm:"not null"`
	Measurement   Measurement
	Type          byte `gorm:"not null"`
	CommitID      uint
	COmmit        Commit
	ClassName     string
	TestsRun      int
	Failures      int
	Errors        int
	Skipped       int
	TimeElapsed   float64
}

func (r *Maven) TableName() string {
	return "maven"
}

func CreateMaven(db *gorm.DB, results *Maven) (uint, error) {
	err := db.Create(results).Error
	if err != nil {
		return 0, err
	}
	return results.ID, nil
}
