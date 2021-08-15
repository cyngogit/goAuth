package database

import (
	"path"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"example.com/go/auth/models"
)

var DBCon *gorm.DB

func Connect() {
	user := "root:Admin@"
	db := "goauthdb"
	dsn := path.Join(user, db)
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	showError(err)
	DBCon = connection

	connection.AutoMigrate(&models.User{})
}

func showError(err error) {
	if err != nil {
		panic("could not connect to the database")
	}
}
