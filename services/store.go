package services

import (
	"errors"

	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

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