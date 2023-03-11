package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	m "github.com/forumGamers/store-service/models"
	s "github.com/forumGamers/store-service/services"
	validate "github.com/forumGamers/store-service/validations"
	"github.com/gin-gonic/gin"
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
			fmt.Println(item.Price)

			tx := getDb().Begin()
			var amounts int

			transaction.Item_id = item.ID

			transaction.Store_id = store.ID

			if userId,amounts,err := validate.CheckDataTransaction(id,amount) ; err != nil {
				errCh <- err
				return
			}else {
				transaction.User_id = uint(userId)
				transaction.Amount = amounts

				if err := getDb().Model(m.Store{}).Where("id = ?",store.ID).Update("exp",(store.Exp + s.TransactionExpForStore(transaction.Value,&v))).Error ; err != nil {
					errCh <- err
					tx.Rollback()
					return
				}
			}

			if err := <- voucherCheck ; err != nil && err.Error() != "skip" {
				v = <- voucherCh

				test := s.VoucherCheck(v,store.ID)

				if !test {
					errCh <- errors.New("voucher is not registered")
					return
				}

				transaction.Value = transaction.Amount * item.Price - v.Discount
				fmt.Println(transaction.Value,"disc")
			}else {
				transaction.Value = transaction.Amount * item.Price
				fmt.Println(transaction.Value)
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

			if err := getDb().Model(m.Item{}).Where("id = ?",item.ID).Update(map[string]interface{}{"stock":item.Stock - amounts,"sold":item.Sold + 1}).Error ; err != nil {
				errCh <- err
				tx.Rollback()
				return
			}

			errCh <- nil
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

	c.JSON(http.StatusCreated,gin.H{"message" : "success"})

}