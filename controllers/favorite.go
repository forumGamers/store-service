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

func GetMyFavorite(c *gin.Context){
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
	dataCh := make(chan []m.Favorite)

	go func(userId int,page int,limit int){
		var data []m.Favorite
		if err := getDb().Model(m.Favorite{}).Where("user_id = ?",userId).Offset((page - 1) * limit).Limit(limit).Find(&data).Error ; err != nil {
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

	favorite := <- dataCh

	c.JSON(http.StatusOK,favorite)
}