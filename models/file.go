package models

import (
	"fmt"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

// https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type File struct {
	Model
	CommitID uint `gorm:"index"`
	Commit   Commit
	Hash     string `gorm:"not null"`
	Name     string `gorm:"index;not null;size:2048"`
	Size     int64
	// Contents string `gorm:"type:text"`
	IsBinary bool `gorm:"not null;default: false"`
	// Lines       []FileLine
	IsMalformed bool `gorm:"not null;default: false"`
	HasChanged  bool `gorm:"not null;default: false"`
	Methods     []Method
}

type FileLine struct {
	gorm.Model
	FileID uint
	File   File
	Line   string
}

func (r *File) TableName() string {
	return "files"
}

func CreateFile(db *gorm.DB, f *File) (uint, error) {
	err := db.Create(f).Error
	if err != nil {
		fmt.Printf("Error creating file %s: %s\n", f.Name, err.Error())
		return 0, err
	}
	return f.ID, nil
}

func FindFileByName(db *gorm.DB, filename string) (*File, error) {
	var file File
	// res := db.Find(&file, &File{Name: filename, HasChanged: false})
	res := db.Where("name like ?", "%"+filename)
	return &file, res.Error
}

func FindFileByEndsWithNameAndCommit(db *gorm.DB, filename string, commitID uint) (*File, error) {
	var file File
	// res := db.Find(&file, &File{Name: filename, CommitID: commitID})
	res := db.Where("name like ? and commit_id=?", "%"+filename, commitID).First(&file)
	return &file, res.Error
}
