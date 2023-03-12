package models

import g "github.com/jinzhu/gorm"

type StoreStatus struct {
	g.Model
	Name			string		`json:"name" gorm:"varchar(255);not null"`
	Minimum_exp		int			`json:"exp" gorm:"not null"`
	Maker_id		int
	Store			[]Store
}