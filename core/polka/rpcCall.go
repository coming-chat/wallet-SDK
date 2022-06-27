package polka

func (c *Chain) RpcCall(result interface{}, method string, params ...any) error {
	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return err
	}
	return client.api.Client.Call(result, method, params...)
}
