package line

import "github.com/line-api/model/go/model"

type E2EEService struct {
	client *Client
	conn   *model.FTalkServiceClient
}

func (cl *Client) newE2EEService() *E2EEService {
	return &E2EEService{
		client: cl,
		conn:   cl.thriftFactory.newTalkServiceClient(),
	}
}

func (e *E2EEService) NegotiateE2EEPublicKey(mid string) (*model.E2EENegotiation, error) {
	return e.conn.NegotiateE2EEPublicKey(e.client.ctx, mid)
}
