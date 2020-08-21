package paymentgateway

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
)

//GetMD5Hash excute md5 string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//create secret key with sha256 algothirm
func GenerateHashMacSha256(secret string, data string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}

// NewSHA256 ...
func NewSHA256(data string) string {
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash[:])
}

func indexOfString(arr []string, str string) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == str {
			return i
		}
	}
	return -1
}

func stringifyQuery(data interface{}) string {
	jsonEn, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	var jsonStr map[string]interface{}
	json.Unmarshal(jsonEn, &jsonStr)
	params := url.Values{}
	for key, value := range jsonStr {
		params.Add(key, fmt.Sprintf("%s", value))
	}
	return params.Encode()
}

func generateQueryBasic(data interface{}) string {
	jsonEn, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	var jsonStr map[string]interface{}
	json.Unmarshal(jsonEn, &jsonStr)
	keys := make([]string, 0)
	str := ""
	for key := range jsonStr {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		str += (key + "=" + fmt.Sprintf("%s", jsonStr[key]) + "&")
	}
	return str[:len(str)-1]
}
