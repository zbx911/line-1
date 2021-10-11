package line

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/imroc/req"
	"github.com/line-api/model/go/model"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//TODO:req moduleを使用しないようにする

type ChannelToken map[string]string

func (c ChannelToken) get(key string) (string, bool) {
	v, ok := c[key]
	return v, ok
}

type ChannelService struct {
	client       *Client
	conn         *model.FChannelServiceClient
	ChannelToken ChannelToken
}

func (cl *Client) newChannelService() *ChannelService {
	return &ChannelService{
		client:       cl,
		conn:         cl.ThriftFactory.newChannelServiceClient(),
		ChannelToken: ChannelToken{},
	}
}

func (s *ChannelService) IssueChannelToken(token string) (*model.ChannelToken, error) {
	t, err := s.conn.IssueChannelToken(s.client.ctx, token)
	return t, s.client.afterError(err)
}

func (s *ChannelService) InitChannelToken() error {
	channelToken, err := s.IssueChannelToken("1341209950")
	if err != nil {
		return err
	}
	s.ChannelToken["1341209950"] = channelToken.ChannelAccessToken
	return nil
}

func (s *ChannelService) UpdateGroupPicture(gid, imagePath string) error {
	header := make(http.Header)
	header.Set("X-Line-Application", s.client.GetLineApplicationHeader())
	header.Set("X-Line-Access", s.client.TokenManager.AccessToken)
	header.Set("X-Lal", "ja_jp")
	header.Set("Quality", "95")
	header.Set("X-Line-Region", "CA")
	header.Set("X-Line-Carrier", "44070")
	header.Set("User-Agent", s.client.GetLineUserAgentHeader())
	header.Set("Content-Type", "image/jpeg")
	file, _ := os.Open(imagePath)
	if s.client.ClientSetting.Proxy != "" {
		req.SetProxyUrl(s.client.ClientSetting.Proxy)
	}
	host := "https://obs-jp.line-apps.com/os/g/" + gid
	_, err := req.Post(host, header, file)
	return err
}

func (s *ChannelService) DownloadObjMessage(msgId, path string) error {
	r, err := http.NewRequest("GET", "https://obs-jp.line-apps.com/r/talk/m/"+msgId, nil)
	if err != nil {
		return err
	}
	r.Host = "obs-jp.line-apps.com"
	r.Header.Set("X-Line-Application", s.client.GetLineApplicationHeader())
	r.Header.Set("X-Line-Carrier", "44070")
	r.Header.Set("User-Agent", s.client.GetLineUserAgentHeader())
	r.Header.Set("X-Line-Access", s.client.TokenManager.AccessToken)
	r.Header.Set("X-Lal", "ja_jp")
	r.Header.Set("X-Line-Region", "CA")
	//r.Header.Set("Range", "bytes=0-22700")

	resp, err := s.client.ThriftFactory.HttpClient().Do(r)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChannelService) UpdateProfilePicture(path string) error {
	host := "https://obs-jp.line-apps.com/os/p/" + s.client.Profile.Mid
	header := make(http.Header)
	header.Set("X-Line-Application", s.client.GetLineApplicationHeader())
	header.Set("X-Line-Access", s.client.TokenManager.AccessToken)
	header.Set("X-Lal", "ja_jp")
	header.Set("Quality", "95")
	header.Set("X-Line-Region", "CA")
	header.Set("X-Line-Carrier", "44070")
	header.Set("User-Agent", s.client.GetLineUserAgentHeader())
	header.Set("Content-Type", "image/jpeg")

	file, _ := os.Open(path)
	if s.client.ClientSetting.Proxy != "" {
		req.SetProxyUrl(s.client.ClientSetting.Proxy)
	}
	_, err := req.Post(host, header, file)
	return err
}

func (s *ChannelService) UpdateProfileCover(path string) error {
	oid, err := s.UploadObjHome(path)
	if err != nil {
		return err
	}
	err = s.UpdateProfileCoverById(oid)
	return err
}

func (s *ChannelService) UpdateProfileCoverById(objId string) error {
	data := map[string]string{
		"homeId":        s.client.Profile.Mid,
		"coverObjectId": objId,
		"storyShare":    "true",
	}
	header := make(http.Header)
	for k, v := range s.client.ThriftFactory.header() {
		header.Set(k, v)
	}
	for k, v := range map[string]string{
		"x-line-access":             s.client.TokenManager.AccessToken,
		"x-lpv":                     "1",
		"x-line-global-config":      "discover.enable=false; follow.enable=true",
		"x-line-bdbtemplateversion": "v1",
		"user-agent":                "androidapp.line/11.4.1 (Linux; U; Android 5.1.1; en-GB; GA00747-UK Build/LMY48Z)",
		"x-lsr":                     "CA",
		"content-type":              "application/json; charset=UTF-8",
	} {
		header.Set(k, v)
	}
	if s.client.ClientSetting.Proxy != "" {
		req.SetProxyUrl(s.client.ClientSetting.Proxy)
	}
	_, err := req.Post("https://ga2.model.naver.jp/hm/api/v1/home/cover.json", header, req.BodyJSON(data))
	return err
}

func (s *ChannelService) UploadObjHome(path string) (string, error) {
	header := make(http.Header)
	for k, v := range s.client.ThriftFactory.header() {
		header.Set(k, v)
	}
	hstr := fmt.Sprintf("Line_%d", time.Now().Unix())
	objId := fmt.Sprintf("%x", md5.Sum([]byte(hstr)))
	file, _ := os.Open(path)
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	for k, v := range map[string]string{
		"x-obs-params": genObsParam(map[string]string{
			"name":   fmt.Sprintf("%d", time.Now().Unix()),
			"userid": s.client.Profile.Mid,
			"oid":    objId,
			"type":   "image",
			"ver":    "1.0",
		}),
		"Content-Type":   "image/jpeg",
		"Content-Length": fmt.Sprintf("%d", fi.Size()),
	} {
		header.Set(k, v)
	}
	if s.client.ClientSetting.Proxy != "" {
		req.SetProxyUrl(s.client.ClientSetting.Proxy)
	}
	_, err = req.Post("https://obs-jp.line-apps.com/myhome/c/upload.nhn", file, header)
	if err != nil {
		return "", err
	}
	return objId, nil
}

func genObsParam(dict map[string]string) string {
	marshal, _ := json.Marshal(dict)
	return base64.StdEncoding.EncodeToString(marshal)
}

func (s *ChannelService) DownloadGroupPicture(picPath, path string) error {
	r, err := http.NewRequest("GET", "https://profile.line-scdn.net/"+picPath, nil)
	if err != nil {
		return err
	}
	r.Header.Set("User-Agent", "okhttp/3.12.6")

	resp, err := s.client.ThriftFactory.HttpClient().Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil

}

func (s *ChannelService) DownloadContactIcon(url, path string) error {
	r, err := http.NewRequest("GET", "https://profile.line-scdn.net/"+url, nil)
	if err != nil {
		// handle err
	}
	r.Host = "profile.line-scdn.net"
	r.Header.Set("User-Agent", "okhttp/3.12.6")

	resp, err := s.client.ThriftFactory.HttpClient().Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChannelService) GetProfileCoverId(mid string) (string, error) {
	r, err := http.NewRequest("GET", "https://ga2.model.naver.jp/hm/api/v1/home/profile.json?homeId="+mid+"&styleMediaVersion=v2&storyVersion=v6", nil)
	if err != nil {
		return "", err
	}
	r.Header.Set("X-Lsr", "CA")
	channelToken, ok := s.ChannelToken.get("1341209950")
	if !ok {
		err := s.InitChannelToken()
		if err != nil {
			return "", err
		}
	}
	r.Header.Set("X-Line-Channeltoken", channelToken)
	r.Header.Set("X-Line-Application", s.client.GetLineApplicationHeader())
	r.Header.Set("Content-Type", "application/json; charset=UTF-8")
	r.Header.Set("X-Lal", "ja_JP")
	r.Header.Set("X-Line-Global-Config", "discover.enable=true; follow.enable=true")
	r.Header.Set("X-Line-Bdbtemplateversion", "v1")
	r.Header.Set("User-Agent", "androidapp.line/11.5.2 (Linux; U; Android 5.1.1; ja-JP; G011A Build/LMY48Z)")
	r.Header.Set("X-Line-Mid", s.client.Profile.Mid)
	r.Header.Set("X-Lpv", "1")
	resp, err := s.client.ThriftFactory.HttpClient().Do(r)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	t := new(ProfileCoverStruct)
	if err := json.Unmarshal(bytes, t); err != nil {
		return "", err
	}
	return t.Result.CoverObsInfo.ObjectId, nil
}

type ProfileCoverStruct struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		HomeId       string `json:"homeId"`
		HomeType     string `json:"homeType"`
		HasNewPost   bool   `json:"hasNewPost"`
		CoverObsInfo struct {
			ObsNamespace string `json:"obsNamespace"`
			ServiceName  string `json:"serviceName"`
			ObjectId     string `json:"objectId"`
		} `json:"coverObsInfo"`
		FollowSummaryInfo struct {
			FollowingCount int  `json:"followingCount"`
			FollowerCount  int  `json:"followerCount"`
			Following      bool `json:"following"`
			AllowFollow    bool `json:"allowFollow"`
			ShowFollowList bool `json:"showFollowList"`
		} `json:"followSummaryInfo"`
		GiftShopInfo struct {
			GiftShopScheme         string `json:"giftShopScheme"`
			BirthdayGiftShopScheme string `json:"birthdayGiftShopScheme"`
			GiftShopUrl            string `json:"giftShopUrl"`
			IsGiftShopAvailable    bool   `json:"isGiftShopAvailable"`
		} `json:"giftShopInfo"`
		UserStyleMedia struct {
			MenuInfo struct {
				LatestEditTime int64 `json:"latestEditTime"`
			} `json:"menuInfo"`
		} `json:"userStyleMedia"`
		Meta struct {
		} `json:"meta"`
	} `json:"result"`
}
