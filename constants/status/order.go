package constants

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusExpired   OrderStatus = "expired"
)

var statusToString = map[OrderStatus]string{
	OrderStatusPending:   "pending",
	OrderStatusPaid:      "paid",
	OrderStatusCancelled: "cancelled",
	OrderStatusExpired:   "expired",
}

var stringToStatus = map[string]OrderStatus{
	"pending":   OrderStatusPending,
	"paid":      OrderStatusPaid,
	"cancelled": OrderStatusCancelled,
	"expired":   OrderStatusExpired,
}

func (s OrderStatus) GetStatusString() string {
	if val, ok := statusToString[s]; ok {
		return val
	}
	return "unknown"
}
