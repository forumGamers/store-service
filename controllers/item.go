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
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)


func CreateItem(c *gin.Context){
	id := c.Param("storeId")
	var item m.Item

	name,stock,price,description,discount := 
	c.PostForm("name"),
	c.PostForm("stock"),
	c.PostForm("price"),
	c.PostForm("description"),
	c.PostForm("discount")

	if name == "" {
		panic("Invalid data")
	}else {
		checkCh := make(chan error)
		go func(name string,id string){
			var data m.Item
			if err := getDb().Model(m.Item{}).Where("name = ? and store_id = ?",name,id).First(&data).Error ; err != nil {
				if err == gorm.ErrRecordNotFound {
					checkCh <- nil
					return
				}
			}

			checkCh <- errors.New("Conflict")
		}(name,id)

		if err := <- checkCh ; err != nil {
			panic(err.Error())
		}
	}

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
	storeId := c.Param("storeId")
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

		go func (user string,id string) {
			var data m.Item
			getDb().Model(m.Item{}).Where("id = ?",id).Preload("Store").First(&data)

			itemCh <- data
		}(user,id)

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