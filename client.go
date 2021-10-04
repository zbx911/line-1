package line

import (
	"fmt"
	"github.com/bot-sakura/frugal"
	"github.com/google/uuid"
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
	"strings"
)

// ClientSetting line client setting
type ClientSetting struct {
	AppType        model.ApplicationType
	Proxy          string
	LocalAddr      string
	KeeperDir      string
	AfterTalkError map[model.TalkErrorCode]func(err *model.TalkException) error
}

type ClientInfo struct {
	Device      *model.Device
	PhoneNumber *model.UserPhoneNumber
}

// Client line client
type Client struct {
	*PollService
	*ChannelService
	*TalkService
	*AccessTokenRefreshService

	opts            []ClientOption
	ctx             frugal.FContext
	ClientSetting   *ClientSetting
	ClientInfo      *ClientInfo
	RequestSequence int32
	thriftFactory   *thriftFactory

	Profile  *model.Profile
	Settings *model.Settings

	TokenManager *TokenManager
}

func (cl *Client) setupSessions() error {
	cl.PollService = cl.newPollService()
	cl.ChannelService = cl.newChannelService()
	cl.TalkService = cl.newTalkService()
	cl.AccessTokenRefreshService = cl.newAccessTokenRefreshService()
	return nil
}

func (cl *Client) executeOpts() error {
	for idx, opt := range cl.opts {
		err := opt(cl)
		if err != nil {
			return xerrors.Errorf("failed to execute %v option: %w", idx, err)
		}
	}
	return nil
}

func (cl *Client) getLineApplicationHeader() string {
	switch cl.ClientSetting.AppType {
	case model.ApplicationType_ANDROID:
		return fmt.Sprintf("ANDROID\t%v\tAndroid OS\t%v", AndroidAppVersion, AndroidVersion)
	case model.ApplicationType_ANDROIDLITE:
		return fmt.Sprintf("ANDROIDLITE\t%v\tAndroid OS\t%v", AndroidLiteAppVersion, AndroidVersion)
	case model.ApplicationType_IOS:
		return "IOS\t11.9.0\tiOS\t14.5.1"
	}
	panic("unsupported app type")
}

func (cl *Client) getLineUserAgentHeader() string {
	switch cl.ClientSetting.AppType {
	case model.ApplicationType_ANDROID:
		return fmt.Sprintf("Line/%v", AndroidAppVersion)
	case model.ApplicationType_ANDROIDLITE:
		return fmt.Sprintf("LLA/%v %v %v", AndroidLiteAppVersion, cl.ClientInfo.Device.DeviceModel, AndroidVersion)
	case model.ApplicationType_IOS:
		return "Line/11.9.0"
	}
	panic("unsupported app type")
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
		TokenManager: &TokenManager{},
		Profile:      &model.Profile{},
		Settings:     &model.Settings{},
	}
	return cl
}

// New create new line client
func New(opts ...ClientOption) (*Client, error) {
	cl := newDefaultClient()
	cl.opts = opts
	return cl, nil
}
