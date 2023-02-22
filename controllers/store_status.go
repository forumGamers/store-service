package controllers

import (
	"net/http"
	"regexp"
	"strconv"

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

	err := make(chan bool)

	go func (){
		res := getDb().Create(&store_status)

		if res.Error != nil {
			err <- true
		}else {
			err <- false
		}
	}()

	if <- err == false {
		c.JSON(http.StatusCreated,gin.H{"message":"success"})
		return
	}else {
		c.JSON(http.StatusInternalServerError,gin.H{"message" : "Internal Server Error"})
		return
	}
}

func GetAllStoreStatus(c *gin.Context){
	name,page := c.Query("name"),c.Query("page")

	ch := make(chan []m.StoreStatus)
	errCh := make(chan string)

	go func(){
		var store_status []m.StoreStatus

		tx := getDb().Model(&m.StoreStatus{})

		if name != "" {
			r := regexp.MustCompile(`\W`)
			result := r.ReplaceAllString(name," ")
			tx.Where("name ILIKE ?",result)
		}
	
		limit := 10
	
		if page != "" {
			p ,err:= strconv.ParseInt(page,10,64)
	
			if err != nil {
				panic("Invalid data")
			}
	
			offset := (int(p) - 1) * limit
	
			tx.Offset(offset)
		}
	
		tx.Limit(limit)
	
		
		tx.Find(&store_status)

		if len(store_status) < 1 {
			errCh <- "Data not found"
			return
		}

		ch <- store_status
	}()

	select {
	case storeStatus := <- ch :
		c.JSON(http.StatusOK,gin.H{"data":storeStatus})
		return
	case err := <- errCh : 
		if err == "Data not found" {
			panic(err)
		}else {
			panic("Internal Server Error")
		}
	}
}

func UpdateStoreStatusName(c *gin.Context){
	var store_status m.StoreStatus

	name := c.PostForm("name")

	id := c.Param("id")

	if name == "" {
		panic("Invalid data")
	}

	if err := getDb().Where("id = ?", id).First(&store_status).Error; err != nil {
		panic("Data not found")
	}

	store_status.Name = name

	if err := getDb().Save(&store_status).Error ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message": "success"})
}