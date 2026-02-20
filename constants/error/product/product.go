package error

import "errors"

var (
	ErrProductNotFound         = errors.New("product not found")
	ErrInsufficientStock       = errors.New("insufficient stock")
	ErrStockCannotBeNegative   = errors.New("stock cannot be negative")
	ErrProductNameAlreadyExist = errors.New("product name already exist")
)

var ProductErrors = []error{
	ErrProductNotFound,
	ErrInsufficientStock,
	ErrStockCannotBeNegative,
	ErrProductNameAlreadyExist,
}
