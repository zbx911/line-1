package line

import "github.com/line-api/model/go/model"

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

// KeeperPath set line client keepers path
func KeeperPath(path string) ClientOption {
	return func(client *Client) error {
		client.ClientSetting.KeeperPath = path
		return nil
	}
}
