package controllers

import (
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
)

func CreateStoreStatus(c *gin.Context){
	var store_status m.StoreStatus

	name,minimum_exp := c.PostForm("name"),c.PostForm("minimum_exp")

	maker_id := c.Request.Header.Get("id")

	if name == "" || minimum_exp == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"message" : "Invalid data"})
		return
	}

	if maker_id == "" {
		c.AbortWithStatusJSON(http.StatusForbidden,gin.H{"message" : "Forbidden"})
		return
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