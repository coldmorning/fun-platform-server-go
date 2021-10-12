package categorypostresql

import (
	"fmt"
	"log"
	
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"github.com/coldmorning/fun-platform/config"
	"github.com/coldmorning/fun-platform/model"
)

var db *gorm.DB

func init() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		config.GetString("DB_POSTGRESQL.HOST"),
		config.GetString("DB_POSTGRESQL.PORT"),
		config.GetString("DB_POSTGRESQL.USER"),
		config.GetString("DB_POSTGRESQL.PASSWORD"),
		config.GetString("DB_POSTGRESQL.DBNAME"),
		config.GetString("DB_POSTGRESQL.SSLMODEL"),
		config.GetString("DB_POSTGRESQL.TimeZone"),
	)

	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Println("connection to postgres failed:", err)
	} else {
		log.Println("Database connected ...")
	}

}

func CreateCategory(category *model.CreateCategoryRequest) (err error) {

	//get user from database.

	err = db.Create(&category).Error
	if err != nil {
		log.Println("CreateCategory:", err)
		return  err
	}
	return  nil
}

