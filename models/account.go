package models

import (
	"gorm.io/gorm"
)

type Account struct {
	Model
	Name             string   `gorm:"not null" json:"name"	validate:"required"`
	Email            string   `gorm:"not null" json:"email"	validate:"required,email"`
	CommitsAuthor    []Commit `gorm:"foreignKey:AuthorID"`
	CommitsCommitter []Commit `gorm:"foreignKey:CommitterID"`
}

func (u *Account) TableName() string {
	return "accounts"
}

func CreateAccount(db *gorm.DB, account *Account) (uint, error) {
	err := db.Create(account).Error
	if err != nil {
		return 0, err
	}
	return account.ID, nil
}

func FindAccountByEmail(db *gorm.DB, email string) (*Account, error) {
	var account Account
	res := db.Where("email = ?", email).First(&account)
	return &account, res.Error
}
