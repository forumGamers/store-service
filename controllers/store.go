package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"

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
	name,minDate,maxDate,owner,active,minExp,maxExp := 
	c.Query("name"),
	c.Query("minDate"),
	c.Query("maxDate"),
	c.Query("owner"),
	c.Query("active"),
	c.Query("minExp"),
	c.Query("maxExp")

	var store []m.Store

	var wg sync.WaitGroup

	var res string

	tx := getDb().Model(m.Store{})

	wg.Add(1)

	if name != "" {
		r := regexp.MustCompile(`\W`)
		res = r.ReplaceAllString(name,"")
		tx.Where("name ILIKE ?",res)
	}

	if minDate != "" && maxDate != "" {
		if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
			panic(err.Error())
		}

		if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
			panic(err.Error())
		}

		tx.Where("created_at BETWEEN ? and ?",minDate,maxDate)

	}else if minDate != "" {
		if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
			panic(err.Error())
		}

		tx.Where("created_at >= ?",minDate)

	}else if maxDate != "" {
		if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
			panic(err.Error())
		}

		tx.Where("created_at <= ?",maxDate)
	}

	if owner != "" {
		if _,err := strconv.ParseInt(owner,10,64) ; err != nil {
			panic(err.Error())
		}

		tx.Where("owner = ?",owner)
	}

	if active != "" {
		if _,err := strconv.ParseBool(active) ; err != nil {
			panic(err.Error())
		}

		tx.Where("active = ?",active)

	}

	if minExp != "" && maxExp != "" {
		if _,err := strconv.ParseInt(minExp,10,64) ; err != nil {
			panic(err.Error())
		}

		if _,err := strconv.ParseInt(maxExp,10,64) ; err != nil {
			panic(err.Error())
		}

		tx.Where("exp BETWEEN ? and ?",minExp,maxExp)

	}else if minExp != "" {
		if _,err := strconv.ParseInt(minExp,10,64) ; err != nil {
			panic(err.Error())
		}

		tx.Where("exp >= ?",minExp)

	}else if maxExp != "" {
		if _,err := strconv.ParseInt(maxExp,10,64) ; err != nil {
			panic(err.Error())
		}

		tx.Where("exp <= ?",maxExp)
	}

	go func ()  {
		defer wg.Done()
		tx.Find(&store)
	}()

	wg.Wait()

	if len(store) < 1 {
		panic("Data not found")
	}

	c.JSON(http.StatusOK,gin.H{"data" : store})
}