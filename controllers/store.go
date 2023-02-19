package controllers

import (
	"net/http"
	"strconv"

	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getDb() *gorm.DB {
	return l.GetDb()
}


func CreateStore(c *gin.Context){
	var store m.Store

	name,image,description := c.PostForm("name"),c.PostForm("image"),c.PostForm("description")

	owner_id := c.Request.Header.Get("id")

	if name == ""  || owner_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest,gin.H{"message" : "invalid data"})
		return
	}

	store.Name = name

	store.Image = image

	store.Description = description

	id,_ := strconv.ParseInt(owner_id,10,64)

	store.Owner_id = int(id)

	err := make(chan bool)

	go func () {
		res := getDb().Create(&store)

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