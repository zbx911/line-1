package line

import "github.com/line-api/model/go/model"

func (cl *Client) afterError(err error) error {
	if err == nil {
		return nil
	}
	talkErr, ok := err.(*model.TalkException)
	if !ok {
		return err
	}
	fnc, ok := cl.ClientSetting.AfterTalkError[talkErr.Code]
	if !ok {
		return err
	}
	return fnc(talkErr)
}
