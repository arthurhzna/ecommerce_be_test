package error

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderAlreadyPaid  = errors.New("order already paid")
	ErrOrderNotPending   = errors.New("order is not pending")
	ErrOrderUnauthorized = errors.New("unauthorized to access this order")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrOrderAlreadyPaid,
	ErrOrderNotPending,
	ErrOrderUnauthorized,
}
