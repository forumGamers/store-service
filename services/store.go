package services

import (
	"errors"

	i "github.com/forumGamers/store-service/interfaces"
	"github.com/forumGamers/store-service/loaders"
	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func getDb() *gorm.DB {
	return loaders.GetDb()
}

func GetStore(id interface{},storeCh chan m.Store,errCh chan error) {
	var data m.Store
	
	if err := l.GetDb().Model(m.Store{}).Where("id = ?",id).Preload("Vouchers").Preload("Items").First(&data).Error ; err != nil {
		if err == gorm.ErrRecordNotFound {
			errCh <- errors.New("Data not found")
			storeCh <- m.Store{}
			return
		}else {
			errCh <- err
			storeCh <- m.Store{}
			return
		}
	}

	errCh <- nil
	storeCh <- data
}

func VoucherCheck(voucher m.Voucher,storeId uint) bool {
	if voucher.Store_id == int(storeId) {
		return true
	}
	return false
}

func GetStoreByCondition(data *i.Store,cond string,id int) error {
	if err := getDb().Model(m.Store{}).Where(cond,id).
			Select(`stores.*, AVG(store_ratings.rate) AS avg_rating, COUNT(store_ratings.*) AS rating_count`).
			Joins("LEFT JOIN store_ratings ON store_ratings.store_id = stores.id").
			Group("stores.id").
			Preload("Items",func(db *gorm.DB) *gorm.DB {
				return db.Select("items.*, NULL as store")
			}).Preload("StoreStatus").First(&data).Error ; err != nil {
				return err
			}
	return nil
}