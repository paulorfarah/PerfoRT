package models

import (
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

type JSON []byte

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type Commit struct {
	Model
	RepositoryFK uint `gorm:"not null;unique_index:idx_commit"`
	CommitHash string `gorm:"unique;not null;unique_index:idx_commit"`
	TreeHash string `gorm:"not null"`
	ParentHashes JSON  `sql:"type:json" json:"parent_hashes,omitempty"` 
	Author  uint `gorm:"not null"` 
	AuthorDate time.Time `gorm:"not null"`
	Committer uint  `gorm:"not null"`
	CommitterDate time.Time `gorm:"not null"`
	Subject string `gorm:"not_null"`
	Branch string `gorm:"not_null"`
	Changes []Change
}

func (r *Commit) TableName() string {
	return "commits"
}

func CreateCommit(db *gorm.DB, commit *Commit) (uint, error) {
	err := db.Create(commit).Error
	if err != nil {
		return 0, err
	}
	return commit.ID, nil
}

func FindCommitByHash(db *gorm.DB, hash string) (*Commit, error) {
	var commit Commit
	res := db.Find(&commit, &Commit{CommitHash: hash})
	return &commit, res.Error
}
