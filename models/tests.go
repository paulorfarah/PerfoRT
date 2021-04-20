package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Test struct {
	Model
	MeasurementID uint
	Measurement   Measurement
	Type          string `gorm:"not null"`
	CommitID      uint
	Commit        Commit
	ClassName     string
	TestsRun      int
	Failures      int
	Errors        int
	Skipped       int
	TimeElapsed   float64
}

func (r *Test) TableName() string {
	return "tests"
}

func CreateTest(db *gorm.DB, t *Test) (uint, error) {
	err := db.Create(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}
