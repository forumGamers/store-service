package services

import (
	"strconv"

	"errors"

	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func CheckAvailablity(itemId string, amount string,errCh chan error,itemCh chan m.Item) {
	var data m.Item
	Id, err := strconv.ParseInt(itemId, 10, 64)

	if err != nil {
		errCh <- errors.New("Invalid data")
		itemCh <- m.Item{}
		return
	}

	if err := l.GetDb().Model(m.Item{}).Where("id = ?", Id).Preload("Store").Preload("Store.Vouchers").First(&data).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errCh <- errors.New("Data not found")
			itemCh <- m.Item{}
			return
		} else {
			errCh <- err
			itemCh <- m.Item{}
			return
		}
	}

	if a, err := strconv.ParseInt(amount, 10,64); err != nil {
		errCh <- errors.New("Invalid data")
		itemCh <- m.Item{}
		return
	} else {
		if data.Stock < int(a) {
			errCh <- errors.New("Stock is not enough")
			itemCh <- m.Item{}
			return
		}
	}

	errCh <- nil
	itemCh <- data
}

func GetItem(id interface{},dataCh chan m.Item,errCh chan error){
	var data m.Item
	if err := l.GetDb().Model(m.Item{}).Where("id = ?",id).Preload("Store").First(&data).Error ; err != nil {
		if err == gorm.ErrRecordNotFound {
			errCh <- errors.New("Data not found")
			dataCh <- m.Item{}
			return
		}else {
			errCh <- err
			dataCh <- m.Item{}
			return
		}
	}

	errCh <- nil
	dataCh <- data
}