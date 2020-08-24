package paymentgateway

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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

type VnpayRequestCallback struct {
	TmnCode    string `form:"vnp_TmnCode" json:"vnp_TmnCode"`
	Amount     int32  `form:"vnp_Amount" json:"vnp_Amount"`
	Bankcode   string `form:"vnp_BankCode" json:"vnp_BankCode"`
	BankTranNo string `form:"vnp_BankTranNo" json:"vnp_BankTranNo"`
	CardType   string `form:"vnp_CardType" json:"vnp_CardType"`
	PayDate    int64  `form:"vnp_PayDate" json:"vnp_PayDate"`
	// CurrCode      string `form:"vnp_CurrCode" json:"vnp_CurrCode"`
	OrderInfo     string `form:"vnp_OrderInfo" json:"vnp_OrderInfo"`
	TransactionNo int64  `form:"vnp_TransactionNo" json:"vnp_TransactionNo"`
	ResponseCode  string `form:"vnp_ResponseCode" json:"vnp_ResponseCode"`
	TxnRef        string `form:"vnp_TxnRef" json:"vnp_TxnRef"`
	SecureHash    string `form:"vnp_SecureHash" json:"vnp_SecureHash"`
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

func (vp *vnpayService) Callback(data interface{}) (*ResponseCallback, error) {
	request, ok := data.(VnpayRequestCallback)
	if !ok {
		return nil, errors.New("input is not type VnpayRequestCallback")
	}
	if vp.checkSumCallback(request) == false {
		return nil, errors.New("signature is error")
	}
	return &ResponseCallback{request.Amount, request.TxnRef, request.OrderInfo}, nil
}

func (vp *vnpayService) checkSumCallback(data VnpayRequestCallback) bool {
	jsonval, _ := json.Marshal(data)
	var req map[string]interface{}
	json.Unmarshal(jsonval, &req)
	delete(req, "vnp_SecureHash")
	signdata := vp.hashSecret + generateQueryBasic(req)
	return NewSHA256(signdata) == data.SecureHash
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
		encodeStr += "&vnp_SecureHash=" + NewSHA256(signdata)
	}
	return encodeStr
}
