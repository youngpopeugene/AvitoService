package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
)

var db *gorm.DB
var databaseErr error

func initiator() {
	if !db.Migrator().HasTable("reserves") {
		db.Table("reserves").AutoMigrate(&Reserve{})
		db.Table("reserves")
		fmt.Println("Table 'reserves' was created!")
	}
	if !db.Migrator().HasTable("users") {
		db.Table("users").AutoMigrate(&User{})
		fmt.Println("Table 'users' was created!")
	}

}
func connector() {
	host := os.Getenv("pgHost")
	port := os.Getenv("pgPort")
	portDigit, _ := strconv.ParseInt(port, 10, 64)
	user := os.Getenv("pgUser")
	password := os.Getenv("pgPassword")
	dbname := os.Getenv("pgDbName")
	db, databaseErr = gorm.Open(postgres.Open(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, portDigit, user, password, dbname)), &gorm.Config{})
	if databaseErr != nil {
		panic(databaseErr)
	}
	fmt.Println("Successfully connected to database!")
}
