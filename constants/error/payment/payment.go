package error

import "errors"

var (
	ErrPaymentNotFound    = errors.New("payment not found")
	ErrPaymentAlreadyPaid = errors.New("payment already completed")
	ErrPaymentFailed      = errors.New("payment failed")
)

var PaymentErrors = []error{
	ErrPaymentNotFound,
	ErrPaymentAlreadyPaid,
	ErrPaymentFailed,
}
