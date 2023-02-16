package models

import g "github.com/jinzhu/gorm"

type StoreRating struct {
	g.Model
	Store_id		uint		`json:"store_id"`
	Rate			int			`json:"rate"`
	Store			Store		`gorm:"foreignKey:store_id;references:id"`
}