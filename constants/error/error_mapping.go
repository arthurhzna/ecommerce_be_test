package error

import (
	errUser "github.com/arthurhzna/ecommerce_be_test/constants/error/user"
)

func ErrMapping(err error) bool {

	var (
		GeneralErrors = GeneralErrors
		UserErrors    = errUser.UserErrors
	)
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, UserErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
