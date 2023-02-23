package models

import g "github.com/jinzhu/gorm"

type Store struct {
	g.Model
	Name 				string 			`gorm:";varchar(20);NOT NULL" json:"name"`
	Image 				string 			`json:"image"`
	Description 		string 			`gorm:";type:text" json:"description"`
	Owner_id			int				`json:"Owner_id" gorm:";NOT NULL"`
	Exp					int				`json:"exp" gorm:";not null"`
	Active				bool			`json:"active" gorm:"default:true"`
	Items 				[]Item			`gorm:";constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:store_id"`
	Ratings				[]StoreRating	`gorm:";constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:store_id"`
	Transactions 		[]Transaction	`gorm:";foreignKey:store_id"`
}

func (s *Store) BeforeCreate(tx *g.DB) (err error){
	s.Active = true
	s.Exp = 0
	return nil
}