package paymentgateway

import (
	"fmt"

	goservice "github.com/baozhenglab/go-sdk"
	"github.com/gin-gonic/gin"
)

func DemoPayment(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body VnpayRequest
		ctx.ShouldBindJSON(&body)
		vnpay := sc.MustGet(KeyVnPay).(PaymentService)
		res, err := vnpay.UrlPayment(body)
		fmt.Println(err)
		ctx.JSON(200, res)
	}
}

func CallbackDemoPayment(sc goservice.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body VnpayRequestCallback
		ctx.ShouldBind(&body)
		vnpay := sc.MustGet(KeyVnPay).(PaymentService)
		data, err := vnpay.Callback(body)
		fmt.Println(data, err)
		ctx.JSON(200, map[string]interface{}{"RspCode": "00"})
	}
}
