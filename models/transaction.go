package models

import (
	g "github.com/jinzhu/gorm"
)

type Transaction struct {
	g.Model
	User_id				uint		`json:"User_id"`
	Description			string		`json:"description"`
	Payment_method		string		`json:"payment_method" gorm:";NOT NULL"`
	Item_id				uint		`json:"Item_id"`
	Store_id			uint		`json:"Store_id"`
	Store				Store		`gorm:";foreignKey:store_id;references:id"`
	Item				Item		`gorm:";foreignKey:item_id;references:id"`
	Status				string		`json:"status"`
	Value				int			`json:"value"`
	Amount				int			`json:"amount"`
	MessageForSeller	string		`json:"message"`
	Fee					int			`json:"fee"`
}

func (t *Transaction) BeforeCreate(tx *g.DB) error {
	t.Status = "On Process"
	return nil
}