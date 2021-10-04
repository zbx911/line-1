package line

import (
	"math/rand"
	"time"
)

type Path string

const (
	PATH_LONG_POLLING                             Path = "/P4"
	PATH_LONG_POLLING_P5                          Path = "/P5"
	PATH_NORMAL_POLLING                           Path = "/NP4"
	PATH_NORMAL                                   Path = "/S4"
	PATH_COMPACT_MESSAGE                          Path = "/C5"
	PATH_COMPACT_PLAIN_MESSAGE                    Path = "/CA5"
	PATH_COMPACT_E2EE_MESSAGE                     Path = "/ECA5"
	PATH_REGISTRATION                             Path = "/api/v4/TalkService.do"
	PATH_REFRESH_TOKEN                            Path = "/EXT/auth/tokenrefresh/v1"
	PATH_NOTIFY_SLEEP                             Path = "/F4"
	PATH_NOTIFY_BACKGROUND                        Path = "/B"
	PATH_BUDDY                                    Path = "/BUDDY4"
	PATH_SHOP                                     Path = "/SHOP4"
	PATH_SHOP_AUTH                                Path = "/SHOPA"
	PATH_UNIFIED_SHOP                             Path = "/TSHOP4"
	PATH_STICON                                   Path = "/SC4"
	PATH_CHANNEL                                  Path = "/CH4"
	PATH_CANCEL_LONGPOLLING                       Path = "/CP4"
	PATH_SNS_ADAPTER                              Path = "/SA4"
	PATH_SNS_ADAPTER_REGISTRATION                 Path = "/api/v4p/sa"
	PATH_AUTH_EAP                                 Path = "/ACCT/authfactor/eap/v1"
	PATH_USER_INPUT                               Path = ""
	PATH_USER_BEHAVIOR_LOG                        Path = "/L1"
	PATH_AGE_CHECK                                Path = "/ACS4"
	PATH_SPOT                                     Path = "/SP4"
	PATH_CALL                                     Path = "/V4"
	PATH_EXTERNAL_INTERLOCK                       Path = "/EIS4"
	PATH_TYPING                                   Path = "/TS"
	PATH_CONN_INFO                                Path = "/R2"
	PATH_HTTP_PROXY                               Path = ""
	PATH_EXTERNAL_PROXY                           Path = ""
	PATH_PAY                                      Path = "/PY4"
	PATH_WALLET                                   Path = "/WALLET4"
	PATH_AUTH                                     Path = "/RS4"
	PATH_AUTH_REGISTRATION                        Path = "/api/v4p/rs"
	PATH_SEARCH_COLLECTION_MENU_V1                Path = "/collection/v1"
	PATH_SEARCH_V2                                Path = "/search/v2"
	PATH_SEARCH_V3                                Path = "/search/v3"
	PATH_BEACON                                   Path = "/BEACON4"
	PATH_PERSONA                                  Path = "/PS4"
	PATH_SQUARE                                   Path = "/SQS1"
	PATH_SQUARE_BOT                               Path = "/BP1"
	PATH_POINT                                    Path = "/POINT4"
	PATH_COIN                                     Path = "/COIN4"
	PATH_LIFF                                     Path = "/LIFF1"
	PATH_CHAT_APP                                 Path = "/CAPP1"
	PATH_IOT                                      Path = "/IOT1"
	PATH_USER_PROVIDED_DATA                       Path = "/UPD4"
	PATH_NEW_REGISTRATION                         Path = "/acct/pais/v1"
	PATH_SECONDARY_QR_LOGIN                       Path = "/ACCT/lgn/sq/v1"
	PATH_USER_SETTINGS                            Path = "/US4"
	PATH_LINE_SPOT                                Path = "/ex/spot"
	PATH_LINE_HOME_V2_SERVICES                    Path = "/EXT/home/sapi/v4p/hsl"
	PATH_LINE_HOME_V2_CONTENTS_RECOMMENDATIONPath Path = "/EXT/home/sapi/v4p/flex"
	PATH_BIRTHDAY_GIFT_ASSOCIATION                Path = "/EXT/home/sapi/v4p/bdg"
	PATH_SECONDARY_PWLESS_LOGIN_PERMIT            Path = "/ACCT/lgn/secpwless/v1"
	PATH_SECONDARY_AUTH_FACTOR_PIN_CODE           Path = "/ACCT/authfactor/second/pincode/v1"
	PATH_PWLESS_CREDENTIAL_MANAGEMENT             Path = "/ACCT/authfactor/pwless/manage/v1"
	PATH_PWLESS_PRIMARY_REGISTRATION              Path = "/ACCT/authfactor/pwless/v1"
	PATH_GLN_NOTIFICATION_STATUS                  Path = "/gln/webapi/graphql"
	PATH_BOT_EXTERNAL                             Path = "/BOTE"
	PATH_E2EE_KEY_BACKUP                          Path = "/EKBS4"
)

const (
	LINE_SERVER_HOST     = "https://legy-jp-addr.line.naver.jp"
	LINE_SERVER_HOST_gxx = "https://gxx.line.naver.jp"
)

//ToURL get full url for the path
func (p Path) ToURL() string {
	return LINE_SERVER_HOST + string(p)
}

var AndroidVersions = []string{
	"11.0.0", "10.0.0", "9.0.0", "8.1.0", "8.0.0", "7.1.2", "7.1.1", "7.1.0", "7.0.0",
	"6.0.1", "6.0.0", "5.1.1", "5.1.0", "5.0.2", "5.0.1", "5.0.0",
}

func getRandomAndroidVersion() string {
	rand.Seed(time.Now().Unix())
	return AndroidVersions[rand.Intn(len(AndroidVersions))]
}

var (
	AndroidAppVersion     = "11.16.2"
	AndroidVersion        = getRandomAndroidVersion()
	AndroidLiteAppVersion = "2.17.0"
)
