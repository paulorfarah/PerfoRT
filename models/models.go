package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var err error

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
	db, err = gorm.Open(mysql.Open(strConn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Error),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		fmt.Println(err)
	}

	db.Debug().Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&Account{},
		&Platform{},
		&Repository{},
		&Commit{},
		&File{},
		&Method{},
		&Change{},
		// &Issue{},
		// &RandoopMetrics{},
		&Measurement{},
		&TestCase{},
		&Run{},
		&Resource{},
		&Coverage{},
		&CPUTimes{},
		&DiskIOCounters{},
		&NetIOCounters{},
		&Jvm{},
	)

	// db.Model(&Repository{}).AddForeignKey("platform_id", "platforms(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Commit{}).AddForeignKey("repository_id", "repositories(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Commit{}).AddForeignKey("author", "accounts(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Commit{}).AddForeignKey("committer", "accounts(id)", "RESTRICT", "RESTRICT")
	// db.Model(&File{}).AddForeignKey("commit_id", "commits(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Change{}).AddForeignKey("file_from_id", "files(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Measurement{}).AddForeignKey("repository_id", "repositories(id)", "RESTRICT", "RESTRICT")
	// db.Model(&TestCase{}).AddForeignKey("file_id", "files(id)", "RESTRICT", "RESTRICT")

	// db.Model(&Run{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Run{}).AddForeignKey("test_case_id", "testcases(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Resource{}).AddForeignKey("run_id", "runs(id)", "RESTRICT", "RESTRICT")

	// db.Model(&Coverage{}).AddForeignKey("measurement_id", "measurements(id)", "RESTRICT", "RESTRICT")
	// db.Model(&CPUTimes{}).AddForeignKey("resource_id", "resources(id)", "RESTRICT", "RESTRICT")
	// db.Model(&DiskIOCounters{}).AddForeignKey("resource_id", "resources(id)", "RESTRICT", "RESTRICT")
	// db.Model(&NetIOCounters{}).AddForeignKey("resource_id", "resources(id)", "RESTRICT", "RESTRICT")
	// db.Model(&File{}).Related(&FileLine{})

	//set max connections
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("ERROR trying to open database to set max connections")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

}

func GetDB() *gorm.DB {
	return db
}
