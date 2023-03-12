package models

import g "github.com/jinzhu/gorm"

type Favorite struct {
	g.Model
	User_id		uint		`json:"user_id"`
	Item_id		int			`json:"item_id"`
	Item		Item
}