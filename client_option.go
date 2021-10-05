package line

import (
	"github.com/line-api/model/go/model"
	"github.com/phuslu/log"
)

type ClientOption func(client *Client) error

// ApplicationType set line client application type
func ApplicationType(appType model.ApplicationType) ClientOption {
	return func(client *Client) error {
		client.ClientSetting.AppType = appType
		return nil
	}
}

// Proxy set line client proxy
func Proxy(proxy string) ClientOption {
	return func(client *Client) error {
		client.ClientSetting.Proxy = proxy
		return nil
	}
}

// KeeperDir set line client keepers path
func KeeperDir(path string) ClientOption {
	return func(client *Client) error {
		client.ClientSetting.KeeperDir = path
		return nil
	}
}

func LocalAddr(addr string) ClientOption {
	return func(client *Client) error {
		client.ClientSetting.LocalAddr = addr
		return nil
	}
}

func AfterTalkError(fncs map[model.TalkErrorCode]func(err *model.TalkException) error) ClientOption {
	return func(client *Client) error {
		for k, v := range fncs {
			client.ClientSetting.AfterTalkError[k] = v
		}
		return nil
	}
}

func Logger(logger *log.Logger) ClientOption {
	return func(client *Client) error {
		client.ClientSetting.Logger = logger
		return nil
	}
}
