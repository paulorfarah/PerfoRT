package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type Measurement struct {
	Model
	RepositoryID uint //   `gorm:"not null"`
	// Repository   Repository
	Runs            int           `gorm:"default:1"`
	TestcaseTimeout time.Duration `gorm:"default:3600"`
	MonitoringTime  time.Duration `gorm:"default:1"`
	// Maven        []Maven
	// Randoop      []Randoop
}

func (r *Measurement) TableName() string {
	return "measurements"
}

func CreateMeasurement(db *gorm.DB, measurement *Measurement) (uint, error) {
	err := db.Create(measurement).Error
	if err != nil {
		return 0, err
	}
	return measurement.ID, nil
}
