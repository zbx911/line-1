package line

import (
	"encoding/base64"
	"fmt"
	"github.com/line-api/model/go/model"
	"strings"
	"time"
)

type TokenManager struct {
	AuthKey string
}

func parseAuthKey(key string) (string, string) {
	splited := strings.Split(key, ":")
	return splited[0], splited[1]
}

func (cl *Client) GeneratePrimaryToken(authKey string) string {
	switch cl.ClientSetting.AppType {
	case model.ApplicationType_ANDROIDLITE:
		return GenerateLineLiteToken(authKey)
	case model.ApplicationType_ANDROID:
		return GenerateAndroidToken(authKey)
	case model.ApplicationType_IOS:
		return GenerateIOSToken(authKey)
	}
	return ""
}

func GenerateIOSToken(authKey string) string {
	mid, key := parseAuthKey(authKey)
	iat := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("iat: %v\n", time.Now().Unix()*60))) + "."
	keyEnc, _ := base64.StdEncoding.DecodeString(key)
	return mid + ":" + iat + "." + base64.StdEncoding.EncodeToString(SignHmacSha1(keyEnc, []byte(iat)))
}

func GenerateAndroidToken(authKey string) string {
	mid, key := parseAuthKey(authKey)
	iat := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("iat: %v\n", time.Now().UnixNano()/int64(time.Millisecond)))) + "."
	keyEnc, _ := base64.StdEncoding.DecodeString(key)
	return mid + ":" + iat + "." + base64.StdEncoding.EncodeToString(SignHmacSha1(keyEnc, []byte(iat)))
}

func GenerateLineLiteToken(authKey string) string {
	mid, key := parseAuthKey(authKey)
	keyEnc, _ := base64.StdEncoding.DecodeString(key)
	header := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("issuedTo: %v\niat: %v\n", mid, time.Now().UnixNano()/int64(time.Millisecond))))
	header2 := base64.StdEncoding.EncodeToString([]byte("type: YWT\nalg: HMAC_SHA1\n"))

	signature := base64.StdEncoding.EncodeToString(SignHmacSha1(keyEnc, []byte(fmt.Sprintf("%v.%v", header, header2))))
	wToken := fmt.Sprintf("%v.%v.%v", header, header2, signature)
	return mid + ":" + wToken
}
