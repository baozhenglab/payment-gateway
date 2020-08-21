package paymentgateway

import (
	goservice "chat-backend/external/go-sdk"
)

func NewServicePayment(key string) goservice.PrefixConfigure {
	switch key {
	case KeyVnPay:
		return &vnpayService{}
	default:
		panic("not found")
	}
}
