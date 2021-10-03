package line

import (
	"github.com/bot-sakura/frugal"
	"github.com/google/uuid"
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
	"strings"
)

// ClientSetting line client setting
type ClientSetting struct {
	AppType   model.ApplicationType
	Proxy     string
	KeeperDir string
}

type ClientInfo struct {
	Device      *model.Device
	PhoneNumber *model.UserPhoneNumber
}

// Client line client
type Client struct {
	ctx           frugal.FContext
	ClientSetting *ClientSetting
	ClientInfo    *ClientInfo

	thriftFactory *thriftFactory
}

func newLineDevice() *model.Device {
	uuidObj, _ := uuid.NewUUID()
	return &model.Device{
		Udid:        strings.Join(strings.Split(uuidObj.String(), "-"), ""),
		DeviceModel: genRandomDeviceModel(),
	}
}

// create default line client
func newDefaultClient() *Client {
	cl := &Client{
		ctx: frugal.NewFContext(""),
		ClientSetting: &ClientSetting{
			AppType:   model.ApplicationType_ANDROID,
			KeeperDir: "./keepers/",
		},
		ClientInfo: &ClientInfo{
			Device: newLineDevice(),
		},
	}
	cl.thriftFactory = newThriftFactory(cl)
	return cl
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
