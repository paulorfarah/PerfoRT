package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type TestCase struct {
	Model
	FileID       uint
	FileTargetID uint
	Name         string `gorm:"not null;size:2048"`
	Type         string `gorm:"not null"`
	ClassName    string `gorm:"not null;size:2048"`
	Message      string
	Error        bool
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

func SetTestCaseError(db *gorm.DB, t *TestCase) {
	db.Model(&t).Update("Error", true)
}
