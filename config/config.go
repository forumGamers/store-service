package config

import (
	"fmt"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var Db *gorm.DB

func Connection(){
	database,err := gorm.Open("postgres","user=postgres password=qwertyui host=127.0.0.1 port=5432 dbname=store-service sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("connection success")

	database.AutoMigrate(&m.StoreStatus{},&m.Experience{},&m.StoreRating{},&m.ItemRating{},&m.Store{},&m.Item{},&m.Transaction{})

	h.SetFK(database)

	Db = database
}