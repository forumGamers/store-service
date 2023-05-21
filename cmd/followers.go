package cmd

import (
	"errors"
	"net/http"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func FollowStoreById(c *gin.Context){
	store := c.Param("storeId")

	id := h.GetUser(c).Id

	storeId,err := strconv.Atoi(store)

	if err != nil {
		panic("Invalid data")
	}

	checkStore := make(chan error)
	conflictCheck := make(chan error)
	errCh := make(chan error)

	go func(storeId , id int){
		var data m.Store

		if err := getDb().Model(m.Store{}).Where("id = ?",storeId).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				checkStore <- errors.New("Data not found")
				return
			}

			checkStore <- err
			return
		}

		if data.Owner_id == id {
			checkStore <- errors.New("Unauthorize")
			return
		}

		checkStore <- nil
	}(storeId,id)

	go func(storeId,id int){
		var data m.Follower

		if err := getDb().Model(m.Follower{}).Where("store_id = ? and user_id = ?",storeId,id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				conflictCheck <- nil
				return
			}

			conflictCheck <- err
			return
		}

		conflictCheck <- errors.New("Conflict")
	}(storeId,id)

	go func(storeId,id int){
		if err := <- checkStore ; err != nil {
			errCh <- err
			return
		}

		if err := <- conflictCheck ; err != nil {
			errCh <- err
			return
		}

		data := m.Follower{
			User_id: uint(id),
			Store_id: storeId,
		}

		if err := getDb().Model(m.Follower{}).Create(&data).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(storeId,id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}