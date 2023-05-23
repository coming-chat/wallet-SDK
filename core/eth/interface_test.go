package eth

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/stretchr/testify/require"
)

func Test_All_Interface_Type(t *testing.T) {
	{
		var baseAccount base.Account = &Account{}
		account, ok := baseAccount.(*Account)
		require.True(t, ok)
		require.NotNil(t, account)
	}
	{
		var baseChain base.Chain = &Chain{}
		chain, ok := baseChain.(*Chain)
		require.True(t, ok)
		require.NotNil(t, chain)
	}
	{
		var baseToken base.Token = &Token{}
		token, ok := baseToken.(*Token)
		require.True(t, ok)
		require.NotNil(t, token)
	}
	{
		var baseToken base.Token = &Erc20Token{}
		token, ok := baseToken.(*Erc20Token)
		require.True(t, ok)
		require.NotNil(t, token)
	}
	{
		var baseTransaction base.Transaction = &Transaction{}
		transaction, ok := baseTransaction.(*Transaction)
		require.True(t, ok)
		require.NotNil(t, transaction)
	}
}