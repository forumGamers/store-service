package models

import g "github.com/jinzhu/gorm"

type StoreRating struct {
	g.Model
	Store_id		uint		`json:"store_id" gorm:";NOT NULL"`
	User_id			uint		`json:"user_id" gorm:";NOT NULL"`
	Rate			int			`json:"rate" gorm:";NOT NULL"`
	Store			Store		`gorm:";foreignKey:store_id;references:id"`
}