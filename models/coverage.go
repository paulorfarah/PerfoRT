package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type Coverage struct {
	Model
	MeasurementID      uint
	Measurement        Measurement
	Type               string
	Group              string
	Package            string
	Class              string
	InstructionMissed  int
	InstructionCovered int
	BranchMissed       int
	BranchCovered      int
	LineMissed         int
	LineCovered        int
	ComplexityMissed   int
	ComplexityCovered  int
	MethodMissed       int
	MethodCovered      int
}

func (r *Coverage) TableName() string {
	return "coverages"
}

func CreateCoverage(db *gorm.DB, c *Coverage) (uint, error) {
	err := db.Create(c).Error
	if err != nil {
		return 0, err
	}
	return c.ID, nil
}
