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

func (cl *Client) LoginViaKeeper(mid string) error {
	cl.Profile.Mid = mid
	return cl.afterLogin()
}

func (cl *Client) LoginViaPrimaryToken(token string) error {
	return nil
}
