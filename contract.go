package paymentgateway

type ResponseUrl struct {
	URL        string      `json:"url"`
	Message    string      `json:"message"`
	Additional interface{} `json:"additional"`
}

type PaymentService interface {
	UrlPayment(data interface{}) (*ResponseUrl, error)
}

const (
	KeyVnPay   = "vnpay"
	KeyZaloPay = "zalopay"
)
