package services

import (
	"errors"
	"math"

	l "github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/jinzhu/gorm"
)

func CountFee(value int) int {

	if value < 1000000 {
		return int(math.Ceil(0.05 * float64(value)))
	}

	return int(math.Ceil(0.1 * float64(value)))
}

func TransactionExpForStore(value int,voucher *m.Voucher) int {
	return int(math.Ceil(0.01 * float64(value) + float64( VoucherPointForStore(voucher))))
}

func VoucherPointForStore(voucher *m.Voucher) int{
	return voucher.PointForStore
}

func VoucherPointForUser(voucher *m.Voucher) int {
	return voucher.PointForUser
}

func TransactionExpForUser(value int,voucher *m.Voucher) int {
	return int(math.Ceil(0.01 * float64(value) + float64( VoucherPointForUser(voucher))))
}

func AuthorizeTransaction(id int,transactionId int,errCh chan error,dataCh chan m.Transaction){
	var data m.Transaction

		if err := l.GetDb().Model(m.Transaction{}).Where("id = ?",transactionId).First(&data).Error ; err != nil {
			if err == gorm.ErrRecordNotFound {
				errCh <- errors.New("Data not found")
				dataCh <- m.Transaction{}
				return
			}else {
				errCh <- err
				dataCh <- m.Transaction{}
				return
			}
		}

		if data.User_id != uint(id) {
			errCh <- errors.New("Forbidden")
			return
		}

		errCh <- nil
		dataCh <- data
}