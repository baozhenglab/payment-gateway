package paymentgateway

import (
	"flag"
	"fmt"
)

type zalopayService struct {
	appid         int
	bankcode      string
	appSecretKey1 string
	appSecretKey2 string
	version       string
	baseURL       string
}

func (*zalopayService) Name() string {
	return KeyZaloPay
}

func (*zalopayService) GetPrefix() string {
	return KeyZaloPay
}

func (zp *zalopayService) Get() interface{} {
	return zp
}

func (zp *zalopayService) InitFlags() {
	prefix := fmt.Sprintf("%s-", zp.Name())
	flag.IntVar(&zp.appid, prefix+"app-id", 0, "App id of zalopay merchant")
	flag.StringVar(&zp.appSecretKey1, prefix+"app-secret-key-1", "", "Secret key 1 for zalopay")
	flag.StringVar(&zp.appSecretKey2, prefix+"app-secret-key-2", "", "Secret key 2 for zalopay")
	flag.StringVar(&zp.baseURL, prefix+"base-url", "https://sandbox.zalopay.com.vn", "base url request zalopay")
	flag.StringVar(&zp.version, prefix+"version", "v001", "Version default for zalopay")
	flag.StringVar(&zp.bankcode, prefix+"bankcode", "zalopayapp", "Bank code default for zalopay")
}

func (zp *zalopayService) UrlPayment(data interface{}) (*ResponseUrl, error) {
	return nil, nil
}
