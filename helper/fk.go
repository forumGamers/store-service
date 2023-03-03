package helper

import (
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func SetFK(g *gorm.DB){
	g.Model(m.Store{}).AddUniqueIndex("idx_store_name","name")

	g.Model(&m.Item{}).AddUniqueIndex("idx_slug_items","slug").AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.Transaction{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.Transaction{}).AddForeignKey("item_id","items(id)","CASCADE","CASCADE")

	g.Model(&m.StoreRating{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.ItemRating{}).AddForeignKey("item_id","items(id)","CASCADE","CASCADE")
}