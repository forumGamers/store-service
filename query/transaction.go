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

func GetAllTransaction(c *gin.Context){
	minValue,maxValue,status,itemId,storeId,paymentMethod,page,limit :=
	c.Query("minValue"),
	c.Query("maxValue"),
	c.Query("status"),
	c.Query("itemId"),
	c.Query("storeId"),
	c.Query("paymentMethod"),
	c.Query("page"),
	c.Query("limit")

	errCh := make(chan error)
	dataCh := make(chan []m.Transaction)

	go func(
		minValue string,
		maxValue string,
		status string,
		itemId string,
		storeId string,
		paymentMethod string,
		page string,
		limit string,
		){
			var minVal int
			var maxVal int
			var item int
			var store int
			var pg int
			var lmt int
			query := ""
			var args []interface{}
			var data []m.Transaction

			if minValue != "" && maxValue != "" {
				if min ,err := strconv.ParseInt(minValue,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					minVal = int(min)
				}

				if max,err := strconv.ParseInt(maxValue,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					maxVal = int(max)
				}

				query = h.QueryBuild(query,"value BETWEEN ? and ?")
				args = append(args, minVal,maxVal)
			}else if minValue != "" {
				if min ,err := strconv.ParseInt(minValue,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					minVal = int(min)
				}

				query = h.QueryBuild(query,"value >= ?")
				args = append(args, minVal)
			}else if maxValue != "" {
				if max,err := strconv.ParseInt(maxValue,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					maxVal = int(max)
				}

				query = h.QueryBuild(query,"value <= ?")
				args = append(args, maxVal)
			}

			if status != "" {
				r := regexp.MustCompile(`[^\w\s.] `)
				res := r.ReplaceAllString(status,"")
				query = h.QueryBuild(query,"status = ?")
				args = append(args,res)
			}

			if itemId != "" {
				if i,err := strconv.ParseInt(itemId,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return 
				}else {
					item = int(i)
				}

				query = h.QueryBuild(query,"item_id = ?")
				args = append(args, item)
			}

			if storeId != "" {
				if s,err := strconv.ParseInt(storeId,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					store = int(s)
				}

				query = h.QueryBuild(query,"store_id = ?")
				args = append(args, store)
			}

			if limit == "" {
				lmt = 10
			}else {
				if lm,err := strconv.ParseInt(limit,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					lmt = int(lm)
				}
			}
	
			if page == "" {
				pg = 1
			}else {
				if p,err := strconv.ParseInt(page,10,64) ; err != nil {
					errCh <- errors.New("Invalid params")
					dataCh <- nil
					return
				}else {
					pg = int(p)
				}
			}

			if err := getDb().Model(m.Transaction{}).Where(query,args...).Preload("Item").Preload("Store").Offset((pg - 1) * lmt).Limit(lmt).Find(&data).Error ; err != nil {
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

			if len(data) < 1 {
				errCh <- errors.New("Data not found")
				dataCh <- nil
				return
			}

			errCh <- nil
			dataCh <- data
		}(
			minValue,
			maxValue,
			status,
			itemId,
			storeId,
			paymentMethod,
			page,
			limit,
		)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	transaction := <- dataCh 

	c.JSON(http.StatusOK,transaction)
}