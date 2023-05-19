package query

import (
	"errors"
	"net/http"
	"regexp"

	"strconv"
	"time"

	h "github.com/forumGamers/store-service/helper"
	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func getDb() *gorm.DB {
	return l.GetDb()
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
				errCh <- err
				storeCh <- nil
				return
			}
	
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at BETWEEN ? and ?")
			args = append(args,minDate,maxDate)
	
		}else if minDate != "" {
			if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at >= ?")
			args = append(args,minDate)
	
		}else if maxDate != "" {
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at <= ?")
			args = append(args, maxDate)
		}
	
		if owner != "" {
			if _,err := strconv.ParseInt(owner,10,64) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"owner_id = ?")
			args = append(args, owner)
		}
	
		if active != "" {
			if _,err := strconv.ParseBool(active) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"active = ?")
			args = append(args, active)
		}
	
		if minExp != "" && maxExp != "" {
			if _,err := strconv.ParseInt(minExp,10,64) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			if _,err := strconv.ParseInt(maxExp,10,64) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"(exp BETWEEN ? and ?)")
			args = append(args, minExp,maxExp)
	
		}else if minExp != "" {
			if _,err := strconv.ParseInt(minExp,10,64) ; err != nil {
				errCh <- err
				storeCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"exp >= ?")
			args = append(args, minExp)
	
		}else if maxExp != "" {
			if _,err := strconv.ParseInt(maxExp,10,64) ; err != nil {
				errCh <- err
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
				errCh <- err
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
		}
		errCh <- nil
		ch <- store
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	store := <- ch

	c.JSON(http.StatusOK,store)
}

func GetStoreName(c *gin.Context){
	id := c.Request.Header.Get("id")

	if _,err := strconv.ParseInt(id,10,64) ; err != nil {
		panic("Invalid data")
	}

	storeName := make(chan string)
	errCh := make(chan error)

	go func (id string){
		var data m.Store 

		if err := getDb().Model(m.Store{}).Where("owner_id = ?",id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				storeName <- ""
				return
			}
			errCh <- err
			storeName <- ""
			return
		}

		errCh <- nil
		storeName <- data.Name
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	name := <- storeName

	c.JSON(http.StatusOK,name)
}

func GetMyStore(c *gin.Context){
	id := h.GetUser(c).Id

	errCh := make(chan error)
	dataCh := make(chan m.Store)

	go func(id int){
		var data m.Store

		if err := getDb().Model(m.Store{}).Where("owner_id = ?",id).Preload("Items",func(db *gorm.DB) *gorm.DB {
			return db.Select("items.*, NULL as store")
		}).Preload("StoreStatus").First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- m.Store{}
				return
			}

			errCh <- err
			dataCh <- m.Store{}
			return
		}

		errCh <- nil
		dataCh <- data
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh 

	c.JSON(http.StatusOK,data)
}