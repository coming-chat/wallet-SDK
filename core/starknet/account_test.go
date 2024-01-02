package starknet

import (
	"math/big"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func M1Account(t *testing.T) *Account {
	mnemonic := testcase.M1
	acc, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)
	return acc
}

func TestEncodeAddress(t *testing.T) {
	{
		// ArgentX
		// see https://testnet.starkscan.co/tx/0x04001e02d0397ebd0821f7c6865f0914e293955417a33505f99e0e1ec1182ea3
		pub := "0x28081ae2bc3668241b1303df98a61e229ee760eb554f9c7fb21cd968a1b74b1"
		param := deployParamForArgentX(*mustFelt(pub))
		addr, err := param.ComputeContractAddress()
		require.Nil(t, err)
		require.Equal(t, addr.String(), "0x7384b9770dce88ee83a62a8a0ab0fac476e513a9e4b611b80fa08e844ce1f2")
	}
	{
		// Braavos
		// see https://testnet.starkscan.co/tx/0x2d72531b049bcf72dbaa4730161e082798e10fa849763f12b3788f7c275b682
		pub := "0x28081ae2bc3668241b1303df98a61e229ee760eb554f9c7fb21cd968a1b74b1"
		param := deployParamForBraavos(*mustFelt(pub))
		addr, err := param.ComputeContractAddress()
		require.Nil(t, err)
		require.Equal(t, addr.String(), "0x8debaf4740ac184b2e879d4d3fd773f2c7f5d453b795212d4098899a73fc19")
	}
}

func TestAccount(t *testing.T) {
	mnemonic := testcase.M1
	account, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)

	prikey, err := account.PrivateKeyHex()
	require.Nil(t, err)

	account2, err := AccountWithPrivateKey(prikey)
	require.Nil(t, err)
	require.Equal(t, account.PublicKey(), account2.PublicKey())
	require.Equal(t, account.Address(), account2.Address())

	t.Log(prikey)
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())

	require.Equal(t, account.Address(), "0x7384b9770dce88ee83a62a8a0ab0fac476e513a9e4b611b80fa08e844ce1f2")
}

func TestAccount_ImportPrivateKey(t *testing.T) {
	priHex := "0x1234567890"
	priDecimal := "78187493520"
	require.Equal(t, parseNumber(t, priHex), parseNumber(t, priDecimal))

	accountHex, err := AccountWithPrivateKey(priHex)
	require.Nil(t, err)
	accountDecimal, err := AccountWithPrivateKey(priDecimal)
	require.Nil(t, err)

	require.Equal(t, accountHex.PublicKey(), accountDecimal.PublicKey())
	require.Equal(t, accountHex.Address(), accountDecimal.Address())
	require.Equal(t, accountHex.Address(), "0x320d810722501737687ac57ad932e21b1b19d603d131522d5984dd6ca452226")
}

func TestGrindKey(t *testing.T) {
	prikey := "86F3E7293141F20A8BAFF320E8EE4ACCB9D4A4BF2B4D295E8CEE784DB46E0519"
	seed, ok := big.NewInt(0).SetString(prikey, 16)
	require.True(t, ok)
	res, err := grindKey(seed.Bytes())
	require.Nil(t, err)
	require.Equal(t, res.Text(16), "5c8c8683596c732541a59e03007b2d30dbbbb873556fe65b5fb63c16688f941")
}

func parseNumber(t *testing.T, num string) *big.Int {
	bn, err := base.ParseNumber(num)
	require.Nil(t, err)
	return bn
}
