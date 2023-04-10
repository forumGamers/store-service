package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	cfg "github.com/forumGamers/store-service/config"
	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateItem(c *gin.Context) {
	id := c.Param("storeId")
	var item m.Item

	name, stock, price, description, discount :=
		c.PostForm("name"),
		c.PostForm("stock"),
		c.PostForm("price"),
		c.PostForm("description"),
		c.PostForm("discount")

	if name == "" {
		panic("Invalid data")
	} else {
		checkCh := make(chan error)
		go func(name string, id string) {
			var data m.Item
			if err := getDb().Model(m.Item{}).Where("name = ? and store_id = ?", name, id).First(&data).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					checkCh <- nil
					return
				}
			}

			checkCh <- errors.New("Conflict")
		}(name, id)

		if err := <-checkCh; err != nil {
			panic(err.Error())
		}
	}

	if s, r := strconv.ParseInt(stock, 10, 64); r == nil {
		item.Stock = int(s)
	} else {
		panic("Invalid data")
	}

	if p, r := strconv.ParseInt(price, 10, 64); r == nil {
		item.Price = int(p)
	} else {
		panic("Invalid data")
	}

	if d, r := strconv.ParseInt(discount, 10, 64); r == nil {
		item.Discount = int(d)
	} else {
		item.Discount = 0
	}

	if Id, err := strconv.ParseInt(id, 10, 64); err != nil {
		panic("Invalid data")
	} else {
		errCh := make(chan error)
		storeCh := make(chan uint)
		storeNameCh := make(chan string)
		go func(id int64) {
			var store m.Store
			if err := getDb().Where("owner_id = ?", id).First(&store).Error; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- 0
				storeNameCh <- ""
				return
			}

			errCh <- nil
			storeCh <- store.ID
			storeNameCh <- store.Name
		}(Id)

		if err := <-errCh; err != nil {
			panic(err.Error())
		}

		item.Store_id = <-storeCh
		storeName := <-storeNameCh
		item.Slug = h.SlugGenerator(name + " by " + storeName)
	}

	item.Name = name
	item.Description = description

	var img string

	if image, err := c.FormFile("image"); err == nil {

		if err := c.SaveUploadedFile(image, "uploads/"+image.Filename); err != nil {
			panic(err.Error())
		}

		file, _ := os.Open("uploads/" + image.Filename)

		data, errParse := ioutil.ReadAll(file)

		if errParse != nil {
			panic("Failed to parse image")
		}

		urlCh := make(chan string)
		fileIdCh := make(chan string)
		errCh := make(chan error)

		go func(data []byte, image string) {
			url, fileId, errUpload := cfg.UploadImage(data, image, "itemImage")

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
		}(data, image.Filename)

		select {
		case url := <-urlCh:
			if url == "" {
				panic("Internal Server Error")
			} else {
				file.Close()
				item.Image = url
				item.ImageId = <-fileIdCh
			}
		case err := <-errCh:
			if err != nil {
				panic(err.Error())
			}
		}

		img = image.Filename
	}

	errCh := make(chan error)

	go func(item *m.Item) {
		if err := getDb().Create(&item).Error; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		errCh <- nil
	}(&item)

	if err := <-errCh; err != nil {
		panic(err.Error())
	}

	if img != "" {
		if err := os.Remove("uploads/" + img); err != nil {
			fmt.Println(err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func UpdateItemDesc(c *gin.Context){
	id := c.Request.Header.Get("id")
	storeId := c.Request.Header.Get("storeId")
	itemId := c.Param("id")

	if id == "" {
		panic("Forbidden")
	}

	desc := c.PostForm("description")

	if desc == "" {
		panic("Invalid data")
	}

	errCh := make(chan error)

	Id,r := strconv.ParseInt(id,10,64)

	if r != nil {
		panic("Forbidden")
	}

	go func(id int,storeId string,itemId string,desc string){
		var data m.Item
		if err := getDb().Where("store_id = ? and id = ?",storeId,itemId).Preload("Store").First(&data).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		if data.Store.Owner_id != id {
			errCh <- errors.New("Forbidden")
			return
		}

		data.Description = desc

		if err := getDb().Model(m.Item{}).Save(&data).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		errCh <- nil
	}(int(Id),storeId,itemId,desc)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message" : "success"})
}

func UpdateItemImage(c *gin.Context){
	var item m.Item

	image , r := c.FormFile("image")

	id := c.Param("id")

	if r != nil {
		panic("Invalid data")
	}

	itemCh := make(chan m.Item)
	errCheckCh := make(chan error)

	go func(id string){
		var check m.Item

		if err := getDb().Where("id = ?",id).First(&check).Error ; err != nil {
			errCheckCh <- errors.New("Data not found")
			itemCh <- check
			return
		}

		errCheckCh <- nil
		itemCh <- check
	}(id)

	if err := <- errCheckCh ; err != nil {
		panic(err.Error())
	}

	item = <- itemCh

	if err := c.SaveUploadedFile(image,"uploads/"+image.Filename) ; err != nil {
		panic(err.Error())
	}

	file,_ := os.Open("uploads/"+image.Filename)

	data,errParse := ioutil.ReadAll(file)
	
	if errParse != nil {
		panic(errParse.Error())
	}

	urlCh := make(chan string)
	fileIdCh := make(chan string)
	errCh := make(chan error)

	go func(data []byte ,image *multipart.FileHeader,fileId string){

		url,id,err := cfg.UpdateImage(data,image.Filename,"itemImage",fileId)

		if err != nil {
			errCh <- errors.New(err.Error())
			urlCh <- url
			fileIdCh <- id
			return
		}

		urlCh <- url
		fileIdCh <- id
		errCh <- nil
	}(data ,image,item.ImageId)


	select {
		case err := <- errCh :
			if err != nil {
				panic(err.Error())
			}
		case url := <- urlCh :
			file.Close()
			item.Image = url
			item.ImageId = <- fileIdCh
		}

	os.Remove("uploads/"+image.Filename)

	errUpdate := make(chan error)

	go func(item m.Item){
		if err := getDb().Save(&item).Error ; err != nil {
			errUpdate <- errors.New(err.Error())
			return
		}

		errUpdate <- nil
	}(item)

	if err := <- errUpdate ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func AddStock(c *gin.Context){
	id := c.Param("id")

	stock := c.PostForm("stock")

	s,r := strconv.ParseInt(stock,10,64)

	if r != nil {
		panic("Invalid data")
	}

	errCh := make(chan error)

	go func (id string,stock int)  {
		if err := getDb().Model(m.Item{}).Where("id = ?",id).Update("stock",stock).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
			}else {
				errCh <- err
			}
		}

		errCh <- nil
	}(id,int(s))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message" : "success"})
}

func UpdatePrice(c *gin.Context){
	id := c.Param("id")

	price := c.PostForm("price")

	p,r := strconv.ParseInt(price,10,64)

	if r != nil {
		panic("Invalid data")
	}

	errCh := make(chan error)

	go func (id string,price int)  {
		if err := getDb().Model(m.Item{}).Where("id = ?",id).Update("price",price).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
			}else {
				errCh <- err
			}
		}

		errCh <- nil
	}(id,int(p))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message" : "success"})
}

func UpdateName(c *gin.Context){
	storeId := c.Request.Header.Get("storeId")
	id := c.Param("id")
	user := c.Request.Header.Get("id")

	Id,r := strconv.ParseInt(user,10,64)

	if r != nil {
		panic("Forbidden")
	}

	name := c.PostForm("name")
	var item m.Item

	if name == "" {
		panic("Invalid data")
	}else {
		checkCh := make(chan error)

		go func(name string,storeId string){
			var data m.Item
			if err := getDb().Model(m.Item{}).Where("name = ? and store_id = ?",name,storeId).Preload("Store").First(&data).Error ; err != nil {
				if err == gorm.ErrRecordNotFound {
					checkCh <- nil
					return
				}
			}

			checkCh <- errors.New("Conflict")
		}(name,storeId)

		if err := <- checkCh ; err != nil {
			panic(err.Error())
		}

		itemCh := make(chan m.Item)

		go func (id string) {
			var data m.Item
			getDb().Model(m.Item{}).Where("id = ?",id).Preload("Store").First(&data)

			itemCh <- data
		}(id)

		if item = <- itemCh ; item.Store.Owner_id != int(Id) {
			panic("Forbidden")
		}
	}

	errCh := make(chan error)
	slug := h.SlugGenerator(name+" by "+item.Store.Name)

	go func (id string,name string,slug string)  {
		if err := getDb().Model(m.Item{}).Where("id = ?",id).Update(m.Item{Name: name,Slug: slug}).Error ; err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}(id,name,slug)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func UpdateItemDiscount(c *gin.Context){
	itemId := c.Param("id")

	user := c.Request.Header.Get("id")
	storeId := c.Request.Header.Get("storeId")

	discount := c.PostForm("discount")

	disc,er := strconv.ParseInt(discount,10,64) 

	if er != nil {
		panic("Invalid data")
	}

	if user == "" {
		panic("Forbidden")
	}

	Id,r := strconv.ParseInt(user,10,64)

	if r != nil {
		panic("Forbidden")
	}

	errCh := make(chan error)

	go func(id int,storeId string,itemId string,discount int){
		var data m.Item
		if err := getDb().Where("store_id = ? and id = ?",storeId,itemId).Preload("Store").First(&data).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		if data.Store.Owner_id != id {
			errCh <- errors.New("Forbidden")
			return
		}

		data.Discount = discount

		if err := getDb().Model(m.Item{}).Save(&data).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		errCh <- nil
	}(int(Id),storeId,itemId,int(disc))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}