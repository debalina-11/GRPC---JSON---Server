package database

import (
	"fmt"
	"grpc-json-server/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func Connect() {
	dsn := "debalina:1234@tcp(127.0.0.1:3306)/grpc?charset=utf8mb4&parseTime=True&loc=Local" //without docker- 127.0.0.1:3306 //with docker - host.docker.internal//with docker-compose- db
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Successfully connect to the database")

	db.AutoMigrate(&models.User{})
	Database.Db = db
}
