package line

import (
	"github.com/line-api/model/go/model"
	"net"
	"net/http"
	"net/url"
)

type thriftFactory struct {
	client        *Client
	httpClient    *http.Client
	defaultHeader map[string]string
}

func newThriftFactory(cl *Client) *thriftFactory {
	return &thriftFactory{
		client:     cl,
		httpClient: cl.defaultHttpClient(),
		defaultHeader: map[string]string{
			"x-line-application": cl.getLineApplicationHeader(),
			"x-line-access":      cl.TokenManager.AccessToken,
			"user-agent":         cl.getLineUserAgentHeader(),
			"x-lal":              "ja_JP",
			"x-lpv":              "1",
		},
	}
}

func (cl *Client) defaultHttpClient() *http.Client {
	httpClient := &http.Client{Transport: &http.Transport{
		ForceAttemptHTTP2:   true,
		MaxIdleConns:        600,
		MaxIdleConnsPerHost: 200,
	}}
	if cl.ClientSetting.LocalAddr != "" {
		OKAddress, err := net.ResolveTCPAddr("tcp", cl.ClientSetting.LocalAddr)
		if err != nil {
			return httpClient
		}
		httpClient.Transport.(*http.Transport).DialContext = (&net.Dialer{
			LocalAddr: OKAddress,
		}).DialContext
	}
	if cl.ClientSetting.Proxy != "" {
		httpClient.Transport.(*http.Transport).Proxy = parseProxyUrl(cl.ClientSetting.Proxy)
	}
	return httpClient
}

func parseProxyUrl(proxy string) func(*http.Request) (*url.URL, error) {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		return nil
	}
	return http.ProxyURL(proxyUrl)
}

func (f *thriftFactory) header() map[string]string {
	return f.defaultHeader
}

func (f *thriftFactory) newHeaderWithExtra(header map[string]string) map[string]string {
	newHeader := make(map[string]string)
	for k, v := range f.defaultHeader {
		newHeader[k] = v
	}
	for k, v := range header {
		newHeader[k] = v
	}
	return newHeader
}

func (f *thriftFactory) HttpClient() *http.Client {
	return nil
}

func (f *thriftFactory) newPollServiceClient() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newPollTMCPServiceClient() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newTalkServiceClient() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newCompactMessageServiceClient() *model.FCompactMessageServiceClient {
	return nil
}

func (f *thriftFactory) newCompactE2EEMessageServiceClient() *model.FCompactMessageServiceClient {
	return nil
}

func (f *thriftFactory) newBuddyServiceClient() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newRegistrationServiceClient() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newChannelServiceClient() *model.FChannelServiceClient {
	return nil
}

func (f *thriftFactory) newNewRegistrationService() *model.FPrimaryAccountInitService {
	return nil
}

func (f *thriftFactory) newAccessTokenRefreshService() *model.FAccessTokenRefreshServiceClient {
	return nil
}
