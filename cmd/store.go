package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	cfg "github.com/forumGamers/store-service/config"
	h "github.com/forumGamers/store-service/helper"
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

	name,description := c.PostForm("name"),c.PostForm("description")

	owner_id := h.GetUser(c).Id

	if name == ""  {
		panic("Invalid data")
	}

	checkCh := make(chan error)

	go func (name string, ownerId int) {
		var check m.Store

		getDb().Where("name = ? or owner_id = ?",name,ownerId).First(&check)

		if c := check.Name == name ; c != false {
			checkCh <- errors.New("name is already use")
			return
		}

		if c := check.Owner_id == ownerId ; c != false {
			checkCh <- errors.New("you already have a store")
			return
		}

		checkCh <- nil
	}(name,owner_id)

	if err := <- checkCh ; err != nil {
		panic(err.Error())
	}

	var img string

	if image,err := c.FormFile("image") ; err == nil {

		if err := h.IsImage(image) ; err != nil {
			panic(err.Error())
		}

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
	
		go func (data []byte, image string){
			url ,fileId ,errUpload := cfg.UploadImage(data,image,"storeImage")
	
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
				store.Image = url
				store.ImageId = <- fileIdCh
			}
		case err := <- errCh :
			if err != nil {
				panic(err.Error())
			}
		}
		img = image.Filename
	}

	var bgImg string

	if bg,err := c.FormFile("background") ; err == nil {
		if err := c.SaveUploadedFile(bg,"uploads/"+bg.Filename) ; err != nil {
			panic(err.Error())
		}

		file,_ := os.Open("uploads/"+bg.Filename)

		data,errParse := ioutil.ReadAll(file)
		
		if errParse != nil {
			panic(errParse.Error())
		}
	
		bgCh := make (chan string)
		bgIdCh := make(chan string)
		errBgCh := make(chan error)
	
		go func (data []byte, image string){
			url ,fileId ,errUpload := cfg.UploadImage(data,image,"storeImage")
	
			if errUpload != nil {
				bgCh <- ""
				bgIdCh <- ""
				errBgCh <- errors.New("Bad Gateway")
				return
			}
	
			bgCh <- url
			bgIdCh <- fileId
			errBgCh <- nil
			return
		}(data,bg.Filename)
	
		select {
		case url := <- bgCh :
			if url == "" {
				panic("Internal Server Error")
			}else {
				file.Close()
				store.Background = url
				store.BackgroundId = <- bgIdCh
			}
		case err := <- errBgCh :
			if err != nil {
				panic(err.Error())
			}
		}
		bgImg = bg.Filename
	}

	store.Name = name

	store.Description = description

	store.Owner_id = owner_id

	store.Status_id = 1

	err := make(chan error)

	go func () {

		if errCreate := getDb().Create(&store).Error ; errCreate != nil {
			err <- errCreate
		}

		err <- nil
	}()

	if <- err == nil {
		if img != "" {
			if err := os.Remove("uploads/"+img) ; err != nil {
				fmt.Println(err)
			}
		}
		if bgImg != "" {
			if err := os.Remove("uploads/"+bgImg) ; err != nil {
				fmt.Println(err)
			}
		}
		c.JSON(http.StatusCreated,gin.H{"message":"success"})
		return
	}else {
		panic(<- err)
	}
}

func UpdateStoreName(c *gin.Context){
	var store m.Store
	Id := h.GetUser(c).Id

	name := c.PostForm("name")

	errCh := make(chan error)

	if name == "" {
		panic("Invalid data")
	}

	go func (id int)  {
		if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
			errCh <- errors.New("Data not found")
		}
	
		if int(id) != store.Owner_id {
			errCh <- errors.New("Forbidden")
		}
	
		store.Name = name
	
		if err := getDb().Save(&store).Error; err != nil {
			errCh <- err
		}
		
		errCh <- nil
	}(int(Id))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func UpdateStoreDesc(c *gin.Context){
	var store m.Store

	desc := c.PostForm("description")

	Id := h.GetUser(c).Id

	errCh := make(chan error)

	go func (id int)  {
		if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
			errCh <- errors.New("Data not found")
			return
		}

		if Id != store.Owner_id {
			errCh <- errors.New("Forbidden")
			return
		}

		store.Description = desc

		if err := getDb().Save(&store).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(Id)

	if err := <- errCh; err != nil {
		panic(err.Error())
	}
	
	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func UpdateStoreImage(c *gin.Context){
	var store m.Store

	image , err := c.FormFile("image")

	if err != nil {
		panic("Invalid data")
	}

	if err := h.IsImage(image) ; err != nil {
		panic(err.Error())
	}

	id := h.GetUser(c).Id

	storeCh := make(chan m.Store)
	errCheckCh := make(chan error)

	go func(id int){
		var check m.Store

		if err := getDb().Where("owner_id = ?",id).First(&check).Error ; err != nil {
			errCheckCh <- errors.New("Data not found")
			storeCh <- check
			return
		}

		errCheckCh <- nil
		storeCh <- check
	}(id)

	if err := <- errCheckCh ; err != nil {
		panic(err.Error())
	}

	store = <- storeCh

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

		url,id,err := cfg.UpdateImage(data,image.Filename,"storeImage",fileId)

		if err != nil {
			errCh <- errors.New(err.Error())
			urlCh <- url
			fileIdCh <- id
			return
		}

		urlCh <- url
		fileIdCh <- id
		errCh <- nil
	}(data ,image,store.ImageId)


	select {
		case err := <- errCh :
			if err != nil {
				panic(err.Error())
			}
		case url := <- urlCh :
			file.Close()
			store.Image = url
			store.ImageId = <- fileIdCh
		}

	os.Remove("uploads/"+image.Filename)

	errUpdate := make(chan error)

	go func(store m.Store){
		if err := getDb().Save(&store).Error ; err != nil {
			errUpdate <- errors.New(err.Error())
			return
		}

		errUpdate <- nil
	}(store)

	if err := <- errUpdate ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func UpdateStoreBg(c *gin.Context){
	var store m.Store

	id := h.GetUser(c).Id

	image , err := c.FormFile("background")

	if err != nil {
		panic("Invalid data")
	}

	if err := h.IsImage(image) ; err != nil {
		panic(err.Error())
	}

	storeCh := make(chan m.Store)
	errCheckCh := make(chan error)

	go func(id int){
		var check m.Store

		if err := getDb().Where("owner_id = ?",id).First(&check).Error ; err != nil {
			errCheckCh <- errors.New("Data not found")
			storeCh <- check
			return
		}

		errCheckCh <- nil
		storeCh <- check
	}(id)

	if err := <- errCheckCh ; err != nil {
		panic(err.Error())
	}

	store = <- storeCh

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

		url,id,err := cfg.UpdateImage(data,image.Filename,"storeImage",fileId)

		if err != nil {
			errCh <- errors.New(err.Error())
			urlCh <- url
			fileIdCh <- id
			return
		}

		urlCh <- url
		fileIdCh <- id
		errCh <- nil
	}(data ,image,store.ImageId)


	select {
		case err := <- errCh :
			if err != nil {
				panic(err.Error())
			}
		case url := <- urlCh :
			file.Close()
			store.Background = url
			store.BackgroundId = <- fileIdCh
		}

	os.Remove("uploads/"+image.Filename)

	errUpdate := make(chan error)

	go func(store m.Store){
		if err := getDb().Save(&store).Error ; err != nil {
			errUpdate <- errors.New(err.Error())
			return
		}

		errUpdate <- nil
	}(store)

	if err := <- errUpdate ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func DeactiveStore(c *gin.Context){
	Id := h.GetUser(c).Id

	errCh := make(chan error)

	go func(id int){
		var store m.Store

		if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		store.Active = false

		if err := getDb().Save(&store).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		errCh <- nil
	}(Id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func ReactivedStore(c *gin.Context){
	Id := h.GetUser(c).Id

	errCh := make(chan error)

	go func(id int){
		var store m.Store

		if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		store.Active = true

		if err := getDb().Save(&store).Error ; err != nil {
			errCh <- errors.New(err.Error())
			return
		}

		errCh <- nil
	}(Id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}