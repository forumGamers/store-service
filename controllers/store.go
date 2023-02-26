package controllers

import (
	"errors"
	"io/ioutil"
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

	image,er := c.FormFile("image")

	owner_id := c.Request.Header.Get("id")

	if er != nil {
		panic(er.Error())
	}

	if name == ""  {
		panic("Invalid data")
	}

	if owner_id == "" {
		panic("Forbidden")
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
	errCh := make(chan error)

	go func (data []byte, image string){
		if url,err := cfg.UploadImage(data,image) ;err != nil {
			errCh <- errors.New(err.Error())
			urlCh <- ""
		}else {
			errCh <- nil
			urlCh <- url
		}	
		return
	}(data,image.Filename)

	select {
	case err := <- errCh :
		if err != nil {
			panic(err.Error())
		}
	case url := <- urlCh :
		store.Image = url
	}

	store.Name = name

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
		os.Remove("uploads/"+image.Filename)
		c.JSON(http.StatusCreated,gin.H{"message":"success"})
		return
	}else {
		panic(<- err)
	}
}

func UpdateStoreName(c *gin.Context){
	var store m.Store
	id := c.Request.Header.Get("id")

	storeId := c.Param("id")

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

	go func ()  {
		if err := getDb().Where("id = ?",storeId).First(&store).Error ; err != nil {
			errCh <- errors.New("Data not found")
		}
	
		if int(Id) != store.Owner_id {
			errCh <- errors.New("Forbidden")
		}
	
		store.Name = name
	
		if err := getDb().Save(&store).Error; err != nil {
			errCh <- err
		}
		
		errCh <- nil
	}()

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func UpdateStoreDesc(c *gin.Context){
	var store m.Store

	desc := c.PostForm("description")

	id := c.Request.Header.Get("id")

	storeId := c.Param("id")

	if id == "" {
		panic("Forbidden")
	}

	Id,er := strconv.ParseInt(id,10,64)

	if er != nil {
		panic(er.Error())
	}

	errCh := make(chan error)

	go func ()  {
		if err := getDb().Where("id = ?",storeId).First(&store).Error ; err != nil {
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
	}()

	if err := <- errCh; err != nil {
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

		getDb().Model(m.Store{}).Where(query,args...).Offset((pg - 1) * lmt).Limit(lmt).Find(&data)

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
		panic(err.Error())
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

		if err := getDb().Where("id = ?",id).Find(&store).Error ; err != nil {
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