package postgresql

import (
	"fmt"
	"fun-platform-server/config"
	"fun-platform-server/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func FetchUser(user domain.User) (returnUser *domain.User, err error) {

	//get user from database.

	err = db.Table("user").Where("email= ? AND password =? AND delete_time is null ", user.Email, user.Password).First(&returnUser).Error
	if err != nil {
		log.Println("FetchUser:", err)
		return nil, err
	}
	return returnUser, nil
}
