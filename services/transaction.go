package services

import (
	"math"

	m "github.com/forumGamers/store-service/models"
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