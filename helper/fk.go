package helper

import (
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func SetFK(g *gorm.DB){
	g.Model(&m.Item{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.Transaction{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.Transaction{}).AddForeignKey("item_id","items(id)","CASCADE","CASCADE")

	g.Model(&m.StoreRating{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.ItemRating{}).AddForeignKey("item_id","items(id)","CASCADE","CASCADE")
}