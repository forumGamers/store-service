package query

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)


func GetAllItem(c *gin.Context){
	name,
	minDate,
	maxDate,
	store,
	minStock,
	maxStock,
	status,
	minPrice,
	maxPrice,
	minDiscount,
	maxDiscount,
	limit,
	page :=
	c.Query("name"),
	c.Query("minDate"),
	c.Query("maxDate"),
	c.Query("store"),
	c.Query("minStock"),
	c.Query("maxStock"),
	c.Query("status"),
	c.Query("minPrice"),
	c.Query("maxPrice"),
	c.Query("minDiscount"),
	c.Query("maxDiscount"),
	c.Query("limit"),
	c.Query("page")

	var item []m.Item

	errCh := make(chan error)
	itemCh := make(chan []m.Item)

	go func(
		name string,
		minDate string,
		maxDate string,
		store string,
		minStock string,
		maxStock string,
		status string,
		minPrice string,
		maxPrice string,
		minDiscount string,
		maxDiscount string,
		limit string,
		page string,
	){
		var data []m.Item

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
				itemCh <- nil
				return
			}
	
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at BETWEEN ? and ?")
			args = append(args,minDate,maxDate)
		}else if minDate != "" {
			if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at >= ?")
			args = append(args,minDate)
		}else if maxDate != "" {
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at <= ?")
			args = append(args, maxDate)
		}

		if store != "" {
			if _,err := strconv.ParseInt(store,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"store_id = ?")
			args = append(args, store)
		}

		if status != "" {
			r := regexp.MustCompile(`\W`)
			x := r.ReplaceAllString(status,"")
			query = h.QueryBuild(query,"status = ?")
			args = append(args, x)
		}

		if minPrice != "" && maxPrice != "" {
			if _,err := strconv.ParseInt(minPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			if _,err := strconv.ParseInt(maxPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"(price BETWEEN ? and ?)")
			args = append(args, minPrice,maxPrice)
		}else if minPrice != "" {
			if _,err := strconv.ParseInt(minPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"price >= ?")
			args = append(args, minPrice)
		}else if maxPrice != "" {
			if _,err := strconv.ParseInt(maxPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}

			query = h.QueryBuild(query,"price <= ?")
			args = append(args, maxPrice)
		}

		if minDiscount != "" && maxDiscount != "" {
			if _,err := strconv.ParseInt(minDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			if _,err := strconv.ParseInt(maxDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"(discount BETWEEN ? and ?)")
			args = append(args, minDiscount,maxDiscount)
		}else if minDiscount != "" {
			if _,err := strconv.ParseInt(minDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"discount >= ?")
			args = append(args, minDiscount)
		}else if maxDiscount != "" {
			if _,err := strconv.ParseInt(maxDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"discount <= ?")
			args = append(args, maxDiscount)
		}

		if limit == "" {
			lmt = 10
		}else {
			if lm,err := strconv.ParseInt(limit,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				itemCh <- nil
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
				itemCh <- nil
				return
			}else {
				pg = int(p)
			}
		}

		getDb().Model(m.Item{}).Where(query,args...).Preload("Store").Offset((pg - 1) * lmt).Limit(lmt).Find(&data)

		if len(data) < 1 {
			errCh <- errors.New("Data not found")
			itemCh <- nil
			return
		}

		itemCh <- data

		errCh <- nil
	}(
		name,
		minDate,
		maxDate,
		store,
		minStock,
		maxStock,
		status,
		minPrice,
		maxPrice,
		minDiscount,
		maxDiscount,
		limit,
		page,
	)

	select {
	case err := <- errCh :
		if err != nil {
			panic(err.Error())
		}
	case item = <- itemCh :
		c.JSON(http.StatusOK,item)
		return
	}
}

func GetItemBySlug(c *gin.Context){
	slug := c.Param("slug")

	r := regexp.MustCompile(`[^\w%.]`)
	res := r.ReplaceAllString(slug,"")

	errCh := make(chan error)
	itemCh := make(chan m.Item)

	go func (str string){
		var data m.Item
		if err := getDb().Model(m.Item{}).Where("slug = ?",str).Preload("Store").First(&data).Error ; err != nil {
			errCh <- errors.New(err.Error())
			itemCh <- m.Item{}
			return
		}

		errCh <- nil
		itemCh <- data
	}(res)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	item := <- itemCh

	c.JSON(http.StatusOK,item)
}

func GetItemByStoreId(c *gin.Context){
	storeId := c.Param("storeId")
	name,
	minDate,
	maxDate,
	minStock,
	maxStock,
	status,
	minPrice,
	maxPrice,
	minDiscount,
	maxDiscount,
	limit,
	page := 
	c.Query("name"),
	c.Query("minDate"),
	c.Query("maxDate"),
	c.Query("minStock"),
	c.Query("maxStock"),
	c.Query("status"),
	c.Query("minPrice"),
	c.Query("maxPrice"),
	c.Query("minDiscount"),
	c.Query("maxDiscount"),
	c.Query("limit"),
	c.Query("page")

	errCh := make(chan error)
	dataCh := make(chan []m.Item)

	go func (
		id string,
		name string,
		minDate string,
		maxDate string,
		minStock string,
		maxStock string,
		status string,
		minPrice string,
		maxPrice string,
		minDiscount string,
		maxDiscount string,
		limit string,
		page string,
		)  {
		var data []m.Item 

		var res string

		var args []interface{}

		query := "store_id = ?"
		args = append(args, id)

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
				dataCh <- nil
				return
			}
	
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at BETWEEN ? and ?")
			args = append(args,minDate,maxDate)
		}else if minDate != "" {
			if _,err := time.Parse("30-12-2022",minDate) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at >= ?")
			args = append(args,minDate)
		}else if maxDate != "" {
			if _,err := time.Parse("30-12-2022",maxDate) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"created_at <= ?")
			args = append(args, maxDate)
		}

		if status != "" {
			r := regexp.MustCompile(`\W`)
			x := r.ReplaceAllString(status,"")
			query = h.QueryBuild(query,"status = ?")
			args = append(args, x)
		}

		if minPrice != "" && maxPrice != "" {
			if _,err := strconv.ParseInt(minPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			if _,err := strconv.ParseInt(maxPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"(price BETWEEN ? and ?)")
			args = append(args, minPrice,maxPrice)
		}else if minPrice != "" {
			if _,err := strconv.ParseInt(minPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"price >= ?")
			args = append(args, minPrice)
		}else if maxPrice != "" {
			if _,err := strconv.ParseInt(maxPrice,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}

			query = h.QueryBuild(query,"price <= ?")
			args = append(args, maxPrice)
		}

		if minDiscount != "" && maxDiscount != "" {
			if _,err := strconv.ParseInt(minDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			if _,err := strconv.ParseInt(maxDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"(discount BETWEEN ? and ?)")
			args = append(args, minDiscount,maxDiscount)
		}else if minDiscount != "" {
			if _,err := strconv.ParseInt(minDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"discount >= ?")
			args = append(args, minDiscount)
		}else if maxDiscount != "" {
			if _,err := strconv.ParseInt(maxDiscount,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
				return
			}
	
			query = h.QueryBuild(query,"discount <= ?")
			args = append(args, maxDiscount)
		}

		if limit == "" {
			lmt = 10
		}else {
			if lm,err := strconv.ParseInt(limit,10,64) ; err != nil {
				errCh <- errors.New(err.Error())
				dataCh <- nil
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
				dataCh <- nil
				return
			}else {
				pg = int(p)
			}
		}

		if err := getDb().Model(m.Item{}).Where(query,args...).Preload("Store").Offset((pg - 1) * lmt).Limit(lmt).Find(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- []m.Item{}
				return
			}else {
				errCh <- errors.New(err.Error())
				dataCh <- []m.Item{}
				return
			}
		}

		errCh <- nil
		dataCh <- data
	}(
		storeId,
		name,
		minDate,
		maxDate,
		minStock,
		maxStock,
		status,
		minPrice,
		maxPrice,
		minDiscount,
		maxDiscount,
		limit,
		page,
	)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	items := <- dataCh

	c.JSON(http.StatusOK,items)
}

func GetItemSlugByStoreId(c *gin.Context){
	storeId := c.Param("storeId")

	id,err := strconv.Atoi(storeId)

	type itemSlug struct {
		Slug	string
	}

	if err != nil {
		panic("Invalid data")
	}

	dataCh := make(chan []itemSlug)
	errCh := make(chan error)

	go func(id int){
		var data []itemSlug

		if err := getDb().Raw("select slug from items i where store_id = ?",id).Find(&data).Error ; err != nil {
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
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh
	
	c.JSON(http.StatusOK,data)
}