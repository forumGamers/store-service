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

	if name == ""  {
		panic("Invalid data")
	}

	if owner_id == "" {
		panic("Forbidden")
	}

	store.Name = name

	store.Image = image

	store.Description = description

	id,_ := strconv.ParseInt(owner_id,10,64)

	store.Owner_id = int(id)

	err := make(chan error)

	go func () {
		res := getDb().Create(&store)

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

func UpdateStoreName(c *gin.Context){
	var store m.Store
	id := c.Request.Header.Get("id")

	name := c.PostForm("name")

	if id == "" {
		panic("Forbidden")
	}

	if name == "" {
		panic("Invalid data")
	}

	Id,er := strconv.ParseInt(id,10,64)

	if er != nil {
		panic(er.Error())
	}

	if err := getDb().Where("id = ?",Id).First(&store).Error ; err != nil {
		panic("Data not found")
	}

	if int(Id) != store.Owner_id {
		panic("Forbidden")
	}

	store.Name = name

	if err := getDb().Save(&store).Error; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusContinue,gin.H{"message":"success"})
}