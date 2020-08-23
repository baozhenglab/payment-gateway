package paymentgateway

import (
	"net/http"
	"time"

	goservice "github.com/baozhenglab/go-sdk"
)

var client = http.Client{
	Timeout: time.Duration(10 * time.Second),
}

func NewServicePayment(key string) goservice.PrefixConfigure {
	switch key {
	case KeyVnPay:
		return &vnpayService{}
	default:
		panic("not found")
	}
}
