package models

import g "github.com/jinzhu/gorm"

type Store struct {
	g.Model
	Name 				string 			`gorm:";varchar(20);NOT NULL" json:"name"`
	Image 				string 			`json:"image"`
	Description 		string 			`gorm:";type:text" json:"description"`
	Owner_id			int				`json:"Owner_id" gorm:";NOT NULL"`
	Items 				[]Item			`gorm:";constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:store_id"`
	Ratings				[]StoreRating	`gorm:";constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:store_id"`
	Transactions 		[]Transaction	`gorm:";foreignKey:store_id"`
}