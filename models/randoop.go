package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Randoop struct {
	Model
	MeasurementID uint //   `gorm:"not null"`
	Measurement   Measurement
	Type          byte `gorm:"not null"`
	CommitID      uint
	Commit        Commit
	ClassName     string `gorm:"size:2048"`
	TestsRun      int
	// Failures    int
	// Errors      int
	// Skipped    int
	TimeElapsed float64
}

func (r *Randoop) TableName() string {
	return "randoop"
}

func CreateRandoop(db *gorm.DB, results *Randoop) (uint, error) {
	err := db.Create(results).Error
	if err != nil {
		return 0, err
	}
	return results.ID, nil
}
