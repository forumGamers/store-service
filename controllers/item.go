package controllers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	cfg "github.com/forumGamers/store-service/config"
	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
)


func CreateItem(c *gin.Context){
	id := c.Request.Header.Get("id")
	var item m.Item

	name,stock,price,description,discount := 
	c.PostForm("name"),
	c.PostForm("stock"),
	c.PostForm("price"),
	c.PostForm("description"),
	c.PostForm("discount")

	if s,r := strconv.ParseInt(stock,10,64) ; r == nil {
		item.Stock = int(s)
	}else {
		panic("Invalid data")
	}

	if p,r := strconv.ParseInt(price,10,64) ; r == nil {
		item.Price = int(p)
	}else {
		panic("Invalid data")
	}

	if d,r := strconv.ParseInt(discount,10,64) ; r == nil {
		item.Discount = int(d)
	}else {
		item.Discount = 0
	}

	if Id,err := strconv.ParseInt(id,10,64) ; err != nil {
		panic("Invalid data")
	}else {
		errCh := make(chan error)
		storeCh := make(chan uint)
		go func (id int64) {
			var store m.Store
			if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- 0
			}

			errCh <- nil
			storeCh <- store.ID
		}(Id)

		if err := <- errCh ; err != nil {
			panic(err.Error())
		}

		item.Store_id = <- storeCh
	}

	item.Name = name
	item.Description = description
	item.Slug = h.SlugGenerator(name)

	if image,err := c.FormFile("image") ; err == nil {
		if err := c.SaveUploadedFile(image,"uploads/"+image.Filename) ; err != nil {
			panic(err.Error())
		}

		file,_ := os.Open("uploads/"+image.Filename)

		data,errParse := ioutil.ReadAll(file)
		
		if errParse != nil {
			panic(errParse.Error())
		}

		urlCh := make (chan string)
		fileIdCh := make(chan string)
		errCh := make(chan error)

		go func(data []byte ,imageName string) {
			url,fileId,err := cfg.UploadImage(data,imageName,"itemImage") 

			if err != nil {
				urlCh <- ""
				fileIdCh <- ""
				errCh <- errors.New(err.Error())
				return
			}

			urlCh <- url
			fileIdCh <- fileId
			errCh <- nil
		}(data,image.Filename)

		if err := <- errCh ; err != nil {
			panic(err.Error())
		}else {
			item.Image = <- urlCh
			item.ImageId = <- fileIdCh
		}

		os.Remove("uploads/"+image.Filename)
	}

	errCh := make(chan error)

	go func (item *m.Item){
		if err := getDb().Create(&item).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		errCh <- nil
	}(&item)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message" : "success"})
}