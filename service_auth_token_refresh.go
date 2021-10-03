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

func (s *AccessTokenRefreshService) Refresh(token string) (*model.RefreshAccessTokenResponse, error) {
	return s.conn.Refresh(s.client.ctx, &model.RefreshAccessTokenRequest{
		RefreshToken: token,
	})
}

func (s *AccessTokenRefreshService) ReportRefreshedAccessToken(token string) (*model.ReportRefreshedAccessTokenResponse, error) {
	return s.conn.ReportRefreshedAccessToken(s.client.ctx, &model.ReportRefreshedAccessTokenRequest{
		AccessToken: token,
	})
}
