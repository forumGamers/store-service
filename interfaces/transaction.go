package interfaces

type SuccessTransaction struct {
	ExpForUser    int
	Discount      int
	TransactionId int
	TotalPayment  int
	Voucher       bool
}