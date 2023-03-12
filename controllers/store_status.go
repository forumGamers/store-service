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

func CreateStoreStatus(c *gin.Context){
	var store_status m.StoreStatus

	name,minimum_exp := c.PostForm("name"),c.PostForm("minimum_exp")

	maker_id := c.Request.Header.Get("id")

	if name == "" || minimum_exp == "" {
		panic("Invalid data")
	}

	if maker_id == "" {
		panic("Forbidden")
	}

	store_status.Name = name

	exp,_ := strconv.ParseInt(minimum_exp,10,64)

	store_status.Minimum_exp = int(exp)

	id,_ := strconv.ParseInt(maker_id,10,64)

	store_status.Maker_id = int(id)

	err := make(chan error)

	go func (){
		res := getDb().Create(&store_status)

		if res.Error != nil {
			err <- res.Error
		}else {
			err <- nil
		}
	}()

	if <- err == nil {
		c.JSON(http.StatusCreated,gin.H{"message":"success"})
		return
	}else {
		panic(<- err)
	}
}

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

func UpdateStoreStatusName(c *gin.Context){
	var store_status m.StoreStatus

	name := c.PostForm("name")

	id := c.Param("id")

	if name == "" {
		panic("Invalid data")
	}

	errCh := make(chan error)

	go func ()  {
		if err := getDb().Where("id = ?", id).First(&store_status).Error; err != nil {
			errCh <- errors.New("Data not found")
			return
		}
	
		store_status.Name = name
	
		if err := getDb().Save(&store_status).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}
	}()

	if err := <- errCh ;err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message": "success"})
}

func UpdateStoreStatusExp(c *gin.Context){
	var store_status m.StoreStatus

	exp := c.PostForm("exp")

	id := c.Param("id")

	if exp == "" {
		panic("Invalid data")
	}

	errCh := make(chan error)

	e,er := strconv.ParseInt(exp,10,64)

	if er != nil {
		panic(er.Error())
	}

	go func ()  {
		if err := getDb().Where("id = ?", id).First(&store_status).Error; err != nil {
			errCh <- errors.New("Data not found")
			return
		}

		store_status.Minimum_exp = int(e)

		if err := getDb().Save(&store_status).Error ; err != nil {
			errCh <- err
		}

		errCh <- nil
	}()

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}