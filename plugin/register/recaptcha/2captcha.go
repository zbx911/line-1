package recaptcha

import (
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	gSiteKey = "6Lfo_XYUAAAAAFQbdsuk6tETqnpKIg5gNxJy4xM0"
)

type TwoCaptcha struct {
	apiKey string
}

func NewTwoCaptcha(key string) *TwoCaptcha {
	return &TwoCaptcha{
		apiKey: key,
	}
}

func (c *TwoCaptcha) Solve(details *model.WebAuthDetails, client *http.Client) (string, error) {
	headers := map[string]string{
		"Connection":                "keep-alive",
		"Upgrade-Insecure-Requests": "1",
		"user-agent":                "Mozilla/5.0 (Linux; Android 8.0.0; Nexus 6P Build/OPR6.170623.019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.93 Mobile Safari/537.36",
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"authorization":             details.Token,
		"Accept-Encoding":           "gzip, deflate",
		"Accept-Language":           "ja-JP,ja;q=0.9,en-US;q=0.8,en;q=0.7",
		"Host":                      "2captcha.com",
	}
	params := map[string]string{
		"method":      "userrecaptcha",
		"pageurl":     details.BaseUrl,
		"soft_id":     "3029",
		"googlekey":   gSiteKey,
		"key":         c.apiKey,
		"header_acao": "1",
	}
	request, err := http.NewRequest("GET", "https://2captcha.com/in.php", nil)
	if err != nil {
		return "", err
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}
	rParam := request.URL.Query()
	for k, v := range params {
		rParam.Add(k, v)
	}
	request.URL.RawQuery = rParam.Encode()
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	if string(body)[:2] != "OK" {
		return "", xerrors.New("something went wrong on starting solve capture process" + string(body))
	}
	captchaId := string(body)[3:]
	for i := 0; i < 20; i++ {
		time.Sleep(time.Second * 5)
		request, err = http.NewRequest("GET", "https://2captcha.com/res.php?key="+c.apiKey+"&action=get&id="+captchaId, nil)
		if err != nil {
			return "", err
		}
		response, err = client.Do(request)
		if err != nil {
			if strings.Contains(err.Error(), "timeout awaiting") {
				continue
			}
			return "", err
		}
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		if string(body)[:2] == "OK" {
			return string(body)[3:], nil
		}
		if string(body) == "CAPCHA_NOT_READY" {
			continue
		}
		response.Body.Close()
	}
	return "", xerrors.New("failed solve recaptcha")
}
