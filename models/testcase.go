package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type TestCase struct {
	Model
	// CommitID      uint
	// Commit        Commit
	// MeasurementID uint
	// Measurement   Measurement
	FileID       uint
	FileTargetID uint
	Name         string `gorm:"not null"`
	Type         string `gorm:"not null"`
	// Status        string        `gorm:"not null"`
	ClassName string `gorm:"not null"`
	// Duration      time.Duration `gorm:"not null"`
	// Error         string
	Message string
	// Properties    string
	// SystemErr string
	// SystemOut string
	// Methods []Method
}

func (r *TestCase) TableName() string {
	return "testcases"
}

func CreateTestCase(db *gorm.DB, t *TestCase) (uint, error) {
	// fmt.Println("creating testcase...")
	err := db.Create(t).Error
	if err != nil {
		return 0, err
	}
	return t.ID, nil
}
