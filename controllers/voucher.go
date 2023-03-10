package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	h "github.com/forumGamers/store-service/helper"
	m "github.com/forumGamers/store-service/models"
	s "github.com/forumGamers/store-service/services"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func AddVoucher(c *gin.Context){
	name,discount,cashback,period,stock := c.PostForm("name"),c.PostForm("discount"),c.PostForm("cashback"),c.PostForm("periode"),c.PostForm("stock")
	
	id,storeId := c.Request.Header.Get("id"),c.Request.Header.Get("storeId")

	checkCh := make(chan error)

	Id,r := strconv.ParseInt(id,10,64)

	if r != nil {
		panic("Forbidden")
	}

	sId,er := strconv.ParseInt(storeId,10,64)

	if er != nil {
		panic("Forbidden")
	}

	go func(id int,storeId int) {
		var data m.Store
		if err := getDb().Model(m.Voucher{}).Where("id = ? and owner_id = ?",storeId,id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				checkCh <- errors.New("You have not store yet")
				return
			}else {
				checkCh <- err
				return
			}
		}

		if data.Owner_id != id {
			checkCh <- errors.New("Forbidden")
			return
		}

		checkCh <- nil
	}(int(Id),int(sId))

	if err := <- checkCh ; err != nil {
		panic(err.Error())
	}

	errCh := make(chan error)

	go func(name string,discount string,cashback string,period string,storeId string,stock string){
		var data m.Voucher

		disc,errDisc := strconv.ParseInt(discount,10,64)
		
		if errDisc != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		cb,errCb := strconv.ParseInt(cashback,10,64)

		if errCb != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		p,errP := strconv.ParseInt(period,10,64)

		if errP != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		id,errId := strconv.ParseInt(storeId,10,64)

		if errId != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		stck,errStck := strconv.ParseInt(stock,10,64)

		if errStck != nil {
			errCh <- errors.New("Invalid data")
			return
		}

		data.Name = name
		data.Discount = int(disc)
		data.Cashback = int(cb)
		data.Period = int(p)
		data.Store_id = int(id)
		data.Stock = int(stck)
		data.PointForStore = s.ExpForStore(int(disc),int(cb),int(p))
		data.PointForUser = s.ExpForUser(int(disc),int(cb),int(p))

		if err := getDb().Model(m.Voucher{}).Create(&data).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(name,discount,cashback,period,storeId,stock)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}

func GetAllVoucher(c *gin.Context){
	name,store,page,limit := c.Query("name"),c.Query("store"),c.Query("page"),c.Query("limit")

	errCh := make(chan error)
	dataCh := make(chan []m.Voucher)

	go func(name string,store string,page string,limit string){
		var data []m.Voucher
		var args []interface{}
		var pg int
		var lmt int
		query := ""
		if name != "" {
			r := regexp.MustCompile(`\W`)
			res := r.ReplaceAllString(name,"")
			query = h.QueryBuild(query,"name ILIKE ?")
			args = append(args, "%"+res+"%")
		}

		if store != "" {
			if _,err := strconv.ParseInt(store,10,64) ; err != nil {
				errCh <- errors.New("Invalid params")
				dataCh <-nil
				return
			}
		}

		if page != "" {
			p,err := strconv.ParseInt(page,10,64) 
			if err != nil {
				pg = 1
			}else {
				pg = int(p)
			}
		}else {
			pg = 1
		}

		if limit != "" {
			l,err := strconv.ParseInt(limit,10,64)
			if err != nil {
				lmt = 10
			}else {
				lmt = int(l)
			}
		}else {
			lmt = 10
		}
		query = h.QueryBuild(query,"store_id = ?")
		args = append(args, store)
		if err := getDb().Model(m.Voucher{}).Where(query,args...).Offset((pg - 1) * lmt).Limit(lmt).Find(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- nil
				return
			}else {
				errCh <- err
				dataCh <- nil
				return
			}
		}

		errCh <- nil
		dataCh <- data
	}(name,store,page,limit)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh

	c.JSON(http.StatusOK,data)
}

func GetVoucherById(c *gin.Context){
	id := c.Param("id")

	errCh := make(chan error)
	dataCh := make(chan m.Voucher)

	go func(id string){
		var data m.Voucher
		if err := getDb().Model(m.Voucher{}).Where("id = ?",id).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- m.Voucher{}
				return
			}else {
				errCh <- err
				dataCh <- m.Voucher{}
				return
			}
		}

		dataCh <- data
		errCh <- nil
	}(id)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	data := <- dataCh

	c.JSON(http.StatusOK,data)
}

func DeleteVoucher(c *gin.Context){
	storeId := c.Request.Header.Get("store")
	id := c.Request.Header.Get("id")
	voucherId := c.Param("id")

	store,r := strconv.ParseInt(storeId,10,64)

	if r != nil {
		panic("Forbidden")
	}

	Id , er := strconv.ParseInt(id,10,64)

	if er != nil {
		panic("Forbidden")
	}

	errCh := make(chan error)

	go func (id int,voucher string,storeId int)  {
		var data m.Voucher
		if err := getDb().Model(m.Voucher{}).Where("id = ?",voucher).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				return
			}else {
				errCh <- err
				return
			}
		}

		if data.Store_id != storeId || data.Store.Owner_id != id {
			errCh <- errors.New("Forbidden")
			return
		}

		if err := getDb().Model(m.Voucher{}).Delete(m.Voucher{},voucher).Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}(int(Id),voucherId,int(store))

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK,gin.H{"message":"success"})
}
