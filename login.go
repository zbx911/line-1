package line

func (cl *Client) afterLoginWithKeeperLoad() error {
	_ = cl.LoadKeeper()
	return cl.afterLogin()
}

func (cl *Client) afterLogin() error {
	if err := cl.executeOpts(); err != nil {
		return err
	}
	cl.ThriftFactory = newThriftFactory(cl)
	if err := cl.setupSessions(); err != nil {
		return err
	}
	if cl.TokenManager.IsV3Token {
		cl.runTokenUpdaterBackGround()
	}
	return nil
}

func (cl *Client) LoginViaAuthKey(key string) error {
	cl.Profile.Mid = key[:33]
	_ = cl.LoadKeeper()
	//AccessToken should be PrimaryToken
	if cl.TokenManager.AccessToken != "" {
		return cl.LoginViaPrimaryToken(cl.TokenManager.AccessToken)
	}
	return cl.LoginViaPrimaryToken(cl.GeneratePrimaryToken(key))
}

func (cl *Client) LoginViaKeeper(mid string) error {
	cl.Profile.Mid = mid
	if err := cl.LoadKeeper(); err != nil {
		return err
	}
	return cl.afterLoginWithKeeperLoad()
}

func (cl *Client) LoginViaPrimaryToken(token string) error {
	cl.Profile.Mid = token[:33]
	cl.TokenManager.AccessToken = token
	cl.TokenManager.IsV3Token = false
	return cl.afterLoginWithKeeperLoad()
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

	oldClient := New(cl.opts...)
	err = oldClient.LoginViaKeeper(cl.Profile.Mid)
	if err == nil {
		if !oldClient.TokenManager.IsV3Token {
			return cl.afterLoginWithKeeperLoad()
		}
		oldToken, err := oldClient.TokenManager.parseV3Token()
		if err != nil {
			return err
		}
		if oldToken.IssuedAt > token.IssuedAt {
			cl = oldClient
			return cl.afterLogin()
		}
	}
	return cl.afterLoginWithKeeperLoad()
}
