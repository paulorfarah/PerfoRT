package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type TestCase struct {
	Model
	CommitID      uint
	Commit        Commit
	MeasurementID uint
	Measurement   Measurement
	FileID        uint
	TestSuiteID   uint
	Name          string  `gorm:"not null"`
	Type          string  `gorm:"not null"`
	Status        string  `gorm:"not null"`
	ClassName     string  `gorm:"not null"`
	Duration      float64 `gorm:"not null"`
	Error         string
	Message       string
	Properties    string
	SystemErr     string
	SystemOut     string
}

func (r *TestCase) TableName() string {
	return "testcases"
}

func CreateTestCase(db *gorm.DB, t *TestCase) (uint, error) {
	err := db.Create(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}
