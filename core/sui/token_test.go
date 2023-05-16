package sui

import (
	"context"
	"strconv"
	"testing"

	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/stretchr/testify/require"
)

type SUI float32

func (s SUI) Int64() int64 {
	return int64(s * 1e9)
}
func (s SUI) Uint64() uint64 {
	return uint64(s * 1e9)
}
func (s SUI) String() string {
	return strconv.FormatUint(s.Uint64(), 10)
}

func TestBalance(t *testing.T) {
	address := "0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e"

	chain := TestnetChain()
	b, err := chain.BalanceOfAddress(address)
	require.Nil(t, err)

	t.Log(b, "")
}

func TestTokenBalance(t *testing.T) {
	chain := DevnetChain()
	token, err := NewToken(chain, "0x2d79a3c70aa3f3a3feabbf54b7b520f956c4ef8d::AAA::AAA")
	require.NoError(t, err)

	balance, err := token.BalanceOfAddress("0x2ecb102385afd954bf06f2a3a4ac648eb7a536e0")
	require.NoError(t, err)
	require.Equal(t, "0", balance.Total) // invalid address
}

func TestTokenInfo(t *testing.T) {
	chain := DevnetChain()
	token, err := NewToken(chain, "0x2d79a3c70aa3f3a3feabbf54b7b520f956c4ef8d::AAA::AAA")
	require.NoError(t, err)

	info, err := token.TokenInfo()
	require.Error(t, err) // token not found
	t.Log(info)

	mainToken := NewTokenMain(chain)
	tokenInfo, err := mainToken.TokenInfo()
	require.NoError(t, err)

	t.Log(tokenInfo)
}

func TestToken_EstimateGas(t *testing.T) {
	chain := TestnetChain()
	token := NewTokenMain(chain)

	account := M1Account(t)

	gasPrice, err := chain.GasPrice()
	require.Nil(t, err)
	t.Log(gasPrice)

	gas, err := token.EstimateFees(account, account.Address(), "1000000")
	require.Nil(t, err)
	t.Log(gas.Value)
}

func TestTransfer(t *testing.T) {
	// account := M1Account(t)
	account := M3Account(t)
	t.Log("m3 address ", account.Address())
	chain := TestnetChain()
	token := NewTokenMain(chain)

	toAddress := M3Account(t).Address()
	amount := strconv.FormatUint(uint64(25e9), 10)
	// toAddress := account.Address()
	// amount := strconv.FormatUint(4e9, 10) // test big amount transfer

	txn, err := token.BuildTransferTransaction(account, toAddress, amount)
	require.Nil(t, err)

	simulateCheck(t, chain, &txn.Txn, false)
	// executeTransaction(t, chain, &txn.Txn, account.account)
}

func TestToken_Transfer_Use_Pay(t *testing.T) {
	ownerStr := "0x7d20dcdb2bca4f508ea9613994683eb4e76e9c4ed371169677c1be02aaf0b58e"
	recipientStr := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"
	chain := TestnetChain()

	coinCount := int(5)

	owner, err := sui_types.NewAddressFromHex(ownerStr)
	require.Nil(t, err)
	recipient, err := sui_types.NewAddressFromHex(recipientStr)
	require.Nil(t, err)
	cli, err := chain.Client()
	require.Nil(t, err)
	coins, err := cli.GetCoins(context.Background(), *owner, nil, nil, uint(coinCount))
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(coins.Data), coinCount)

	inputCoins := make([]sui_types.ObjectID, 0)
	totalAmount := uint64(0)
	for _, coin := range coins.Data {
		inputCoins = append(inputCoins, coin.CoinObjectId)
		totalAmount += coin.Balance.Uint64()
	}

	sendAmount := types.NewSafeSuiBigInt(totalAmount / 2)
	gasBudget := types.NewSafeSuiBigInt[uint64](30000000)
	txn, err := cli.Pay(context.Background(), *owner, inputCoins, []sui_types.SuiAddress{*recipient}, []types.SafeSuiBigInt[uint64]{sendAmount}, nil, gasBudget)
	require.Nil(t, err)

	simulateCheck(t, chain, txn, true)
}

func Test_TokenTransfer_Gas_Compare(t *testing.T) {
	ownerStr := "0x7d20dcdb2bca4f508ea9613994683eb4e76e9c4ed371169677c1be02aaf0b58e"
	recipientStr := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"
	chain := TestnetChain()

	owner, err := sui_types.NewAddressFromHex(ownerStr)
	require.Nil(t, err)
	recipient, err := sui_types.NewAddressFromHex(recipientStr)
	require.Nil(t, err)
	cli, err := chain.Client()
	require.Nil(t, err)
	coins, err := cli.GetCoins(context.Background(), *owner, nil, nil, 5)
	require.Nil(t, err)
	require.GreaterOrEqual(t, len(coins.Data), 1)

	sendCoin := coins.Data[0]
	sendAmount := types.NewSafeSuiBigInt(sendCoin.Balance.Uint64())
	gasBudget := types.NewSafeSuiBigInt[uint64](10000000)

	{
		txn, err := cli.Pay(context.Background(), *owner,
			[]sui_types.ObjectID{sendCoin.CoinObjectId},
			[]sui_types.SuiAddress{*recipient},
			[]types.SafeSuiBigInt[uint64]{sendAmount},
			nil, gasBudget)
		require.Nil(t, err)

		resp := simulateCheck(t, chain, txn, false)
		t.Log("gas use Pay = ", resp.Effects.Data.GasFee())
	}
	{
		txn, err := cli.TransferObject(context.Background(), *owner, *recipient,
			sendCoin.CoinObjectId, nil, gasBudget)
		require.Nil(t, err)

		resp := simulateCheck(t, chain, txn, false)
		t.Log("gas use TransferObject = ", resp.Effects.Data.GasFee())
	}
	{
		txn, err := cli.PaySui(context.Background(), *owner,
			[]sui_types.ObjectID{sendCoin.CoinObjectId},
			[]sui_types.SuiAddress{*recipient},
			[]types.SafeSuiBigInt[uint64]{sendAmount},
			gasBudget)
		require.Nil(t, err)

		resp, err := cli.DryRunTransaction(context.Background(), txn.TxBytes)
		require.Nil(t, err)
		require.False(t, resp.Effects.Data.IsSuccess())
		// InsufficientCoinBalance, because the sendCoin need balance=amount+gasfee
	}
}
