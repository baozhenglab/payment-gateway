package paymentgateway

import (
	goservice "github.com/baozhenglab/go-sdk"
)

func NewServicePayment(key string) goservice.PrefixConfigure {
	switch key {
	case KeyVnPay:
		return &vnpayService{}
	default:
		panic("not found")
	}
}
