package models

import g "github.com/jinzhu/gorm"

type Item struct {
	g.Model
	Name 			string 			`gorm:";type:varchar(255);NOT NULL" json:"name"`
	Image			string			`json:"image"`
	Store_id 		uint 			`gorm:";NOT NULL"`
	Status			string			`json:"status"`		
	Ratings 		[]ItemRating	`gorm:";foreignKey:item_id;references:id"`
	Store			Store			`gorm:";foreignKey:store_id;references:id"`	
}