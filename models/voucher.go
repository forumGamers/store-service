package models

import g "github.com/jinzhu/gorm"

type Voucher struct {
	g.Model
	Name		string		`json:"name"`
	Discount	int			`json:"discount"`
	Cashback	int			`json:"cashback"`
	Store_id	int			`json:"store_id"`
	Store		Store
}