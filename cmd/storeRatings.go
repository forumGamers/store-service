package cmd

import (
	"errors"
	"net/http"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	v "github.com/forumGamers/store-service/validations"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func RateStore(c *gin.Context){
	store := c.Param("storeId")
	id := h.GetUser(c).Id

	r := c.PostForm("rate")

	if r == "" {
		panic("Invalid data")
	}

	rate,err := v.CheckRates(r)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"message": err.Error()})
		return
	}

	storeId,err := strconv.Atoi(store)

	if err != nil {
		panic("Invalid data")
	}

	checkStore := make(chan error)
	errCheck := make(chan error)
	errCh := make(chan error)

	go func(storeId int) {
		var store m.Store

		if err := getDb().Model(m.Store{}).Where("id = ?",storeId).First(&store).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				checkStore <- errors.New("Data not found")
				return
			}

			checkStore <- err
			return
		}

		checkStore <- nil
	} (storeId)

	go func (id,storeId int)  {
		var data m.StoreRating

		if err := getDb().Model(m.StoreRating{}).Where("store_id = ? and user_id = ?",storeId,id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCheck <- nil
				return
			}

			errCheck <- err
			return
		}

		errCheck <- errors.New("Conflict")
	} (id,storeId)

	go func(id,storeId int,rate int){
		if err := <- checkStore ; err != nil {
			errCh <- err
			return
		}

		if err := <- errCheck ; err != nil {
			errCh <- err
			return
		}

		data := m.StoreRating{
			User_id: uint(id),
			Store_id: uint(storeId),
			Rate: rate,
		}

		if err := getDb().Model(m.StoreRating{}).Create(&data).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(id,storeId,rate)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}