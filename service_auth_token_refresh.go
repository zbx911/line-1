package line

import "github.com/line-api/model/go/model"

type AccessTokenRefreshService struct {
	client *Client

	conn *model.FAccessTokenRefreshServiceClient
}

func (cl *Client) newAccessTokenRefreshService() *AccessTokenRefreshService {
	return &AccessTokenRefreshService{
		client: cl,
		conn:   cl.thriftFactory.newAccessTokenRefreshService(),
	}
}

func (cl *Client) RefreshV3AccessToken() error {
	if !cl.TokenManager.IsV3Token {
		return nil
	}
	response, err := cl.Refresh(cl.TokenManager.RefreshToken)
	if err != nil {
		return err
	}
	cl.TokenManager.RefreshToken = response.RefreshToken
	cl.TokenManager.AccessToken = response.AccessToken
	cl.TokenManager.IsV3Token = true
	return cl.ReportRefreshedAccessToken(response.AccessToken)
}

func (s *AccessTokenRefreshService) Refresh(token string) (*model.RefreshAccessTokenResponse, error) {
	return s.conn.Refresh(s.client.ctx, &model.RefreshAccessTokenRequest{
		RefreshToken: token,
	})
}

func (s *AccessTokenRefreshService) ReportRefreshedAccessToken(token string) error {
	_, err := s.conn.ReportRefreshedAccessToken(s.client.ctx, &model.ReportRefreshedAccessTokenRequest{
		AccessToken: token,
	})
	return err
}
