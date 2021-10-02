package line

import "github.com/line-api/model/go/model"

type thriftFactory struct {
	client *Client
}

func newThriftFactory(cl *Client) *thriftFactory {
	return &thriftFactory{client: cl}
}

func (f *thriftFactory) newPollService() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newPollTMCPService() *model.FTalkServiceClient {
	return nil
}

func (f *thriftFactory) newTalkService() *model.FTalkServiceClient {
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
