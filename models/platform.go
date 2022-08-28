package models

import (
	"gorm.io/gorm"
)

type Platform struct {
	Model
	Repositories []Repository //`gorm:"foreignkey:PlatformFK"`
	Name         string       `gorm:"unique_index;not null;size:2048"`
}

func (p *Platform) TableName() string {
	return "platforms"
}

func CreatePlatform(db *gorm.DB, platform *Platform) (uint, error) {
	err := db.Create(platform).Error
	if err != nil {
		return 0, err
	}
	// fmt.Println("New platform added: " + platform.Name)
	return platform.ID, nil
}

func FindPlatformByName(db *gorm.DB, name string) (*Platform, error) {
	var platform Platform
	res := db.Where("name = ?", name).First(&platform)
	return &platform, res.Error
}
