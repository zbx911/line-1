package register

import (
	"github.com/bot-sakura/frugal"
	"github.com/line-api/line"
	phone2 "github.com/line-api/line/plugin/register/phone"
	recaptcha2 "github.com/line-api/line/plugin/register/recaptcha"
	"github.com/line-api/model/go/model"
	"os"
)

type Client struct {
	lineClient *line.Client

	recaptchaSolver recaptcha2.Solver
	phoneService    phone2.Service

	sessionId string
	ctx       frugal.FContext
	userPhone *model.UserPhoneNumber

	Password     string
	Name         string
	Debug        bool
	afterCreates []func(client *line.Client) error
}

func New(lineCl *line.Client, opts ...ClientOption) *Client {
	cl := &Client{
		lineClient:      lineCl,
		recaptchaSolver: recaptcha2.NewTwoCaptcha(os.Getenv("TWO_CAPTCHA_API_KEY")),
		phoneService:    phone2.NewFiveSim(os.Getenv("FIVE_SIM_API_KEY")),
	}
	for _, op := range opts {
		op(cl)
	}
	return cl
}
