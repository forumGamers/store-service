package cmd

import (
	"errors"
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddCart(c *gin.Context){
	user := c.Request.Header.Get("id")
	itemId := c.Param("itemId")

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

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	errCreate := make(chan error)

	go func(userId int,itemId int){
		var data m.Cart
		data.Item_id = itemId
		data.User_id = uint(userId)

		if err := getDb().Model(m.Cart{}).Create(&data).Error ; err != nil {
			errCreate <- err
			return
		}

		errCreate <- nil
	}(int(id),int(item))

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