package register

import (
	"github.com/line-api/line"
	"github.com/line-api/line/plugin/register/phone"
	"github.com/line-api/line/plugin/register/recaptcha"
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

func Password(pswd string) ClientOption {
	return func(client *Client) {
		client.Password = pswd
	}
}

func Name(name string) ClientOption {
	return func(client *Client) {
		client.Name = name
	}
}

func Debug(flag bool) ClientOption {
	return func(client *Client) {
		client.Debug = flag
	}
}

func AfterCreates(funcs ...func(client *line.Client) error) ClientOption {
	return func(client *Client) {
		client.afterCreates = append(client.afterCreates, funcs...)
	}
}
