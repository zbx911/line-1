package register

import (
	"github.com/davegardnerisme/phonegeocode"
	"github.com/line-api/model/go/model"
	"golang.org/x/xerrors"
	"strings"
)

func (c *Client) GetValidPhoneNumber() (string, error) {
	phone, err := c.phoneService.GetNumber()
	if err != nil {
		return "", err
	}
	err = c.checkPhoneNumber(phone)
	if err == nil {
		return phone, nil
	}
	switch authErr := toAuthError(xerrors.Unwrap(err)); authErr.Code {
	case model.AuthErrorCode_HUMAN_VERIFICATION_REQUIRED:
		if err := c.solveHumanVerification(authErr.WebAuthDetails); err != nil {
			return "", xerrors.Errorf("failed to solve human verification on checking phone number(HUMAN_VERIFICATION_REQUIRED): %w", err)
		}
		return phone, nil
	case model.AuthErrorCode_ILLEGAL_ARGUMENT:
		if strings.Contains(authErr.AlertMessage, "不正な電話番号") {
			c.phoneService.BanNumber()
		} else {
			return "", xerrors.Errorf("unknown error occurred on checking phone number(ILLEGAL_ARGUMENT): %w", authErr)
		}
	}
	c.phoneService.CancelNumber()
	return c.GetValidPhoneNumber()
}

func (c *Client) solveHumanVerification(detail *model.WebAuthDetails) error {
	token, err := c.recaptchaSolver.Solve(detail, c.lineClient.ThriftFactory.HttpClient())
	if err != nil {
		return err
	}
	return c.lineClient.SendGoogleRecaptchaTokenToLine(detail, token)
}

func (c *Client) checkPhoneNumber(number string) error {
	code, err := phonegeocode.New().Country(number)
	if err != nil {
		return xerrors.Errorf("wrong phone number: %v", number)
	}
	c.userPhone = &model.UserPhoneNumber{PhoneNumber: number, CountryCode: code}

	result, err := c.lineClient.GetPhoneVerifMethodV2(c.sessionId, c.lineClient.ClientInfo.Device, c.userPhone)
	if err != nil {
		return xerrors.Errorf("failed to get phone verify method v2: %w", err)
	}
	if !checkSMSMethod(result.AvailableMethods) {
		return xerrors.New("error, sms method not available")
	}
	return nil
}

func checkSMSMethod(methods []model.PhoneVerifMethodV2) bool {
	for _, val := range methods {
		if val == model.PhoneVerifMethodV2_SMS {
			return true
		}
	}
	return false
}
