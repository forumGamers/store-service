package models

import g "github.com/jinzhu/gorm"

type Experience struct {
	g.Model
	Exp			int 		`json:"exp" gorm:"default:0"`
}