package line

import (
	"github.com/bot-sakura/frugal"
	"github.com/bot-sakura/thrift"
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
	return f.httpClient
}

func (f *thriftFactory) newPollServiceClient() *model.FTalkServiceClient {
	return model.NewFTalkServiceClient(f.newFrugalClient(PATH_LONG_POLLING.ToURL()))
}

func (f *thriftFactory) newPollTMCPServiceClient() *model.FTalkServiceClient {
	return model.NewFTalkServiceClient(f.newTMCPFrugalClient(PATH_LONG_POLLING_P5.ToURL()))
}

func (f *thriftFactory) newTalkServiceClient() *model.FTalkServiceClient {
	return model.NewFTalkServiceClient(f.newFrugalClient(PATH_NORMAL.ToURL()))
}

func (f *thriftFactory) newCompactMessageServiceClient() *model.FCompactMessageServiceClient {
	return model.NewFCompactMessageServiceClient(f.newFrugalClient(PATH_COMPACT_MESSAGE.ToURL()))
}

func (f *thriftFactory) newCompactE2EEMessageServiceClient() *model.FCompactMessageServiceClient {
	return model.NewFCompactMessageServiceClient(f.newFrugalClient(PATH_COMPACT_E2EE_MESSAGE.ToURL()))
}

func (f *thriftFactory) newBuddyServiceClient() *model.FTalkServiceClient {
	return model.NewFTalkServiceClient(f.newFrugalClient(PATH_BUDDY.ToURL()))
}

func (f *thriftFactory) newRegistrationServiceClient() *model.FTalkServiceClient {
	return model.NewFTalkServiceClient(f.newFrugalClient(PATH_REGISTRATION.ToURL()))
}

func (f *thriftFactory) newChannelServiceClient() *model.FChannelServiceClient {
	return model.NewFChannelServiceClient(f.newFrugalClient(PATH_CHANNEL.ToURL()))
}

func (f *thriftFactory) newNewRegistrationService() *model.FPrimaryAccountInitServiceClient {
	return model.NewFPrimaryAccountInitServiceClient(f.newFrugalClient(PATH_NEW_REGISTRATION.ToURL()))
}

func (f *thriftFactory) newAccessTokenRefreshService() *model.FAccessTokenRefreshServiceClient {
	return model.NewFAccessTokenRefreshServiceClient(f.newFrugalClient(PATH_REFRESH_TOKEN.ToURL()))
}

func (f *thriftFactory) newFrugalClient(hostUrl string) *frugal.FServiceProvider {
	fProtoc := frugal.NewFProtocolFactory(thrift.NewTCompactProtocolFactoryConf(&thrift.TConfiguration{}))
	httpTrans := frugal.NewFHTTPTransportBuilder(f.httpClient, hostUrl).WithRequestHeaders(f.header()).Build()
	provider := frugal.NewFServiceProvider(httpTrans, fProtoc)
	return provider
}

func (f *thriftFactory) newTMCPFrugalClient(hostUrl string) *frugal.FServiceProvider {
	fProtoc := frugal.NewFProtocolFactory(thrift.NewTMoreCompactProtocolFactoryConfAndroidLITE(&thrift.TConfiguration{}))
	httpTrans := frugal.NewFHTTPTransportBuilder(f.httpClient, hostUrl).WithRequestHeaders(f.header()).Build()
	provider := frugal.NewFServiceProvider(httpTrans, fProtoc)
	return provider
}
