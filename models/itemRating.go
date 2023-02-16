package models

import g "github.com/jinzhu/gorm"

type ItemRating struct {
	g.Model
	Item_id		uint
	Rate		int			`json:"rate"`
	Item		Item		`gorm:"foreignKey:item_id;references:id"`
}