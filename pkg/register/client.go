package register

import (
	"github.com/line-api/line"
	"github.com/line-api/line/pkg/register/phone"
	"github.com/line-api/line/pkg/register/recaptcha"
	"os"
)

type Client struct {
	lineClient *line.Client

	recaptchaSolver recaptcha.Solver
	phoneService    phone.Service
}

func New(lineCl *line.Client, opts ...ClientOption) *Client {
	cl := &Client{
		lineClient:      lineCl,
		recaptchaSolver: recaptcha.NewTwoCaptcha(os.Getenv("TWO_CAPTCHA_API_KEY")),
		phoneService:    phone.NewFiveSim(os.Getenv("FIVE_SIM_API_KEY")),
	}
	for _, op := range opts {
		op(cl)
	}
	return cl
}
