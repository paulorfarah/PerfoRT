package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type File struct {
	Model
	CommitID    uint
	Commit      Commit
	Hash        string `gorm:"not null;idx_commit"`
	Name        string `gorm:"not null"`
	Size        int64
	Contents    string `gorm:"type:text"`
	IsBinary    bool   `gorm:"not null;default: false"`
	Lines       []FileLine
	IsMalformed bool `gorm:"not null;default: false"`
	HasChanged  bool `gorm:"not null;default: false"`
}

type FileLine struct {
	gorm.Model
	Line string
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
	res := db.Find(&file, &File{Name: filename, HasChanged: false})
	return &file, res.Error
}

func FindFileByNameAndCommit(db *gorm.DB, filename string, commitID uint) (*File, error) {
	var file File
	res := db.Find(&file, &File{Name: filename, CommitID: commitID})
	return &file, res.Error
}
