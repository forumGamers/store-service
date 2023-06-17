package helper

import (
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func SetFK(g *gorm.DB){
	g.Model(m.Store{}).AddUniqueIndex("idx_store_name","name")

	g.Model(m.Store{}).AddForeignKey("status_id","store_statuses(id)","CASCADE","CASCADE")

	g.Model(&m.Item{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE").AddUniqueIndex("idx_slug_items","slug")

	g.Model(&m.Cart{}).AddIndex("idx_cart","item_id").AddForeignKey("item_id","items(id)","CASCADE","CASCADE")

	g.Model(m.Favorite{}).AddIndex("idx_favorite","item_id").AddForeignKey("item_id","items(id)","CASCADE","CASCADE")

	g.Model(&m.StoreRating{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")

	g.Model(&m.ItemRating{}).AddForeignKey("item_id","items(id)","CASCADE","CASCADE")

	g.Model(&m.Voucher{}).AddForeignKey("store_id","stores(id)","CASCADE","CASCADE")
}