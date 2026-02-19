package constants

type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "paid"
	PaymentStatusFailed  PaymentStatus = "failed"
)

var paymentStatusToString = map[PaymentStatus]string{
	PaymentStatusPending: "pending",
	PaymentStatusPaid:    "paid",
	PaymentStatusFailed:  "failed",
}

var stringToPaymentStatus = map[string]PaymentStatus{
	"pending": PaymentStatusPending,
	"paid":    PaymentStatusPaid,
	"failed":  PaymentStatusFailed,
}

func (s PaymentStatus) GetStatusString() string {
	if val, ok := paymentStatusToString[s]; ok {
		return val
	}
	return "unknown"
}
