package polka

import "github.com/coming-chat/wallet-SDK/core/base"

// MARK - Implement the protocol Chain.SendTx

func (c *Chain) SendRawTransaction(signedTx string) (s string, err error) {
	defer func() {
		err = base.MapToBasicError(err)
	}()
	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	var hashString string
	err = client.api.Client.Call(&hashString, "author_submitExtrinsic", signedTx)
	if err != nil {
		return
	}

	return hashString, nil
}
