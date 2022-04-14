package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type TestCase struct {
	Model
	FileID       uint
	FileTargetID uint
	Name         string `gorm:"not null"`
	Type         string `gorm:"not null"`
	ClassName    string `gorm:"not null"`
	Message      string
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
