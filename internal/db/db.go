package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectToDatabase(connectionString string) (db *gorm.DB, err error) {

	db, err = gorm.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	return
}
