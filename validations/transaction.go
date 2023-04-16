package validations

import (
	"errors"
	"strconv"
)

func CheckDataTransaction(userId int, amount string) (int, int, error) {
	var id int
	var amounts int

	if a,err := strconv.ParseInt(amount,10,64) ; err != nil {
		return id,amounts,errors.New("Invalid data")
	}else {
		amounts = int(a)
	}

	return id,amounts,nil
}

func CheckEndTransactionData(id int,transactionId string) (int,int,error) {
	var Id int
	var tId int

	if t,err := strconv.ParseInt(transactionId,10,64) ; err != nil {
		return Id,tId,errors.New("Invalid data")
	}else {
		tId = int(t)
	}

	return Id,tId,nil
}