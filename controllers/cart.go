package controllers

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

func GetCart(c *gin.Context){
	user := c.Request.Header.Get("id")
	limit,page := 
	c.Query("limit"),
	c.Query("page")

	var lmt int
	var pg int

	id,err := strconv.ParseInt(user,10,64)

	if err != nil {
		panic("Forbidden")
	}

	if limit == "" {
		lmt = 10
	}else {
		l,r := strconv.ParseInt(limit,10,64)

		if r != nil {
			panic("Invalid params")
		}

		lmt = int(l)
	}

	if page == "" {
		pg = 1
	}else {
		p,r := strconv.ParseInt(page,10,64)

		if r != nil {
			panic("Invalid params")
		}

		pg = int(p)
	}

	errCh := make(chan error)
	dataCh := make(chan []m.Cart)

	go func(userId int,page int,limit int){
		var data []m.Cart
		if err := getDb().Model(m.Cart{}).Where("user_id = ?",userId).Offset((page - 1) * limit).Limit(limit).Find(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- nil
				return
			}else {
				errCh <- err
				dataCh <- nil
				return
			}
		}

		errCh <- nil
		dataCh <- data
	}(int(id),int(pg),int(lmt))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	cart := <- dataCh

	c.JSON(http.StatusOK,cart)
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