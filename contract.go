package paymentgateway

type ResponseUrl struct {
	URL        string      `json:"url"`
	Message    string      `json:"message"`
	Additional interface{} `json:"additional"`
}

type ResponseCallback struct {
	Amount      int32
	TransID     string
	ExtractData string
}

type PaymentService interface {
	UrlPayment(data interface{}) (*ResponseUrl, error)
	Callback(data interface{}) (*ResponseCallback, error)
}

const (
	KeyVnPay   = "vnpay"
	KeyZaloPay = "zalopay"
)
