package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
)

func GetAllStoreStatus(c *gin.Context){
	name,page,limit := 
	c.Query("name"),
	c.Query("page"),
	c.Query("limit")

	ch := make(chan []m.StoreStatus)
	errCh := make(chan error)

	go func(name string,page string,limit string){
		var store_status []m.StoreStatus

		tx := getDb().Model(&m.StoreStatus{})

		var args []interface{}

		var query string

		var lmt int

		if name != "" {
			r := regexp.MustCompile(`\W`)
			result := r.ReplaceAllString(name,"")
			query = h.QueryBuild(query,"name ILIKE ?")
			args = append(args, "%"+result+"%")
		}
	
		if limit == "" {
			lmt = 10
		}else {
			if l,err := strconv.ParseInt(limit,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				ch <- nil
				return
			}else {
				lmt = int(l)
			}
		}
	
		if page != "" {
			p ,err:= strconv.ParseInt(page,10,64)
	
			if err != nil {
				errCh <- errors.New(err.Error())
				ch <- nil
				return
			}
	
			offset := (int(p) - 1) * lmt
	
			tx.Offset(offset)
		}
	
		tx.Limit(lmt)
	
		
		tx.Where(query,args...).Find(&store_status)

		if len(store_status) < 1 {
			errCh <- errors.New("Data not found") 
			ch <- nil
			return
		}

		ch <- store_status
	}(name,page,limit)

	select {
	case storeStatus := <- ch :
		c.JSON(http.StatusOK,storeStatus)
		return
	case err := <- errCh : 
		panic(err.Error())
	}
}