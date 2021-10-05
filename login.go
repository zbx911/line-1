package line

func (cl *Client) afterLogin() error {
	_ = cl.LoadKeeper()
	if err := cl.executeOpts(); err != nil {
		return err
	}
	cl.thriftFactory = newThriftFactory(cl)
	if err := cl.setupSessions(); err != nil {
		return err
	}
	if cl.TokenManager.IsV3Token {
		cl.tokenUpdater()
	}
	return nil
}

func (cl *Client) LoginViaAuthKey(key string) error {
	return cl.LoginViaPrimaryToken(cl.GeneratePrimaryToken(key))
}

func (cl *Client) LoginViaKeeper(mid string) error {
	cl.Profile.Mid = mid
	if err := cl.LoadKeeper(); err != nil {
		return err
	}
	return cl.afterLogin()
}

func (cl *Client) LoginViaPrimaryToken(token string) error {
	cl.Profile.Mid = token[:33]
	cl.TokenManager.AccessToken = token
	cl.TokenManager.IsV3Token = false
	return cl.afterLogin()
}

func (cl *Client) LoginViaV3Token(accessToken, refreshToken string) error {
	cl.TokenManager.AccessToken = accessToken
	cl.TokenManager.RefreshToken = refreshToken
	cl.TokenManager.IsV3Token = true
	token, err := cl.TokenManager.parseV3Token()
	if err != nil {
		return err
	}
	cl.Profile.Mid = token.AuthorId
	return cl.afterLogin()
}
