package paymentgateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
)

type MomoRequestCreateOrder struct {
	RequestID string `json:"requestId"`
	OrderID   string `json:"orderId"`
	Amount    int32  `json:"amount"`
	OrderInfo string `json:"orderInfo"`
	ReturnURL string `json:"returnUrl"`
	ExtraData string `json:"extraData"`
}

type MomoResponseCreateOrder struct {
	RequestID        string `json:"requestId"`
	ErrorCode        int    `json:"errorCode"`
	Message          string `json:"message"`
	LocalMessage     string `json:"localMessage"`
	PayUrl           string `json:"payUrl"`
	QRCodeURL        string `json:"qrCodeUrl"`
	DeepLink         string `json:"deeplink"`
	DeepLinkWebInApp string `json:"deeplinkWebInApp"`
	Signature        string `json:"signature"`
}

type MomoRequestCallback struct {
	PartnerCode  string `json:"partnerCode" form:"partnerCode"`
	AccessKey    string `json:"accessKey" form:"accessKey"`
	RequestID    string `json:"requestId" form:"requestId"`
	Amount       int32  `json:"amount" form:"amount"`
	OrderID      string `json:"orderId" form:"orderId"`
	OrderInfo    string `json:"orderInfo" form:"orderInfo"`
	OrderType    string `json:"orderType" form:"orderType"`
	TransID      string `json:"transId" form:"transId"`
	ErrorCode    int    `json:"errorCode" form:"errorCode"`
	Message      string `json:"message" form:"message"`
	LocalMessage string `json:"localMessage" form:"localMessage"`
	PayType      string `json:"payType" form:"payType"`
	ResponseTime string `json:"responseTime" form:"responseTime"`
	ExtraData    string `json:"extraData" form:"extraData"`
	Signature    string `json:"signature" form:"signature"`
}

type momoService struct {
	partnerCode string
	accessKey   string
	secretKey   string
	baseURL     string
	notifyURL   string
}

const (
	endpointAPIMomo = "gw_payment/transactionProcessor"
)

func (m2 *momoService) Name() string {
	return KeyMomo
}

func (m2 *momoService) GetPrefix() string {
	return KeyMomo
}

func (m2 *momoService) Get() interface{} {
	return m2
}

func (m2 *momoService) InitFlags() {
	prefix := fmt.Sprintf("%s-", m2.Name())
	flag.StringVar(&m2.partnerCode, prefix+"partner-code", "", "Partnercode for momo app")
	flag.StringVar(&m2.accessKey, prefix+"accesskey", "", "AccessKey for momo payment app")
	flag.StringVar(&m2.secretKey, prefix+"secretkey", "", "secretkey for momo")
	flag.StringVar(&m2.notifyURL, prefix+"notify-url", "", "notify url for momo")
	flag.StringVar(&m2.baseURL, prefix+"base-url", "https://test-payment.momo.vn", "base url for payment momo")
}

func (m2 *momoService) UrlPayment(data interface{}) (*ResponseUrl, error) {
	request, ok := data.(MomoRequestCreateOrder)
	if !ok {
		return nil, fmt.Errorf("data input must type of MomoRequestCreateOrder,have type %s", getType(data))
	}
	res, err := m2.requestPayment(request)
	if err != nil {
		return nil, err
	}
	checkres := fmt.Sprintf("requestId=%s&orderId=%s&message=%s&localMessage=%s&payUrl=%s&errorCode=%d&requestType=captureMoMoWallet",
		res.RequestID, request.OrderID, res.Message, res.LocalMessage, res.PayUrl, res.ErrorCode)
	if GenerateHashMacSha256(m2.secretKey, checkres) != res.Signature {
		return nil, errors.New("signature is invalid")
	}
	return &ResponseUrl{
		URL:     res.PayUrl,
		Message: "Get sucessfully",
		Additional: map[string]string{
			"qrCodeUrl":        res.QRCodeURL,
			"deeplink":         res.DeepLink,
			"deeplinkWebInApp": res.DeepLinkWebInApp,
		},
	}, nil
}

func (m2 *momoService) Callback(data interface{}) (*ResponseCallback, error) {
	request, ok := data.(MomoRequestCallback)
	if !ok {
		return nil, fmt.Errorf("data input  must is typeof MomoRequestCallback,have %s", getType(data))
	}
	if m2.checkSignatureCallback(request) == false {
		return nil, errors.New("Signature is invalid")
	}
	return &ResponseCallback{
		Amount:      request.Amount,
		TransID:     request.OrderID,
		ExtractData: request.ExtraData,
	}, nil
}

func (m2 *momoService) checkSignatureCallback(data MomoRequestCallback) bool {
	hashinput := fmt.Sprintf("partnerCode=%s&accessKey=%s&requestId=%s&amount=%d&orderId=%s&orderInfo=%s&orderType=%s&transId=%s&message=%s&localMessage=%s&responseTime=%s&errorCode=%d&payType=%s&extraData=%s",
		data.PartnerCode, data.AccessKey, data.RequestID, data.Amount, data.OrderID, data.OrderInfo, data.OrderType, data.TransID, data.Message,
		data.LocalMessage, data.ResponseTime, data.ErrorCode, data.PayType, data.ExtraData)
	return GenerateHashMacSha256(m2.secretKey, hashinput) == data.Signature
}

func (m2 *momoService) createSignature(data MomoRequestCreateOrder) string {
	hashinput := fmt.Sprintf("partnerCode=%s&accessKey=%s&requestId=%s&amount=%d&orderId=%s&orderInfo=%s&returnUrl=%s&notifyUrl=%s&extraData=%s",
		m2.partnerCode, m2.accessKey, data.RequestID, data.Amount, data.OrderID, data.OrderInfo, data.ReturnURL, m2.notifyURL, data.ExtraData)
	return GenerateHashMacSha256(m2.secretKey, hashinput)
}

func (m2 *momoService) requestPayment(data MomoRequestCreateOrder) (*MomoResponseCreateOrder, error) {
	jsonVal, _ := json.Marshal(data)
	var body map[string]interface{}
	json.Unmarshal(jsonVal, &body)
	body["amount"] = fmt.Sprintf("%.0f", body["amount"])
	body["partnerCode"] = m2.partnerCode
	body["accessKey"] = m2.accessKey
	body["notifyUrl"] = m2.notifyURL
	body["requestType"] = "captureMoMoWallet"
	body["signature"] = m2.createSignature(data)
	dataReq, _ := json.Marshal(body)
	url := fmt.Sprintf("%s/%s", m2.baseURL, endpointAPIMomo)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dataReq))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request create order momo return status %d,want status is 200", res.StatusCode)
	}
	var response MomoResponseCreateOrder
	json.NewDecoder(res.Body).Decode(&response)
	if response.ErrorCode != 0 {
		return nil, fmt.Errorf("request create order momo return ErrorCode %d,want ErrorCode is 0", response.ErrorCode)
	}
	return &response, nil
}
