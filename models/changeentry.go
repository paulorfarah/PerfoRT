package models

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"

)

//https://git-scm.com/book/en/v2/Git-Basics-Viewing-the-Commit-History
type ChangeEntry struct {
	Model
	CommitFK uint `gorm:"not null;unique_index:idx_commit"`
}

