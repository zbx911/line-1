package line

import (
	"github.com/line-api/model/go/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	cl := New()
	assert.Equal(t, model.ApplicationType_ANDROID, cl.ClientSetting.AppType)
	assert.Equal(t, "", cl.ClientSetting.Proxy)
	assert.Equal(t, "C:\\Users\\Sakura/.line-keepers/", cl.ClientSetting.KeeperDir)
}

func TestNewWithOpts(t *testing.T) {
	cl := New(KeeperDir("./piyo/"), ApplicationType(model.ApplicationType_ANDROIDLITE), Proxy("https://example.com:6666"))
	assert.Equal(t, "./piyo/", cl.ClientSetting.KeeperDir)
	assert.Equal(t, model.ApplicationType_ANDROIDLITE, cl.ClientSetting.AppType)
	assert.Equal(t, "https://example.com:6666", cl.ClientSetting.Proxy)
}
