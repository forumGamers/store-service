package cmd

import (
	"errors"
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"

	"github.com/gin-gonic/gin"
)

func CreateStoreStatus(c *gin.Context) {
	var store_status m.StoreStatus

	name, minimum_exp := c.PostForm("name"), c.PostForm("minimum_exp")

	maker_id := c.Request.Header.Get("id")

	if name == "" || minimum_exp == "" {
		panic("Invalid data")
	}

	if maker_id == "" {
		panic("Forbidden")
	}

	store_status.Name = name

	exp, _ := strconv.ParseInt(minimum_exp, 10, 64)

	store_status.Minimum_exp = int(exp)

	id, _ := strconv.ParseInt(maker_id, 10, 64)

	store_status.Maker_id = int(id)

	err := make(chan error)

	go func() {
		res := getDb().Create(&store_status)

		if res.Error != nil {
			err <- res.Error
		} else {
			err <- nil
		}
	}()

	if <-err == nil {
		c.JSON(http.StatusCreated, gin.H{"message": "success"})
		return
	} else {
		panic(<-err)
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