package models

import g "github.com/jinzhu/gorm"

type Voucher struct {
	g.Model
	Name			string		`json:"name"`
	Discount		int			`json:"discount"`
	Cashback		int			`json:"cashback"`
	Store_id		int			`json:"store_id"`
	Period			int			`json:"period"`
	Status			string		`json:"status"`
	Stock			int			`json:"stock"`
	PointForStore	int			`json:"point_for_store"`
	PointForUser	int			`json:"point_for_user"`
	Store			Store
}

func (v *Voucher) BeforeCreate(tx *g.DB) error {
	v.Status = "Active"
	return nil
}