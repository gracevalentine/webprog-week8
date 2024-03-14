package Controller

import (
	"database/sql"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func connectForGorm() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp?parseTime=true&loc=Asia%2FJakarta")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
		return nil, err
	}
	db.LogMode(true)
	return db, err
}
