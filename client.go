package line

import (
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
)

// ClientSetting line client setting
type ClientSetting struct {
	AppType    model.ApplicationType
	Proxy      string
	KeeperPath string
}

// Client line client
type Client struct {
	ClientSetting *ClientSetting
}

// create default line client
func newDefaultClient() *Client {
	return &Client{
		ClientSetting: &ClientSetting{
			AppType: model.ApplicationType_ANDROID,
		},
	}
}

// New create new line client
func New(opts ...ClientOption) (*Client, error) {
	cl := newDefaultClient()
	for idx, opt := range opts {
		err := opt(cl)
		if err != nil {
			return nil, xerrors.Errorf("failed to execute %v option: %w", idx, err)
		}
	}
	return cl, nil
}
