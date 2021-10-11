package line

import (
	"fmt"
	"github.com/bot-sakura/frugal"
	"github.com/google/uuid"
	"github.com/line-api/line/pkg/logger"
	"github.com/line-api/model/go/model"
	"github.com/phuslu/log"
	"golang.org/x/xerrors"
	"os"
	"strings"
)

// ClientSetting line client setting
type ClientSetting struct {
	AppType        model.ApplicationType
	Proxy          string
	LocalAddr      string
	KeeperDir      string
	AfterTalkError map[model.TalkErrorCode]func(err *model.TalkException) error `json:"-"`

	Logger *log.Logger `json:"-"`
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
	*E2EEService
	*NewRegistrationService

	HeaderFactory *HeaderFactory

	opts            []ClientOption
	ctx             frugal.FContext
	ClientSetting   *ClientSetting
	ClientInfo      *ClientInfo
	RequestSequence int32
	ThriftFactory   *ThriftFactory `json:"-"`

	Profile  *model.Profile
	Settings *model.Settings

	TokenManager *TokenManager
}

func (cl *Client) setupSessions() error {
	cl.PollService = cl.newPollService()
	cl.ChannelService = cl.newChannelService()
	cl.TalkService = cl.newTalkService()
	cl.AccessTokenRefreshService = cl.newAccessTokenRefreshService()
	cl.E2EEService = cl.newE2EEService()
	cl.NewRegistrationService = cl.newNewRegistrationService()
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

func (cl *Client) GetLineApplicationHeader() string {
	switch cl.ClientSetting.AppType {
	case model.ApplicationType_ANDROID:
		return fmt.Sprintf("ANDROID\t%v\tAndroid OS\t%v", cl.HeaderFactory.AndroidAppVersion, cl.HeaderFactory.AndroidVersion)
	case model.ApplicationType_ANDROIDLITE:
		return fmt.Sprintf("ANDROIDLITE\t%v\tAndroid OS\t%v", cl.HeaderFactory.AndroidLiteAppVersion, cl.HeaderFactory.AndroidVersion)
	case model.ApplicationType_IOS:
		return "IOS\t11.9.0\tiOS\t14.5.1"
	}
	panic("unsupported app type")
}

func (cl *Client) GetLineUserAgentHeader() string {
	switch cl.ClientSetting.AppType {
	case model.ApplicationType_ANDROID:
		return fmt.Sprintf("Line/%v", cl.HeaderFactory.AndroidAppVersion)
	case model.ApplicationType_ANDROIDLITE:
		return fmt.Sprintf("LLA/%v %v %v", cl.HeaderFactory.AndroidLiteAppVersion, cl.ClientInfo.Device.DeviceModel, cl.HeaderFactory.AndroidVersion)
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

func getHomeDir() string {
	home, _ := os.UserHomeDir()
	return home
}

// create default line client
func newDefaultClient() *Client {
	cl := &Client{
		ctx: frugal.NewFContext(""),
		ClientSetting: &ClientSetting{
			AppType:   model.ApplicationType_ANDROID,
			KeeperDir: getHomeDir() + "/.line-keepers/",
			Logger:    logger.New(),
		},
		ClientInfo: &ClientInfo{
			Device: newLineDevice(),
		},
		TokenManager: &TokenManager{},
		Profile:      &model.Profile{},
		Settings:     &model.Settings{},
		HeaderFactory: &HeaderFactory{
			AndroidVersion:        getRandomAndroidVersion(),
			AndroidAppVersion:     getRandomAndroidAppVersion(),
			AndroidLiteAppVersion: getRandomAndroidLiteAppVersion(),
		},
	}
	cl.ClientSetting.AfterTalkError = map[model.TalkErrorCode]func(err *model.TalkException) error{
		model.TalkErrorCode_MUST_REFRESH_V3_TOKEN: func(talkErr *model.TalkException) error {
			err := cl.RefreshV3AccessToken()
			if err != nil {
				return err
			}
			return xerrors.Errorf("update v3 access token done: %w", talkErr)
		},
	}
	return cl
}

// New create new line client
func New(opts ...ClientOption) *Client {
	cl := newDefaultClient()
	cl.opts = opts
	return cl
}
