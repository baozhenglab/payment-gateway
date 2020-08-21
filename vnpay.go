package paymentgateway

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type VnpayRequest struct {
	Amount    uint32 `json:"vnp_Amount"`
	BankCode  string `json:"vnp_BankCode"`
	IpAddr    string `json:"vnp_IpAddr"`
	OrderInfo string `json:"vnp_OrderInfo"`
	OrderType string `json:"vnp_OrderType"`
	ReturnUrl string `json:"vnp_ReturnUrl"`
	TxnRef    string `json:"vnp_TxnRef"`
	Locale    string `json:"vnp_Locale"`
}

var client = http.Client{
	Timeout: time.Duration(10 * time.Second),
}

type VnpayRequestFull struct {
	VnpayRequest
	Amount     string `json:"vnp_Amount"`
	Version    string `json:"vnp_Version"`
	Command    string `json:"vnp_Command"`
	TmnCode    string `json:"vnp_TmnCode"`
	CreateDate string `json:"vnp_CreateDate"`
	Currency   string `json:"vnp_CurrCode"`
}

type vnpayService struct {
	tmnCode    string
	version    string
	command    string
	currency   string
	local      string
	returnUrl  string
	hashType   string
	hashSecret string
	baseURL    string
}

func (*vnpayService) Name() string {
	return KeyVnPay
}

func (*vnpayService) GetPrefix() string {
	return KeyVnPay
}

func (vp *vnpayService) Get() interface{} {
	return vp
}

func (vp *vnpayService) InitFlags() {
	prefix := fmt.Sprintf("%s-", vp.Name())
	flag.StringVar(&vp.version, prefix+"version", "2", "Version for api payment of vnpay")
	flag.StringVar(&vp.tmnCode, prefix+"tmn-code", "", "TmnCode for vnpay")
	flag.StringVar(&vp.hashSecret, prefix+"hash-secret", "", "hash secret for vnpay")
	flag.StringVar(&vp.command, prefix+"command", "pay", "comannd pay for vnpay, default pay")
	flag.StringVar(&vp.baseURL, prefix+"base-url", "", "base url for vnpay")
	flag.StringVar(&vp.currency, prefix+"curency", "VND", "currency for payment vnpay, default is VND")
	flag.StringVar(&vp.local, prefix+"local", "en", "Langaue show for payment vnpay, default is en")
	flag.StringVar(&vp.hashType, prefix+"hash-type", "sha256", "hash type for vnpay default is sha256")
	flag.StringVar(&vp.returnUrl, prefix+"return-url", "", "return url default for vnpay")
}

func (vp *vnpayService) UrlPayment(data interface{}) (*ResponseUrl, error) {
	request, ok := data.(VnpayRequest)
	if !ok {
		return nil, errors.New("Error is not type VnpayRequest")
	}
	params, err := vp.getRequestParam(request)
	if err != nil {
		return nil, err
	}
	encodeStr := vp.genrateQueryString(params)
	return &ResponseUrl{
		Message: "Get sucessfully",
		URL:     vp.baseURL + "?" + encodeStr,
	}, nil
}

func (vp *vnpayService) getRequestParam(request VnpayRequest) (VnpayRequestFull, error) {
	body := VnpayRequestFull{}
	body.Amount = fmt.Sprintf("%d", request.Amount*100)
	body.BankCode = request.BankCode
	body.IpAddr = request.IpAddr
	body.OrderInfo = request.OrderInfo
	body.OrderType = request.OrderType
	body.TxnRef = request.TxnRef
	body.Currency = vp.currency
	if request.ReturnUrl == "" {
		body.ReturnUrl = vp.returnUrl
	} else {
		body.ReturnUrl = request.ReturnUrl
	}
	if request.Locale == "" {
		body.Locale = vp.local
	} else {
		body.Locale = request.Locale
	}
	body.Version = vp.version
	body.Command = vp.command
	body.TmnCode = vp.tmnCode
	if indexOfString([]string{"sha256", "md5"}, vp.hashType) < 0 {
		return body, fmt.Errorf("hash type %s in not belong to sha256 or md5", vp.hashType)
	}

	body.CreateDate = time.Now().UTC().Format("20060102150405")
	return body, nil

}

func (vp *vnpayService) genrateQueryString(data VnpayRequestFull) string {
	encodeStr := stringifyQuery(data) + "&vnp_SecureHashType=" + strings.ToUpper(vp.hashType)
	signdata := vp.hashSecret + generateQueryBasic(data)
	switch vp.hashType {
	case "md5":
		encodeStr += "&vnp_SecureHash=" + GetMD5Hash(signdata)
	default:
		fmt.Println(signdata)
		encodeStr += "&vnp_SecureHash=" + NewSHA256(signdata)
	}
	return encodeStr
}
