package line

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/line-api/model/go/model"
	"strings"
	"time"
)

type TokenManager struct {
	AuthKey      string
	AccessToken  string
	RefreshToken string
	IsV3Token    bool
}

type V3TokenContent struct {
	JwtId                 string `json:"jti"`
	Audience              string `json:"aud"`
	IssuedAt              int64  `json:"iat"`
	ExpiredAt             int64  `json:"exp"`
	Scope                 string `json:"scp"`
	RefreshTokenId        string `json:"rtid"`
	RefreshTokenExpiredAt int64  `json:"rexp"`
	Version               string `json:"ver"`
	AuthorId              string `json:"aid"`
	LineSessionId         string `json:"lsid"`
	Did                   string `json:"did"`
	ClientType            string `json:"ctype"`
	ClientMode            string `json:"cmode"`
	ClientId              string `json:"cid"`
}

func (t *TokenManager) parseV3Token() (*V3TokenContent, error) {
	jsonData, _ := base64.StdEncoding.DecodeString(strings.Split(t.AccessToken, ".")[1] + "==")
	token := &V3TokenContent{}
	err := json.Unmarshal(jsonData, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (cl *Client) runTokenUpdaterBackGround() {
	go func() {
		for {
			token, err := cl.TokenManager.parseV3Token()
			if err != nil {
				cl.ClientSetting.Logger.Error().Err(err).Msgf("stopped token updater: failed to parse v3 token")
				return
			}
			if time.Unix(token.ExpiredAt, 0).Add(-time.Hour*24).Unix() >= time.Now().Unix() {
				time.Sleep(time.Hour * 1)
				continue
			}
			err = cl.RefreshV3AccessToken()
			if err != nil {
				cl.ClientSetting.Logger.Error().Err(err).Msgf("stopped token updater: failed to refresh v3 token")
				return
			}
			cl.ClientSetting.Logger.Info().Msgf("refreshed access token")
			time.Sleep(time.Hour * 12)
		}
	}()
}

func parseAuthKey(key string) (string, string) {
	splited := strings.Split(key, ":")
	return splited[0], splited[1]
}

func (cl *Client) GeneratePrimaryToken(authKey string) string {
	switch cl.ClientSetting.AppType {
	case model.ApplicationType_ANDROIDLITE:
		return generateLineLiteToken(authKey)
	case model.ApplicationType_ANDROID:
		return generateAndroidToken(authKey)
	case model.ApplicationType_IOS:
		return generateIOSToken(authKey)
	}
	return ""
}

func generateIOSToken(authKey string) string {
	mid, key := parseAuthKey(authKey)
	iat := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("iat: %v\n", time.Now().Unix()*60))) + "."
	keyEnc, _ := base64.StdEncoding.DecodeString(key)
	return mid + ":" + iat + "." + base64.StdEncoding.EncodeToString(SignHmacSha1(keyEnc, []byte(iat)))
}

func generateAndroidToken(authKey string) string {
	mid, key := parseAuthKey(authKey)
	iat := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("iat: %v\n", time.Now().UnixNano()/int64(time.Millisecond)))) + "."
	keyEnc, _ := base64.StdEncoding.DecodeString(key)
	return mid + ":" + iat + "." + base64.StdEncoding.EncodeToString(SignHmacSha1(keyEnc, []byte(iat)))
}

func generateLineLiteToken(authKey string) string {
	mid, key := parseAuthKey(authKey)
	keyEnc, _ := base64.StdEncoding.DecodeString(key)
	header := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("issuedTo: %v\niat: %v\n", mid, time.Now().UnixNano()/int64(time.Millisecond))))
	header2 := base64.StdEncoding.EncodeToString([]byte("type: YWT\nalg: HMAC_SHA1\n"))

	signature := base64.StdEncoding.EncodeToString(SignHmacSha1(keyEnc, []byte(fmt.Sprintf("%v.%v", header, header2))))
	wToken := fmt.Sprintf("%v.%v.%v", header, header2, signature)
	return mid + ":" + wToken
}
