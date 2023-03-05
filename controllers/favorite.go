package controllers

import (
	"errors"
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddFavorite(c *gin.Context){
	user := c.Request.Header.Get("id")
	itemId := c.Param("id")

	id,r := strconv.ParseInt(user,10,64)

	if r != nil {
		panic("Forbidden")
	}

	item,er := strconv.ParseInt(itemId,10,64)

	if er != nil {
		panic("Invalid data")
	}

	errCh := make(chan error)

	go func(itemId int){
		if err := getDb().Model(m.Item{}).Where("id = ?",itemId).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				return
			}else {
				errCh <- err
			}
		}
		errCh <- nil
	}(int(item))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	errCreate := make(chan error)

	go func(userId int,itemId int){
		var data m.Favorite
		data.Item_id = itemId
		data.User_id = uint(userId)

		if err := getDb().Model(m.Favorite{}).Create(&data).Error ; err != nil {
			errCreate <- err
			return
		}

		errCreate <- nil
	}(int(id),int(item))

	if err := <- errCreate ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK,gin.H{"message":"success"})
}