package models

import g "github.com/jinzhu/gorm"

type ItemRating struct {
	g.Model
	Item_id		uint		`json:"item_id" gorm:";NOT NULL"`
	User_id		uint		`json:"user_id" gorm:";NOT NULL"`
	Rate		int			`json:"rate" gorm:";NOT NULL"`
	Item		Item		`gorm:";foreignKey:item_id;references:id"`
}