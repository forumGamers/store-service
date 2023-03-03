package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

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
			query = h.QueryBuild(query,"status = ?")
			args = append(args, status)
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