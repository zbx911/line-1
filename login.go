package line

func (cl *Client) afterLogin() error {
	if err := cl.LoadKeeper(); err != nil {
		return err
	}
	if err := cl.executeOpts(); err != nil {
		return err
	}
	if err := cl.setupSessions(); err != nil {
		return err
	}
	return nil
}

func (cl *Client) LoginViaAuthKey(key string) error {
	return cl.LoginViaPrimaryToken(cl.GeneratePrimaryToken(key))
}

func (cl *Client) LoginViaKeeper(mid string) error {
	cl.Profile.Mid = mid
	return cl.afterLogin()
}

func (cl *Client) LoginViaPrimaryToken(token string) error {
	cl.TokenManager.IsV3Token = false
	cl.TokenManager.AccessToken = token
	return cl.afterLogin()
}

func (cl *Client) LoginViaV3Token(accessToken, refreshToken string) error {
	cl.TokenManager.AccessToken = accessToken
	cl.TokenManager.RefreshToken = refreshToken
	cl.TokenManager.IsV3Token = true
	return cl.afterLogin()
}
