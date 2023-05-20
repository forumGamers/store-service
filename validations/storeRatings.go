package validations

import "errors"

func CheckRates(rate string) (int, error) {
	switch rate {
	case "1":
		return 1, nil
	case "2":
		return 2, nil
	case "3":
		return 3, nil
	case "4":
		return 4, nil
	case "5":
		return 5, nil
	default:
		return 0, errors.New("rate is must between 1 - 5")
	}
}