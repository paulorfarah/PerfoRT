package models

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

var db *gorm.DB

type Model struct {
	ID	uint	`gorm:"primary_key" json:"id, omitempty"`
	CreatedAt	time.Time	`gorm:"not null" json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt	time.Time	`gorm:"not null" json:"updated_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	DeletedAt	*time.Time	`sql:"index" json:"deleted_at,omitempty"`
}

func init() {
	e := godotenv.Load()
	if e !=  nil {
		log.Fatal(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName   := os.Getenv("db_name")
	dbHost	 := os.Getenv("db_host")
	dbPort	 := os.Getenv("db_port")

	msql := mysql.Config{}
	log.Println(msql)

	strConn := username + ":" + password + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	aux, err := gorm.Open("mysql", strConn)

	db = aux 
	if err != nil {
		fmt.Println(err)
	}

	db.LogMode(true)


	db.Debug().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&Account{},
		&Platform{},
		&Repository{},
		&Commit{},
	)
	
	db.Model(&Repository{}).AddForeignKey("platform_fk", "platform(id)", "RESTRICT", "RESTRICT")
	db.Model(&Commit{}).AddForeignKey("repository_fk", "repositories(id)", "RESTRICT", "RESTRICT")
	db.Model(&Commit{}).AddForeignKey("author", "accounts(id)", "RESTRICT", "RESTRICT")
	db.Model(&Commit{}).AddForeignKey("committer", "accounts(id)", "RESTRICT", "RESTRICT")

}

func GetDB() *gorm.DB {
	return db
}
