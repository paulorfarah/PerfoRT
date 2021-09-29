package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type Change struct {
	Model
	// FileID         uint //   `gorm:"not null"`
	// File           File
	ChangeHash     string `gorm:"not null"`
	FileFromID     uint   `gorm:"not null"`
	FileFrom       File
	FileToID       uint `gorm:"not null"`
	FileTo         File
	Action         string `gorm:"not null"`
	Patch          string `gorm:"type:text;not null"`
	RandoopMetrics []RandoopMetrics
}

func (r *Change) TableName() string {
	return "changes"
}

func CreateChange(db *gorm.DB, change *Change) (uint, error) {
	err := db.Create(change).Error
	if err != nil {
		return 0, err
	}
	return change.ID, nil
}

func FindChangeByHash(db *gorm.DB, hash string, fileFromID uint) (*Change, error) {
	var change Change
	res := db.Find(&change, &Change{FileFromID: fileFromID, ChangeHash: hash})
	return &change, res.Error
}
