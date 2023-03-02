package controllers

import (
	"errors"
	"fmt"
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
		storeNameCh := make(chan string)
		go func (id int64) {
			var store m.Store
			if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- 0
				storeNameCh <- ""
				return
			}

			errCh <- nil
			storeCh <- store.ID
			storeNameCh <- store.Name
		}(Id)

		if err := <- errCh ; err != nil {
			panic(err.Error())
		}

		item.Store_id = <- storeCh
		storeName := <- storeNameCh
		item.Slug = h.SlugGenerator(name+" by "+storeName)
	}

	item.Name = name
	item.Description = description

	var img string

	if image,err := c.FormFile("image") ; err == nil {

		if err := c.SaveUploadedFile(image,"uploads/"+image.Filename) ; err != nil {
			panic(err.Error())
		}
	
		file,_ := os.Open("uploads/"+image.Filename)
	
		data,errParse := ioutil.ReadAll(file)
		
		if errParse != nil {
			panic("Failed to parse image")
		}
	
		urlCh := make (chan string)
		fileIdCh := make(chan string)
		errCh := make(chan error)
	
		go func (data []byte, image string){
			url ,fileId ,errUpload := cfg.UploadImage(data,image,"itemImage")
	
			if errUpload != nil {
				urlCh <- ""
				fileIdCh <- ""
				errCh <- errors.New("Bad Gateway")
				return
			}
	
			urlCh <- url
			fileIdCh <- fileId
			errCh <- nil
			return
		}(data,image.Filename)
	
		select {
		case url := <- urlCh :
			if url == "" {
				panic("Internal Server Error")
			}else {
				file.Close()
				item.Image = url
				item.ImageId = <- fileIdCh
			}
		case err := <- errCh :
			if err != nil {
				panic(err.Error())
			}
		}

		img = image.Filename
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

	if err := os.Remove("uploads/"+img) ; err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusCreated,gin.H{"message" : "success"})
}