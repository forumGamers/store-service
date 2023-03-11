package services

import (
	"errors"

	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func GetVoucher(id int,dataCh chan m.Voucher,errCh chan error) {
	var data m.Voucher

	if id == 0 {
		errCh <- errors.New("skip")
		dataCh <- m.Voucher{}
		return
	}

	if err := l.GetDb().Model(m.Voucher{}).Where("id = ?",id).First(&data).Error ; err != nil {
		if err == gorm.ErrRecordNotFound {
			errCh <- errors.New("Data not found")
			dataCh <- m.Voucher{}
			return
		}else {
			errCh <- err
			dataCh <- m.Voucher{}
			return
		}
	}

	errCh <- nil
	dataCh <- data
}

func ExpForStore(discount int,cashback int,stock int) int {
	return int((discount + cashback ) / stock)
}

func ExpForUser(discount int,cashback int,stock int) int {
	return (discount + cashback) / stock
}