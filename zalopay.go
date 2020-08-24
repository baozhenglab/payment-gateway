package paymentgateway

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type zalopayService struct {
	appid         int
	bankcode      string
	appSecretKey1 string
	appSecretKey2 string
	version       string
	baseURL       string
}

type ResponseCreateOrder struct {
	ReturnCode    int    `json:"returncode"`
	ReturnMessage string `json:"returnmessage"`
	OrderURL      string `json:"orderurl"`
	ZptransToken  string `json:"zptranstoken"`
}

type RequestCreateOrder struct {
	AppUser     string `json:"appuser"`
	Amount      int32  `json:"amount"`
	ApptransID  string `json:"apptransid"`
	EmbedData   string `json:"embeddata"`
	Item        string `json:"item"`
	Description string `json:"description"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	BankCode    string `json:"bankcode"`
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
	request, ok := data.(RequestCreateOrder)
	if !ok {
		return nil, errors.New("data input is type of RequestCreateOrder")
	}
	res, err := zp.requestPayment(request)
	if err != nil {
		return nil, err
	}

	return &ResponseUrl{
		URL:     res.OrderURL,
		Message: "Request successfully",
		Additional: map[string]string{
			"zptranstoken": res.ZptransToken,
		},
	}, nil
}

func (zp *zalopayService) requestPayment(request RequestCreateOrder) (*ResponseCreateOrder, error) {
	body := zp.generateStringCreateOrder(request)
	url := fmt.Sprintf("%s/%s/%s", zp.baseURL, zp.version, "tpe/createorder")

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Can not get request with status code 200")
	}
	var orderUrl ResponseCreateOrder
	json.NewDecoder(res.Body).Decode(&orderUrl)
	if orderUrl.ReturnCode != 1 {
		return nil, fmt.Errorf("request return code %d with message is %s", orderUrl.ReturnCode, orderUrl.ReturnMessage)
	}
	return &orderUrl, nil
}

func (zp *zalopayService) generateStringCreateOrder(request RequestCreateOrder) string {
	jsonValue, _ := json.Marshal(request)
	var mapStr map[string]interface{}
	json.Unmarshal(jsonValue, &mapStr)
	mapStr["apptime"] = fmt.Sprintf("%d", time.Now().Unix())
	mapStr["appid"] = fmt.Sprintf("%d", zp.appid)
	if bankcode, ok := mapStr["bankcode"]; !ok || fmt.Sprintf("%s", bankcode) == "" {
		mapStr["bankcode"] = zp.bankcode
	}
	mapStr["amount"] = fmt.Sprintf("%d", request.Amount)
	mapStr["apptransid"] = convertTransIDZalo(request.ApptransID)
	hmacinput := fmt.Sprintf("%d|%s|%s|%d|%s|%s|%s", zp.appid, mapStr["apptransid"], mapStr["appuser"],
		request.Amount, mapStr["apptime"], request.EmbedData, request.Item)
	mapStr["mac"] = GenerateHashMacSha256(zp.appSecretKey1, hmacinput)
	body := ""
	for key, value := range mapStr {
		body += fmt.Sprintf("%s=%s&", key, value)
	}
	return body[:len(body)-1]
}

func convertTransIDZalo(apptransid string) string {
	var splitStr = strings.Split(apptransid, "_")
	timeUTC := time.Now().UTC().Format("20060102150405")
	if len(splitStr) > 1 {
		r := regexp.MustCompile("^[0-9]{6,6}$")
		if !r.MatchString(splitStr[0]) {
			splitStr[0] = timeUTC[2:8]
		}
	} else {
		splitStr[0] = timeUTC[2:8] + "_" + splitStr[0]
	}
	return strings.Join(splitStr, "_")
}
