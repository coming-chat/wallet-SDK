package cosmos

import (
	"context"
	"errors"
	"math/big"
	"strconv"

	hexTypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

type Token struct {
	chain  *Chain
	Prefix string
	Denom  string
}

func NewToken(chain *Chain, prefix, denom string) *Token {
	return &Token{
		chain:  chain,
		Prefix: prefix,
		Denom:  denom,
	}
}

// MARK - Implement the protocol Token

func (t *Token) Chain() base.Chain {
	return t.chain
}

// Warning: Main token does not support
func (t *Token) TokenInfo() (*base.TokenInfo, error) {
	return nil, errors.New("Cosmos token does not support")
}

func (t *Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.chain.BalanceOfAddressAndDenom(address, t.Denom)
}

// Warning: Unable to use public key to query balance
func (t *Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return base.EmptyBalance(), errors.New("Unable to use public key to query balance")
}

func (t *Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
}

// MARK - Cosmos Token

func (t *Token) BuildTransferTx(privateKey, receiverAddress, gasPrice, gasLimit, amount, memo string) (string, error) {
	priBytes, err := hexTypes.HexDecodeString(privateKey)
	if err != nil {
		return "", err
	}
	priKey := &secp256k1.PrivKey{Key: priBytes}
	return t.buildTransferTx(priKey, receiverAddress, gasPrice, gasLimit, amount, memo)
}

func (t *Token) BuildTransferTxWithAccount(account *Account, receiverAddress, gasPrice, gasLimit, amount, memo string) (string, error) {
	return t.buildTransferTx(account.privKey, receiverAddress, gasPrice, gasLimit, amount, memo)
}

func (t *Token) buildTransferTx(privateKey types.PrivKey, receiverAddress, gasPrice, gasLimit, amount, memo string) (s string, err error) {
	amountInt, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		return
	}
	amountCoin := sdk.NewInt64Coin(t.Denom, amountInt)

	gasPriceFloat, b := new(big.Float).SetString(gasPrice)
	if b == false {
		return s, errors.New("Gas Price format error")
	}
	gasLimitInt, b := new(big.Int).SetString(gasLimit, 10)
	if b == false {
		return s, errors.New("Gas Limit format error")
	}
	gasFloat := new(big.Float).Mul(gasPriceFloat, new(big.Float).SetInt(gasLimitInt))
	gasInt, _ := gasFloat.Int64()
	feeCoin := sdk.NewInt64Coin(t.Denom, gasInt)

	fromAddress := privateKey.PubKey().Address()
	addressString, err := Bech32FromAccAddress(fromAddress.Bytes(), t.Prefix)
	if err != nil {
		return
	}

	client, err := t.chain.GetClient()
	if err != nil {
		return
	}

	accountInfo, err := t.chain.AccountOf(addressString)
	if err != nil {
		return
	}
	accountNumber, err := strconv.ParseUint(accountInfo.AccountNumber, 10, 64)
	if err != nil {
		return
	}
	sequence, err := strconv.ParseUint(accountInfo.Sequence, 10, 64)
	if err != nil {
		return
	}

	blockInfo, err := client.Block(context.Background(), nil)
	if err != nil {
		return
	}
	chainId := blockInfo.Block.ChainID
	latestHeight := blockInfo.Block.Height

	// Create a new TxBuilder.
	encCfg := simapp.MakeTestEncodingConfig()
	txBuilder := encCfg.TxConfig.NewTxBuilder()
	msg1 := &banktypes.MsgSend{
		FromAddress: addressString,
		ToAddress:   receiverAddress,
		Amount:      sdk.NewCoins(amountCoin),
	}
	err = txBuilder.SetMsgs(msg1)
	if err != nil {
		return
	}

	txBuilder.SetGasLimit(gasLimitInt.Uint64())
	txBuilder.SetFeeAmount(sdk.NewCoins(feeCoin))
	txBuilder.SetMemo(memo)
	txBuilder.SetTimeoutHeight(uint64(latestHeight) + 1000)

	sigV2 := signing.SignatureV2{
		PubKey: privateKey.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: sequence,
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return
	}

	signerData := xauthsigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accountNumber,
		Sequence:      sequence,
	}
	sigV2, err = clientTx.SignWithPrivKey(
		encCfg.TxConfig.SignModeHandler().DefaultMode(), signerData,
		txBuilder, privateKey, encCfg.TxConfig, sequence)
	if err != nil {
		return
	}
	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return
	}

	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return
	}

	signedTx := hexTypes.HexEncodeToString(txBytes)
	return signedTx, nil
}

func (t *Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (t *Token) CanTransferAll() bool {
	return false
}
func (t *Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}
