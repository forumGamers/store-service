package controllers

import (
	"errors"
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetStoreFollower(c *gin.Context){
	id := c.Param("storeId")
	limit := c.Query("limit")
	page := c.Query("page")
	
	var pg int
	var lmt int

	if page != "" {
		p,r := strconv.ParseInt(page,10,64)

		if r != nil {
			panic("Invalid params")
		}
		pg = int(p)
	}

	if limit != "" {
		l,r := strconv.ParseInt(limit,10,64)

		if r != nil {
			panic("Invalid params")
		}
		lmt = int(l)
	}

	errCh := make(chan error)
	dataCh := make(chan m.Follower)

	go func(id string,limit int,page int){
		var data m.Follower
		if err := getDb().Model(m.Follower{}).Offset((page - 1) * limit).Limit(limit).Find(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- m.Follower{}
				return
			}else {
				errCh <- err
				dataCh <- m.Follower{}
				return
			}
		}

		dataCh <- data
		errCh <- nil
	}(id,lmt,pg)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	follower := <- dataCh

	c.JSON(http.StatusOK,follower)
}