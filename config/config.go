package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var Db *gorm.DB

func Connection(){
	database,err := gorm.Open("postgres","user=postgres password=qwertyui host=127.0.0.1 port=5432 dbname=store-service sslmode=disable")

	if err != nil {
		panic(err)
	}
	fmt.Println("connection success")

	Db = database
}