package loaders

import (
	c "github.com/forumGamers/store-service/config"
	"github.com/jinzhu/gorm"
)

func GetDb() *gorm.DB {
	return c.Db
}