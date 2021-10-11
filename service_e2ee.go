package line

import "github.com/line-api/model/go/model"

type E2EEService struct {
	client             *Client
	conn               *model.FTalkServiceClient
	connCompactE2EEMsg *model.FCompactMessageServiceClient
}

func (cl *Client) newE2EEService() *E2EEService {
	return &E2EEService{
		client:             cl,
		conn:               cl.ThriftFactory.newTalkServiceClient(),
		connCompactE2EEMsg: cl.ThriftFactory.newCompactE2EEMessageServiceClient(),
	}
}

func (s *E2EEService) NegotiateE2EEPublicKey(mid string) (*model.E2EENegotiation, error) {
	return s.conn.NegotiateE2EEPublicKey(s.client.ctx, mid)
}

func (s *E2EEService) SendE2EEMessage(msg *model.Message) (*model.Message, error) {
	return s.connCompactE2EEMsg.SendE2EEMessageCompact(s.client.ctx, s.client.RequestSequence, msg)
}
