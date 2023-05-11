package polka

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) EstimateTransactionFee(transaction base.Transaction) (fee *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	txn, ok := transaction.(*Transaction)
	if !ok {
		return nil, base.ErrInvalidTransactionType
	}

	account := mockAccount()
	fakeHash := "0x38c5a9f6fabb8d8583ed633c469cdeefb988b0d2384937b15e10e9c0a75aa744"
	signData, err := txn.GetSignData(fakeHash, 0, 0, 0)
	if err != nil {
		return
	}
	signature, err := account.Sign(signData, "")
	if err != nil {
		return
	}
	sendTx, err := txn.GetTx(account.PublicKey(), signature)
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
		return nil, errors.New("get estimated fee result nil")
	}

	return &base.OptionalString{Value: estimateFee}, nil
}

func (c *Chain) EstimateTransactionFeeUsePublicKey(transaction base.Transaction, pubkey string) (fee *base.OptionalString, err error) {
	return c.EstimateTransactionFee(transaction)
}
