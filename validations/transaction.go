package validations

import "strconv"

func CheckDataTransaction(userId string, amount string) (int, int, error) {
	var id int
	var amounts int

	if Id,err :=  strconv.ParseInt(userId,10,64) ; err != nil {
		return id,amounts,err
	}else {
		id = int(Id)
	}

	if a,err := strconv.ParseInt(amount,10,64) ; err != nil {
		return id,amounts,err
	}else {
		amounts = int(a)
	}

	return id,amounts,nil
}