package models

import g "github.com/jinzhu/gorm"

type Follower struct {
	g.Model
	User_id				uint	`json:"userId"`
	Store_id			int		`json:"storeId" gorm:"references:store_id"`
	Store				Store	
}