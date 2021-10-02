package line

import "github.com/line-api/model/go/model"

type TalkService struct {
	client *Client

	conn *model.FTalkServiceClient
}

func (cl *Client) newTalkService() *TalkService {
	return &TalkService{client: cl, conn: cl.thriftFactory.newTalkServiceClient()}
}
