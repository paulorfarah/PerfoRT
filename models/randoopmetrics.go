package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type RandoopMetrics struct {
	Model
	ChangeFK   uint `gorm:"not null"`
	NMEBefore  string
	EMEBefore  string
	AETNBefore string
	AETEBefore string
	AMUBefore  string
	NMEAfter   string
	EMEAfter   string
	AETNAfter  string
	AETEAfter  string
	AMUAfter   string
	NMEDiff    string
	EMEDiff    string
	AETNDiff   string
	AETEDiff   string
	AMUDiff    string
	NMEPerc    string
	EMEPerc    string
	AETNPerc   string
	AETEPerc   string
	AMUPerc    string
}

func (p *RandoopMetrics) TableName() string {
	return "randoopmetrics"
}

func CreateRandoopMetrics(db *gorm.DB, rm *RandoopMetrics) (uint, error) {
	err := db.Create(rm).Error
	if err != nil {
		return 0, err
	}
	fmt.Println("New randoop metrics added: ")
	return rm.ID, nil
}

// func FindPlatformByName(db *gorm.DB, name string) (*Platform, error) {
// 	var platform Platform
// 	res := db.Find(&platform, &Platform{Name: name})
// 	return &platform, res.Error
// }
