package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type Change struct {
	Model
	CommitFK   uint   `gorm:"not null"`
	ChangeHash string `gorm:"not null"`
	FileFrom   string `gorm:"not null"`
	FileTo     string `gorm:"not null"`
	Action     string `gorm:"not null"`
	Patch      string `gorm:"type:text;not null"`
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

func FindChangeByHash(db *gorm.DB, hash string, commitID uint) (*Change, error) {
	var change Change
	res := db.Find(&change, &Change{CommitFK: commitID, ChangeHash: hash})
	return &change, res.Error
}
