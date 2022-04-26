package polka

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) EstimateFeeForTransaction(transaction *Transaction) (s string, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	s = "0"

	account := mockAccount()
	fakeHash := "0x38c5a9f6fabb8d8583ed633c469cdeefb988b0d2384937b15e10e9c0a75aa744"
	signData, err := transaction.GetSignData(fakeHash, 0, 0, 0)
	if err != nil {
		return
	}
	signature, err := account.Sign(signData, "")
	if err != nil {
		return
	}
	pubkey, err := types.HexDecodeString(account.PublicKey())
	if err != nil {
		return
	}
	sendTx, err := transaction.GetTx(pubkey, signature)
	if err != nil {
		return
	}

	cl, err := getConnectedPolkaClient(c.RpcUrl)
	data := make(map[string]interface{})
	err = client.CallWithBlockHash(cl.api.Client, &data, "payment_queryInfo", nil, sendTx)
	if err != nil {
		return
	}

	estimateFee, ok := data["partialFee"].(string)
	if !ok {
		return s, errors.New("get estimated fee result nil")
	}

	return estimateFee, nil
}
