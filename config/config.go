package config

import (
	"fmt"
	"os"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/joho/godotenv"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var Db *gorm.DB

func Connection(){
	godotenv.Load()

	DATABASE_URL := os.Getenv("DATABASE_URL")

	if DATABASE_URL == "" {
		DATABASE_URL = "user=postgres password=qwertyui host=127.0.0.1 port=5432 dbname=store-service sslmode=disable"
	}

	database,err := gorm.Open("postgres",DATABASE_URL)

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("connection success")

	database.AutoMigrate(
		&m.StoreStatus{},
		&m.Cart{},
		&m.Favorite{},
		&m.Voucher{},
		&m.StoreRating{},
		&m.ItemRating{},
		&m.Store{},
		&m.Item{},
		&m.Log{},
	)

	h.SetFK(database)

	Db = database
}