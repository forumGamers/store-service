package controllers

import (
	"errors"
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddVoucher(c *gin.Context){
	name,discount,cashback,period := c.PostForm("name"),c.PostForm("discount"),c.PostForm("cashback"),c.PostForm("periode")
	
	id,storeId := c.Request.Header.Get("id"),c.Request.Header.Get("storeId")

	checkCh := make(chan error)

	Id,r := strconv.ParseInt(id,10,64)

	if r != nil {
		panic("Forbidden")
	}

	go func(id int,storeId string) {
		var data m.Store
		if err := getDb().Model(m.Voucher{}).Where("id = ? and owner_id = ?",storeId,id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				checkCh <- errors.New("You have not store yet")
				return
			}else {
				checkCh <- err
				return
			}
		}

		if data.Owner_id != id {
			checkCh <- errors.New("Forbidden")
			return
		}

		checkCh <- nil
	}(int(Id),storeId)

	if err := <- checkCh ; err != nil {
		panic(err.Error())
	}

	errCh := make(chan error)

	go func(name string,discount string,cashback string,period string,storeId string){
		var data m.Voucher

		disc,errDisc := strconv.ParseInt(discount,10,64)
		
		if errDisc != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		cb,errCb := strconv.ParseInt(cashback,10,64)

		if errCb != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		p,errP := strconv.ParseInt(period,10,64)

		if errP != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		id,errId := strconv.ParseInt(storeId,10,64)

		if errId != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		data.Name = name
		data.Discount = int(disc)
		data.Cashback = int(cb)
		data.Period = int(p)
		data.Store_id = int(id)

		if err := getDb().Model(m.Voucher{}).Create(&data).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(name,discount,cashback,period,storeId)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}
