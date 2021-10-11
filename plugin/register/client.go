package register

import (
	"github.com/bot-sakura/frugal"
	"github.com/line-api/line"
	"github.com/line-api/line/plugin/register/phone"
	"github.com/line-api/line/plugin/register/recaptcha"
	"github.com/line-api/model/go/model"
	"os"
)

type Client struct {
	lineClient *line.Client

	recaptchaSolver recaptcha.Solver
	phoneService    phone.Service

	sessionId string
	ctx       frugal.FContext
	userPhone *model.UserPhoneNumber

	Password        string
	Name            string
	ProfileIconPath string
	Debug           bool
	afterCreates    []func(client *line.Client) error
}

func New(lineCl *line.Client, opts ...ClientOption) (*Client, error) {
	cl := &Client{
		lineClient:      lineCl,
		recaptchaSolver: recaptcha.NewTwoCaptcha(os.Getenv("TWO_CAPTCHA_API_KEY")),
		phoneService:    phone.NewFiveSim(os.Getenv("FIVE_SIM_API_KEY")),
		Password:        line.MakeRandomStr(12),
		Name:            line.MakeRandomStr(7),
	}
	for _, op := range opts {
		op(cl)
	}
	if err := cl.lineClient.NotifyInstalled(cl.lineClient.ClientInfo.Device.Udid, cl.lineClient.GetLineApplicationHeader()); err != nil {
		return nil, err
	}
	sessionId, err := cl.lineClient.OpenSession()
	if err != nil {
		return nil, err
	}
	cl.sessionId = sessionId
	if err := cl.lineClient.GetCountryInfo(sessionId); err != nil {
		return nil, err
	}
	return cl, nil
}

func toAuthError(err error) *model.AuthException {
	aErr, ok := err.(*model.AuthException)
	if ok {
		return aErr
	}
	return &model.AuthException{
		Code:         -1,
		AlertMessage: err.Error(),
	}
}
