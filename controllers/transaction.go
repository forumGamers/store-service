package controllers

import (
	"errors"
	"net/http"
	"strconv"

	i "github.com/forumGamers/store-service/interfaces"
	m "github.com/forumGamers/store-service/models"
	s "github.com/forumGamers/store-service/services"
	validate "github.com/forumGamers/store-service/validations"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func CreateTransaction(c *gin.Context){
	itemId := c.Param("itemId")
	storeId := c.Param("storeId")
	id := c.Request.Header.Get("id")

	payment_method,amount,message :=
	c.PostForm("payment_method"),
	c.PostForm("amount"),
	c.PostForm("message")

	voucherId := c.Query("voucherId")
	var voucher int

	if amount == "" || payment_method == "" {
		panic("Invalid data")
	}

	if voucherId == "" {
		voucher = 0
	}else {
		v,err := strconv.ParseInt(voucherId,10,64)

		if err != nil {
			panic("Invalid params")
		}

		voucher = int(v)
	}

	itemCh := make(chan m.Item)
	storeCh := make(chan m.Store)
	storeCheck := make(chan error)
	checkCh := make(chan error)
	errCh := make(chan error)
	voucherCh := make(chan m.Voucher)
	voucherCheck := make(chan error)
	responseCh := make(chan i.SuccessTransaction)


	go s.CheckAvailablity(itemId,amount,checkCh,itemCh)

	go s.GetStore(storeId,storeCh,storeCheck)

	go s.GetVoucher(voucher,voucherCh,voucherCheck)

	go func(
		itemId string,
		storeId string,
		id string,
		payment_method string,
		amount string,
		message string,
		voucherId int,
		){
			var transaction m.Transaction
			var v m.Voucher
			var response i.SuccessTransaction

			if err := <- storeCheck ;  err != nil {
				errCh <- err
				return
			}

			store := <- storeCh

			if err := <- checkCh ; err != nil {
				errCh <- err
				return
			}

			item := <- itemCh

			tx := getDb().Begin()

			transaction.Item_id = item.ID

			transaction.Store_id = store.ID

			if userId,amounts,err := validate.CheckDataTransaction(id,amount) ; err != nil {
				errCh <- err
				return
			}else {
				transaction.User_id = uint(userId)
				transaction.Amount = amounts
			}

			if err := <- voucherCheck ; err != nil || err.Error() != "skip" {

				transaction.Value = transaction.Amount * item.Price
				response.Discount = 0
				response.TotalPayment = transaction.Amount * item.Price
				response.Voucher = false
			}else {
				v = <- voucherCh
				test := s.VoucherCheck(v,store.ID)

				if !test {
					errCh <- errors.New("voucher is not registered")
					return
				}

				transaction.Value = transaction.Amount * item.Price 
				response.Discount = v.Discount
				response.TotalPayment = transaction.Amount * item.Price - v.Discount
				response.Voucher = true

				if err := getDb().Model(m.Voucher{}).Where("id = ?",v.ID).Update("stock",v.Stock - 1).Error ; err != nil {
					errCh <- err
					tx.Rollback()
					return
				}
			}

			transaction.Payment_method = payment_method
			transaction.Fee = s.CountFee(transaction.Value)
			transaction.MessageForSeller = message
			transaction.Description = "Purchasing " + amount + " " + item.Name

			if err := getDb().Model(m.Transaction{}).Create(&transaction).Error ; err != nil {
				errCh <- err
				tx.Rollback()
				return
			}

			response.TransactionId = int(transaction.ID)
			response.ExpForUser = s.TransactionExpForUser(transaction.Value,&v)

			if err := getDb().Model(m.Store{}).Where("id = ?",store.ID).Update("exp",store.Exp + s.TransactionExpForStore(transaction.Value,&v)).Error ; err != nil {
				errCh <- err
				tx.Rollback()
				return
			}

			if err := getDb().Model(m.Item{}).Where("id = ?",item.ID).Update(map[string]interface{}{"stock":item.Stock - transaction.Amount,"sold":item.Sold + transaction.Amount}).Error ; err != nil {
				errCh <- err
				tx.Rollback()
				return
			}

			errCh <- nil
			responseCh <- response
			tx.Commit()
	}(
		itemId,
		storeId,
		id,
		payment_method,
		amount,
		message,
		voucher,
	)

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	response := <- responseCh

	c.JSON(http.StatusCreated,response)
}

func EndTransaction(c *gin.Context){
	transactionId := c.Param("transactionId")
	id := c.Request.Header.Get("id")

	Id,tId,err := validate.CheckEndTransactionData(id,transactionId)

	if err != nil {
		panic(err.Error())
	}

	checkCh := make (chan error)
	dataCh := make(chan m.Transaction)

	go func(id int,transactionId int) {
		var data m.Transaction

		if err := getDb().Model(m.Transaction{}).Where("id = ?",transactionId).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				checkCh <- errors.New("Data not found")
				dataCh <- m.Transaction{}
				return
			}else {
				checkCh <- err
				dataCh <- m.Transaction{}
				return
			}
		}

		if data.User_id != uint(id) {
			checkCh <- errors.New("Forbidden")
			return
		}

		checkCh <- nil
		dataCh <- data
	}(Id,tId)

	errCh := make(chan error)

	go func ()  {
		if err := <- checkCh ; err != nil {
			errCh <- err
			return
		}

		transaction := <- dataCh

		if err := getDb().Model(m.Transaction{}).Where("id = ?",transaction.ID).Update("status","Finish").Error ; err != nil {
			errCh <- err
			return
		}

		errCh <- nil
	}()

	if err := <- errCh ; err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusCreated,gin.H{"message":"success"})
}