package interfaces

import m "github.com/forumGamers/store-service/models"

type Store struct {
	m.Store
	AvgRating   float64 `json:"avg_rating" gorm:"-"`
	RatingCount int     `json:"rating_count" gorm:"-"`
}