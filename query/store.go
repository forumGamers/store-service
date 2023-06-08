package query

import (
	"errors"
	"net/http"
	"regexp"

	"strconv"
	"time"

	h "github.com/forumGamers/store-service/helper"
	i "github.com/forumGamers/store-service/interfaces"
	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	s "github.com/forumGamers/store-service/services"
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

	errCh := make(chan error)
	storeCh := make(chan []i.Store)

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

		var data []i.Store

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

		getDb().Model(m.Store{}).
		Where(query,args...).
		Select(`stores.*, AVG(store_ratings.rate) AS avg_rating, COUNT(store_ratings.*) AS rating_count`).
		Joins("LEFT JOIN store_ratings ON store_ratings.store_id = stores.id").
		Group("stores.id").
		Preload("Items").
		Preload("StoreStatus").
		Offset((pg - 1) * lmt).
		Limit(lmt).
		Find(&data)

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
	case store := <- storeCh :
		c.JSON(http.StatusOK,store)
		return
	}
}

func GetStoreById(c *gin.Context){
	id := c.Param("id")

	errCh := make(chan error)

	ch := make(chan i.Store)

	Id,err := strconv.Atoi(id)

	if err != nil {
		panic("Invalid data")
	}

	go func (id int){
		var store i.Store

		if err := s.GetStoreByCondition(&store,"stores.id = ?",id) ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				ch <- i.Store{}
				return
			}

		errCh <- err
		ch <- i.Store{}
		return
	}

	errCh <- nil
	ch <- store
	}(Id)

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
	dataCh := make(chan i.Store)

	go func(id int){
		var data i.Store

		if err := s.GetStoreByCondition(&data,"owner_id = ?",id) ; err != nil {
				if err == gorm.ErrRecordNotFound {
					errCh <- errors.New("Data not found")
					dataCh <- i.Store{}
					return
				}

			errCh <- err
			dataCh <- i.Store{}
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

func GetAllStoreId(c *gin.Context){
	type storeId struct {
		ID int
	}

	dataCh := make(chan []storeId)
	errCh := make(chan error)

	go func(){
		var data []storeId
		if err := getDb().Raw(`select id from stores s where deleted_at is null`).Find(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- nil
				return
			}

			errCh <- err
			dataCh <- nil
			return
		}

		errCh <- nil
		dataCh <- data
	}()

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh

	c.JSON(http.StatusOK,data)
}