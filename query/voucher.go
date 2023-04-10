package query

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)



func GetAllVoucher(c *gin.Context){
	name,store,page,limit := c.Query("name"),c.Query("store"),c.Query("page"),c.Query("limit")

	errCh := make(chan error)
	dataCh := make(chan []m.Voucher)

	go func(name string,store string,page string,limit string){
		var data []m.Voucher
		var args []interface{}
		var pg int
		var lmt int
		query := ""
		if name != "" {
			r := regexp.MustCompile(`\W`)
			res := r.ReplaceAllString(name,"")
			query = h.QueryBuild(query,"name ILIKE ?")
			args = append(args, "%"+res+"%")
		}

		if store != "" {
			if _,err := strconv.ParseInt(store,10,64) ; err != nil {
				errCh <- errors.New("Invalid params")
				dataCh <-nil
				return
			}
		}

		if page != "" {
			p,err := strconv.ParseInt(page,10,64) 
			if err != nil {
				pg = 1
			}else {
				pg = int(p)
			}
		}else {
			pg = 1
		}

		if limit != "" {
			l,err := strconv.ParseInt(limit,10,64)
			if err != nil {
				lmt = 10
			}else {
				lmt = int(l)
			}
		}else {
			lmt = 10
		}
		query = h.QueryBuild(query,"store_id = ?")
		args = append(args, store)
		if err := getDb().Model(m.Voucher{}).Where(query,args...).Offset((pg - 1) * lmt).Limit(lmt).Find(&data).Error ; err != nil {
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
	}(name,store,page,limit)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh

	c.JSON(http.StatusOK,data)
}

func GetVoucherById(c *gin.Context){
	id := c.Param("id")

	errCh := make(chan error)
	dataCh := make(chan m.Voucher)

	go func(id string){
		var data m.Voucher
		if err := getDb().Model(m.Voucher{}).Where("id = ?",id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- m.Voucher{}
				return
			}else {
				errCh <- err
				dataCh <- m.Voucher{}
				return
			}
		}

		dataCh <- data
		errCh <- nil
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh

	c.JSON(http.StatusOK,data)
}
