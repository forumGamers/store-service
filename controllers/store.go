package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"

	"strconv"
	"time"

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

	owner_id := c.Request.Header.Get("id")

	if name == ""  {
		panic("Invalid data")
	}

	if owner_id == "" {
		panic("Forbidden")
	}

	id,er := strconv.ParseInt(owner_id,10,64)

	if er != nil {
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
	}(name,int(id))

	if err := <- checkCh ; err != nil {
		panic(err.Error())
	}

	var img string

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

	store.Owner_id = int(id)

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
	id := c.Request.Header.Get("id")

	name := c.PostForm("name")

	errCh := make(chan error)

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

	id := c.Request.Header.Get("id")

	if id == "" {
		panic("Forbidden")
	}

	Id,er := strconv.ParseInt(id,10,64)

	if er != nil {
		panic(er.Error())
	}

	errCh := make(chan error)

	go func (id int)  {
		if err := getDb().Where("owner_id = ?",id).First(&store).Error ; err != nil {
			errCh <- errors.New("Data not found")
			return
		}

		if Id != int64(store.Owner_id) {
			errCh <- errors.New("Forbidden")
			return
		}

		store.Description = desc

		if err := getDb().Save(&store).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(int(Id))

	if err := <- errCh; err != nil {
		panic(err.Error())
	}
	
	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func UpdateStoreImage(c *gin.Context){
	var store m.Store

	image , r := c.FormFile("image")

	id := c.Request.Header.Get("id")

	if r != nil {
		panic("Invalid data")
	}

	storeCh := make(chan m.Store)
	errCheckCh := make(chan error)

	go func(id string){
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

	image , r := c.FormFile("background")

	id := c.Request.Header.Get("id")

	if r != nil {
		panic("Invalid data")
	}

	storeCh := make(chan m.Store)
	errCheckCh := make(chan error)

	go func(id string){
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

func GetAllStores(c *gin.Context){
	name,minDate,maxDate,owner,active,minExp,maxExp,page,limit := 
	c.Query("name"),
	c.Query("minDate"),
	c.Query("maxDate"),
	c.Query("owner"),
	c.Query("active"),
	c.Query("minExp"),
	c.Query("maxExp"),
	c.Query("page"),
	c.Query("limit")

	var store []m.Store

	errCh := make(chan error)
	storeCh := make(chan []m.Store)

	go func (
		name string,
		minDate string,
		maxDate string,
		owner string,
		active string,
		minExp string,
		maxExp string,
		page string,
		limit string,
		){

		var data []m.Store

		var res string

		var args []interface{}

		var query string

		var pg int

		var lmt int

		if name != "" {
			r := regexp.MustCompile(`\W`)
			res = r.ReplaceAllString(name,"")
			query = h.QueryBuild(query,"name ILIKE ?")
			args = append(args, "%"+res+"%")
		}
	
		if minDate != "" && maxDate != "" {
			if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at BETWEEN ? and ?")
			args = append(args,minDate,maxDate)
	
		}else if minDate != "" {
			if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at >= ?")
			args = append(args,minDate)
	
		}else if maxDate != "" {
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at <= ?")
			args = append(args, maxDate)
		}
	
		if owner != "" {
			if _,err := strconv.ParseInt(owner,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"owner_id = ?")
			args = append(args, owner)
		}
	
		if active != "" {
			if _,err := strconv.ParseBool(active) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"active = ?")
			args = append(args, active)
		}
	
		if minExp != "" && maxExp != "" {
			if _,err := strconv.ParseInt(minExp,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			if _,err := strconv.ParseInt(maxExp,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"(exp BETWEEN ? and ?)")
			args = append(args, minExp,maxExp)
	
		}else if minExp != "" {
			if _,err := strconv.ParseInt(minExp,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"exp >= ?")
			args = append(args, minExp)
	
		}else if maxExp != "" {
			if _,err := strconv.ParseInt(maxExp,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"exp <= ?")
			args = append(args, maxExp)
		}

		if limit == "" {
			lmt = 10
		}else {
			if lm,err := strconv.ParseInt(limit,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}else {
				lmt = int(lm)
			}
		}

		if page == "" {
			pg = 1
		}else {
			if p,err := strconv.ParseInt(page,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				storeCh <- nil
				return
			}else {
				pg = int(p)
			}
		}

		getDb().Model(m.Store{}).Where(query,args...).Preload("Items").Offset((pg - 1) * lmt).Limit(lmt).Find(&data)

		if len(data) < 1 {
			errCh <- errors.New("Data not found")
			storeCh <- nil
			return
		}

		storeCh <- data

		errCh <- nil
	}(
		name,
		minDate,
		maxDate,
		owner,
		active,
		minExp,
		maxExp,
		page,
		limit,
	)

	select {
	case err := <- errCh : 
		if err != nil {
			panic(err.Error())
		}
	case store = <- storeCh :
		c.JSON(http.StatusOK,store)
		return
	}
}

func GetStoreById(c *gin.Context){
	id := c.Param("id")

	errCh := make(chan error)

	ch := make(chan m.Store)

	go func (id string){
		var store m.Store

		if err := getDb().Where("id = ?",id).Preload("Items").Find(&store).Error ; err != nil {
			errCh <- errors.New("Data not found")
			ch <- m.Store{}
			return
		}else {
			errCh <- nil
			ch <- store
		}
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	store := <- ch

	c.JSON(http.StatusOK,store)
}

func DeactiveStore(c *gin.Context){
	id := c.Request.Header.Get("id")

	if id == "" {
		panic("Forbidden")
	}

	Id,er := strconv.ParseInt(id,10,64)

	if er != nil {
		panic(er.Error())
	}

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
	}(int(Id))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func ReactivedStore(c *gin.Context){
	id := c.Request.Header.Get("id")

	if id == "" {
		panic("Forbidden")
	}

	Id,er := strconv.ParseInt(id,10,64)

	if er != nil {
		panic(er.Error())
	}

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
	}(int(Id))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}