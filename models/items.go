package models

import g "github.com/jinzhu/gorm"

type Item struct {
	g.Model
	Name 			string 			`gorm:";type:varchar(255);NOT NULL" json:"name"`
	Image			string			`json:"image"`
	ImageId			string			`json:"ImageId"`
	Store_id 		uint 			`gorm:";NOT NULL"`
	Status			string			`json:"status"`		
	Slug			string			`json:"slug"`
	Stock			int				`json:"stock"`
	Price			int				`json:"price"`
	Description		string			`json:"description"`
	Discount		int				`json:"discount"`
	Sold			int				`json:"sold"`
	Active			bool			`json:"active"`
	Ratings 		[]ItemRating	`gorm:";foreignKey:item_id;references:id"`
	Store			Store			`gorm:";foreignKey:store_id;references:id"`	
}

func (i *Item) BeforeCreate(tx *g.DB) error {
	if i.Stock > 0 {
		i.Status = "Available"
	}else {
		i.Status = "Not Available"
	}

	i.Active = true
	i.Sold = 0
	return nil
}