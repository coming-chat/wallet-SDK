package wallet

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/eth"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := NewAccountCache()

	address := "1234567890zxcvbnm"
	cache.SaveAccountInfo("main", "bitcoin", &AccountInfo{Address: address})

	info := cache.GetAccountInfo("main", "bitcoin")
	require.Equal(t, address, info.Address)

	cache.Clean()

	info2 := cache.GetAccountInfo("main", "bitcoin")
	require.Nil(t, info2)

	ethAccount, err := eth.NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	cache.SaveAccount("wallet2", "ethereum", ethAccount)

	info3 := cache.GetAccountInfo("main", "ethereum")
	require.Nil(t, info3)
	info3 = cache.GetAccountInfo("wallet2", "ethereum")
	require.Equal(t, info3.Address, ethAccount.Address())
	require.Equal(t, info3.PublicKey, ethAccount.PublicKey())

	cache.SaveAccount("wallet2", "ethereum", nil)
	info4 := cache.GetAccountInfo("wallet2", "ethereum")
	require.Nil(t, info4)

	// remove not exists key
	cache.SaveAccount("wallet2", "ethereum", nil)
	cache.Delete("xxxxxxx")
	t.Log("delete success")

}
