package register

import (
	"github.com/line-api/line/pkg/register/phone"
	"github.com/line-api/line/pkg/register/recaptcha"
)

type ClientOption func(client *Client)

func RecaptchatService(solver recaptcha.Solver) ClientOption {
	return func(client *Client) {
		client.recaptchaSolver = solver
	}
}
func PhoneService(service phone.Service) ClientOption {
	return func(client *Client) {
		client.phoneService = service
	}
}
