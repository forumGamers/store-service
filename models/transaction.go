package models

import (
	"time"

	g "github.com/jinzhu/gorm"
)

type Transaction struct {
	g.Model
	Date				time.Time	`json:"date"`
	User_id				int			`json:"User_id"`
	Payment_method		string		`json:"payment_method"`
	Item_id				uint		`json:"Item_id"`
	Store_id			uint		`json:"Store_id"`
	Store				Store		`gorm:"foreignKey:store_id;references:id"`
	Item				Item		`gorm:"foreignKey:item_id;references:id"`
}