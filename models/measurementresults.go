package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type MeasurementResults struct {
	Model
	MeasurementID uint //   `gorm:"not null"`
	Measurement   Measurement
	Type          byte `gorm:"not null"`
	//Commit
	ClassName         string
	TestsRunBefore    int
	FailuresBefore    int
	ErrorsBefore      int
	SkippedBefore     int
	TimeElapsedBefore float64
	TestsRunAfter     int
	FailuresAfter     int
	ErrorsAfter       int
	SkippedAfter      int
	TimeElapsedAfter  float64
}

func (r *MeasurementResults) TableName() string {
	return "measurementresults"
}

func CreateMeasurementResults(db *gorm.DB, measurementResults *MeasurementResults) (uint, error) {
	err := db.Create(measurementResults).Error
	if err != nil {
		return 0, err
	}
	return measurementResults.ID, nil
}
