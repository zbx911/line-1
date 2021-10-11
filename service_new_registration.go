package line

import (
	"bytes"
	"encoding/json"
	"github.com/line-api/line/crypt"
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
	"net/http"
)

type NewRegistrationService struct {
	client *Client
	conn   *model.FPrimaryAccountInitServiceClient
}

func (cl *Client) newNewRegistrationService() *NewRegistrationService {
	return &NewRegistrationService{
		client: cl,
		conn:   cl.thriftFactory.newNewRegistrationService(),
	}
}

func (s *NewRegistrationService) OpenSession() (string, error) {
	return s.conn.OpenSession(s.client.ctx, &model.OpenSessionRequest{})
}

func (s *NewRegistrationService) NotifyInstalled(deviceUUID, lineApplication string) error {
	cl := s.client.thriftFactory.newRegistrationServiceClient()
	return cl.NotifyInstalled(s.client.ctx, deviceUUID, lineApplication)
}

func (s *NewRegistrationService) GetPhoneVerifMethodV2(sessionId string, device *model.Device, phone *model.UserPhoneNumber) (*model.GetPhoneVerifMethodV2Response, error) {
	return s.conn.GetPhoneVerifMethodV2(s.client.ctx, &model.GetPhoneVerifMethodV2Request{
		AuthSessionId:   sessionId,
		Device:          device,
		UserPhoneNumber: phone,
	})
}

func (s *NewRegistrationService) SendPinCodeFOrPhone(sessionId string, phone *model.UserPhoneNumber, method model.PhoneVerifMethodV2) (*model.ReqToSendPhonePinCodeResponse, error) {
	return s.conn.RequestToSendPhonePinCode(s.client.ctx, &model.ReqToSendPhonePinCodeRequest{
		AuthSessionId:   sessionId,
		UserPhoneNumber: phone,
		VerifMethod:     method,
	})
}

func (s *NewRegistrationService) VerifyPhonePinCode(sessionId string, phone *model.UserPhoneNumber, pinCode string) (*model.VerifyPhonePinCodeResponse, error) {
	return s.conn.VerifyPhonePinCode(s.client.ctx, &model.VerifyPhonePinCodeRequest{
		AuthSessionId:   sessionId,
		UserPhoneNumber: phone,
		PinCode:         pinCode,
	})
}

func (s *NewRegistrationService) ValidateProfile(session, name string) error {
	_, err := s.conn.ValidateProfile(s.client.ctx, session, name)
	return err
}

func (s *NewRegistrationService) ExchangeEncryptionKey(sessionId string, keyVer model.EncryptionKeyVersion, encKey *crypt.KeyPairForCurve25519) (*model.ExchangeEncryptionKeyResponse, error) {
	return s.conn.ExchangeEncryptionKey(s.client.ctx, sessionId, &model.ExchangeEncryptionKeyRequest{
		AuthKeyVersion: keyVer,
		PublicKey:      encKey.B64EncodedStringPubKey(),
		Nonce:          encKey.B64EncodedStringNonce(),
	})
}

func (s *NewRegistrationService) SetPassword(password string, sessionId string, exRes *model.ExchangeEncryptionKeyResponse, myKey *crypt.KeyPairForCurve25519) error {
	cText, err := crypt.EncryptPassword(password, exRes, myKey)
	if err != nil {
		return xerrors.Errorf("something went wrong on encrypting password: " + err.Error())
	}
	_, err = s.conn.SetPassword(s.client.ctx, sessionId, &model.EncryptedPassword{
		EncryptionKeyVersion: model.EncryptionKeyVersion_V1,
		CipherText:           cText,
	})
	return err
}

func (s *NewRegistrationService) RegisterPrimaryUsingPhoneWithTokenV3(sessionId string) (*model.RegisterPrimaryWithTokenV3Response, error) {
	result, err := s.conn.RegisterPrimaryUsingPhoneWithTokenV3(s.client.ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *NewRegistrationService) SendGoogleRecaptchaTokenToLine(details *model.WebAuthDetails, token string) error {
	jsonStr, err := json.Marshal(map[string]string{"verifier": token})
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", "https://w.line.me/sec/v3/recaptcha/result/verify", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	for k, v := range map[string]string{
		"Connection":       "keep-alive",
		"user-agent":       "Mozilla/5.0 (Linux; Android 8.0.0; Nexus 6P Build/OPR6.170623.019; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/88.0.4324.93 Mobile Safari/537.36",
		"Accept":           "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
		"Referer":          "https://w.line.me/sec/v3/recaptcha",
		"Origin":           "https://w.line.me",
		"X-Requested-With": "jp.naver.line.android",
		"Content-Type":     "application/json;charset=UTF-8",
		"Host":             "w.line.me",
	} {
		request.Header.Set(k, v)
	}
	request.AddCookie(&http.Cookie{
		Name:  "lsct_acct_init",
		Value: details.Token[15:],
	})
	response, err := s.client.thriftFactory.HttpClient().Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == 200 {
		return nil
	}
	return xerrors.New("failed to solve human verify process")

}
