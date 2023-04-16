package query

import (
	"errors"
	"net/http"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func GetCart(c *gin.Context){
	user := h.GetUser(c)
	limit,page := 
	c.Query("limit"),
	c.Query("page")

	var lmt int
	var pg int

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
	}(user.Id,int(pg),int(lmt))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	cart := <- dataCh

	c.JSON(http.StatusOK,cart)
}