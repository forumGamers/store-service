package cmd

import (
	"errors"
	"net/http"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddCart(c *gin.Context){
	user := h.GetUser(c)
	itemId := c.Param("itemId")

	item,err := strconv.ParseInt(itemId,10,64)

	if err != nil {
		panic("Invalid data")
	}

	errCh := make(chan error)

	go func(itemId int){
		if err := getDb().Model(m.Cart{}).Where("item_id = ?",itemId).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- nil
				return
			}else {
				errCh <- err
				return
			}
		}
		errCh <- errors.New("Conflict")
	}(int(item))

	errCreate := make(chan error)

	go func(userId int,itemId int){

		if err := <- errCh ; err != nil {
			errCreate <- err
			return
		}

		var data m.Cart
		data.Item_id = itemId
		data.User_id = uint(userId)

		if err := getDb().Model(m.Cart{}).Create(&data).Error ; err != nil {
			errCreate <- err
			return
		}

		errCreate <- nil
	}(user.Id,int(item))

	if err := <- errCreate ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func RemoveCart(c *gin.Context){
	cart := c.Param("id")
	user := c.Request.Header.Get("id")

	if user == "" {
		panic("Forbidden")
	}

	id,er := strconv.ParseInt(cart,10,64)

	if er != nil {
		panic("Invalid data")
	}

	errCh := make(chan error)

	go func(id int,userId string){
		if err := getDb().Model(m.Cart{}).Where("user_id = ?",userId).Delete(m.Cart{},id).Error ; err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}(int(id),user)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK,gin.H{"message":"success"})
}