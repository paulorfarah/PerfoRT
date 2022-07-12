package models

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

type JSON []byte

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type Commit struct {
	Model
	RepositoryID       uint `gorm:"index;not null"`
	Repository         Repository
	CommitHash         string `gorm:"unique;not null;unique_index:idx_commit"`
	PreviousCommitHash string
	TreeHash           string    `gorm:"not null"`
	ParentHashes       JSON      `sql:"type:json" json:"parent_hashes,omitempty"`
	AuthorID           uint      `gorm:"not null"`
	AuthorDate         time.Time `gorm:"not null"`
	CommitterID        uint      `gorm:"not null"`
	CommitterDate      time.Time `gorm:"index;not null"`
	Subject            string    `gorm:"type:text;not_null"`
	Branch             string    `gorm:"not_null"`
	// Changes            []Change
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
	// res := db.Find(&commit, &Commit{CommitHash: hash})
	res := db.Where("commit_hash = ?", hash).First(&commit)
	return &commit, res.Error
}

type repoCommit struct {
	Name          string
	CommitHash    string
	CommitterDate time.Time
}

// func GetRepoCommitList(db *gorm.DB) {
// 	rows, err := db.Table("repositories").Select("repositories.name, commits.commit_hash, commits.committer_date").Joins("left join commits on commits.repository_fk = repositories.id").Rows()
// 	defer Rows.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	name := ""
// 	commit_hash := ""
// 	var committer_date time.Time
// 	for rows.Next() {
// 		rows.Scan(&name, &commit_hash, &committer_date)
// 		fmt.Printf("%s, %s, %s\n", name, commit_hash, committer_date)
// 	}
// }
