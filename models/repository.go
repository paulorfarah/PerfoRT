package models

import (
	"github.com/jinzhu/gorm"
)

type Repository struct {
	Model
	PlatformFK uint `gorm:"unique_index:idx_repository"`
	Name string `gorm:"not null;unique_index:idx_repository"`
	Description string
	IsPrivate bool `gorm:"not null" sql:"DEFAULT:false"`
}

func (r *Repository) TableName() string {
	return "repositories"
}

func CreateRepository(db *gorm.DB, repository *Repository) (uint, error) {
	err := db.Create(repository).Error
	if err != nil {
		return 0, err
	}
	return repository.ID, nil
}

func FindRepositoryByName(db *gorm.DB, name string) (*Repository, error) {
	var repository Repository
	res := db.Find(&repository, &Repository{Name: name})
	return &repository, res.Error
}
