package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var db *gorm.DB

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id, omitempty"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at" sql:"DEFAULT:CURRENT_TIMESTAMP"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func init() {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")

	msql := mysql.Config{}
	log.Println(msql)

	strConn := username + ":" + password + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	aux, err := gorm.Open("mysql", strConn)

	db = aux
	if err != nil {
		fmt.Println(err)
	}

	db.LogMode(false)

	db.Debug().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&Account{},
		&Platform{},
		&Repository{},
		&Commit{},
		&File{},
		&Change{},
		// &Issue{},
		// &RandoopMetrics{},
		&Measurement{},
		// &Maven{},
		// &Randoop{},
		&TestCase{},
		&MeasurementResources{},
		// &TestResources{},
		&Coverage{},
		&CPUTimes{},
		&DiskIOCounters{},
		&NetIOCounters{},
	)

	db.Model(&Repository{}).AddForeignKey("platform_id", "platforms(id)", "RESTRICT", "RESTRICT")
	db.Model(&Commit{}).AddForeignKey("repository_id", "repositories(id)", "RESTRICT", "RESTRICT")
	db.Model(&Commit{}).AddForeignKey("author", "accounts(id)", "RESTRICT", "RESTRICT")
	db.Model(&Commit{}).AddForeignKey("committer", "accounts(id)", "RESTRICT", "RESTRICT")
	db.Model(&File{}).AddForeignKey("commit_id", "commits(id)", "RESTRICT", "RESTRICT")
	db.Model(&Change{}).AddForeignKey("file_from_id", "files(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Issue{}).AddForeignKey("repository_id", "repositories(id)", "RESTRICT", "RESTRICT")
	//	db.Model(&Issue{}).AddForeignKey("author", "accounts(id)", "RESTRICT", "RESTRICT")
	//	db.Model(&Issue{}).AddForeignKey("editor", "accounts(id)", "RESTRICT", "RESTRICT")
	// db.Model(&RandoopMetrics{}).AddForeignKey("change_id", "changes(id)", "RESTRICT", "RESTRICT")
	db.Model(&Measurement{}).AddForeignKey("repository_id", "repositories(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Maven{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Maven{}).AddForeignKey("commit_id", "commits(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Randoop{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Randoop{}).AddForeignKey("commit_id", "commits(id)", "RESTRICT", "RESTRICT")
	db.Model(&TestCase{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	db.Model(&TestCase{}).AddForeignKey("commit_id", "commits(id)", "RESTRICT", "RESTRICT")
	db.Model(&TestCase{}).AddForeignKey("file_id", "files(id)", "RESTRICT", "RESTRICT")
	db.Model(&TestCase{}).AddForeignKey("test_suite_id", "files(id)", "RESTRICT", "RESTRICT")

	db.Model(&MeasurementResources{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	// db.Model(&TestResources{}).AddForeignKey("test_id", "tests(id)", "RESTRICT", "RESTRICT")
	db.Model(&Coverage{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	db.Model(&CPUTimes{}).AddForeignKey("measurement_resources_id", "measurementresources(id)", "RESTRICT", "RESTRICT")
	db.Model(&DiskIOCounters{}).AddForeignKey("measurement_resources_id", "measurementresources(id)", "RESTRICT", "RESTRICT")
	db.Model(&NetIOCounters{}).AddForeignKey("measurement_resources_id", "measurementresources(id)", "RESTRICT", "RESTRICT")
	db.Model(&File{}).Related(&FileLine{})
}

func GetDB() *gorm.DB {
	return db
}
