package cmd

import (
	l "github.com/forumGamers/store-service/loaders"
	"github.com/jinzhu/gorm"
)

func getDb() *gorm.DB {
	return l.GetDb()
}