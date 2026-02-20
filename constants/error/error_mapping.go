package error

import (
	errOrder "github.com/arthurhzna/ecommerce_be_test/constants/error/order"
	errPayment "github.com/arthurhzna/ecommerce_be_test/constants/error/payment"
	errProduct "github.com/arthurhzna/ecommerce_be_test/constants/error/product"
	errUser "github.com/arthurhzna/ecommerce_be_test/constants/error/user"
)

func ErrMapping(err error) bool {

	var (
		GeneralErrors = GeneralErrors
		UserErrors    = errUser.UserErrors
		ProductErrors = errProduct.ProductErrors
		OrderErrors   = errOrder.OrderErrors
		PaymentErrors = errPayment.PaymentErrors
	)
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, UserErrors...)
	allErrors = append(allErrors, ProductErrors...)
	allErrors = append(allErrors, OrderErrors...)
	allErrors = append(allErrors, PaymentErrors...)

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
