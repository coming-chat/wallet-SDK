package eth

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"strings"

	HexType "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type UrlParam struct {
	RpcUrl string
	WsUrl  string
}

type CallMethodOpts struct {
	Nonce                int64
	Value                string
	GasPrice             string // MaxFeePerGas
	GasLimit             string
	IsPredictError       bool
	MaxPriorityFeePerGas string
}

type CallMethodOptsBigInt struct {
	Nonce                uint64
	Value                *big.Int
	GasPrice             *big.Int // MaxFeePerGas
	GasLimit             uint64
	IsPredictError       bool
	MaxPriorityFeePerGas *big.Int
}

type BuildTxResult struct {
	SignedTx *types.Transaction
	TxHex    string
}

type TransactionByHashResult struct {
	SignedTx    *types.Transaction
	From        common.Address
	IsPending   bool   // 交易是否处于Pending状态
	Status      string // 0: 交易失败, 1: 交易成功
	GasUsed     string // 实际花费gas
	BlockNumber string // 区块高度
}

type Erc20TxParams struct {
	ToAddress string `json:"toAddress"`
	Amount    string `json:"amount"`
	Method    string `json:"method"`
}

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	msg ethereum.CallMsg
}

// NewCallMsg creates an empty contract call parameter list.
func NewCallMsg() *CallMsg {
	return new(CallMsg)
}

func (msg *CallMsg) GetFrom() string     { return msg.msg.From.String() }
func (msg *CallMsg) GetGasLimit() string { return strconv.FormatUint(msg.msg.Gas, 10) }
func (msg *CallMsg) GetGasPrice() string { return msg.msg.GasPrice.String() }
func (msg *CallMsg) GetValue() string    { return msg.msg.Value.String() }
func (msg *CallMsg) GetData() []byte     { return msg.msg.Data }
func (msg *CallMsg) GetDataHex() string  { return HexType.HexEncodeToString(msg.msg.Data) }
func (msg *CallMsg) GetTo() string       { return msg.msg.To.String() }

func (msg *CallMsg) SetFrom(address string) { msg.msg.From = common.HexToAddress(address) }
func (msg *CallMsg) SetGasLimit(gas string) {
	i, _ := strconv.ParseUint(gas, 10, 64)
	msg.msg.Gas = i
}
func (msg *CallMsg) SetGasPrice(price string) {
	i, _ := new(big.Int).SetString(price, 10)
	msg.msg.GasPrice = i
}

// Set amount with decimal number
func (msg *CallMsg) SetValue(value string) {
	i, _ := new(big.Int).SetString(value, 10)
	msg.msg.Value = i
}

// Set amount with hexadecimal number
func (msg *CallMsg) SetValueHex(hex string) {
	hex = strings.TrimPrefix(hex, "0x") // must trim 0x !!
	i, _ := new(big.Int).SetString(hex, 16)
	msg.msg.Value = i
}
func (msg *CallMsg) SetData(data []byte) { msg.msg.Data = common.CopyBytes(data) }
func (msg *CallMsg) SetDataHex(hex string) {
	data, err := HexType.HexDecodeString(hex)
	if err != nil {
		return
	}
	msg.msg.Data = data
}
func (msg *CallMsg) SetTo(address string) {
	if address == "" {
		msg.msg.To = nil
	} else {
		a := common.HexToAddress(address)
		msg.msg.To = &a
	}
}

func (msg *CallMsg) TransferToTransaction() *Transaction {
	return &Transaction{
		GasPrice: msg.GetGasPrice(),
		GasLimit: msg.GetGasLimit(),
		To:       msg.GetTo(),
		Value:    msg.GetValue(),
		Data:     msg.GetDataHex(),
	}
}

type Transaction struct {
	Nonce    string // nonce of sender account
	GasPrice string // wei per gas
	GasLimit string // gas limit
	To       string // receiver
	Value    string // wei amount
	Data     string // contract invocation input data

	// EIP1559, Default is ""
	MaxPriorityFeePerGas string
}

func NewTransaction(nonce, gasPrice, gasLimit, to, value, data string) *Transaction {
	return &Transaction{nonce, gasPrice, gasLimit, to, value, data, ""}
}

func NewTransactionFromHex(hexData string) (*Transaction, error) {
	rawBytes, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	decodeTx := types.NewTx(&types.DynamicFeeTx{})
	err = decodeTx.UnmarshalBinary(rawBytes)
	if err != nil {
		return nil, err
	}
	tx := NewTransaction(
		strconv.Itoa(int(decodeTx.Nonce())),
		decodeTx.GasFeeCap().String(),
		strconv.Itoa(int(decodeTx.Gas())),
		decodeTx.To().String(),
		decodeTx.Value().String(),
		hex.EncodeToString(decodeTx.Data()))
	// not equal, is eip1559; legacy feecap equal tipcap
	if decodeTx.GasTipCap().Cmp(decodeTx.GasFeeCap()) != 0 {
		tx.MaxPriorityFeePerGas = decodeTx.GasTipCap().String()
	}
	return tx, nil
}

// This is an alias property for GasPrice in order to support EIP1559
func (tx *Transaction) MaxFee() string {
	return tx.GasPrice
}

// This is an alias property for GasPrice in order to support EIP1559
func (tx *Transaction) SetMaxFee(maxFee string) {
	tx.GasPrice = maxFee
}

func (tx *Transaction) GetRawTx() (*types.Transaction, error) {
	var (
		gasPrice, value, maxFeePerGas *big.Int // default nil

		nonce     uint64 = 0
		gasLimit  uint64 = 90000 // reference https://eth.wiki/json-rpc/API method eth_sendTransaction
		toAddress common.Address
		data      []byte
		valid     bool
		err       error
	)
	if tx.GasPrice != "" {
		if gasPrice, valid = big.NewInt(0).SetString(tx.GasPrice, 10); !valid {
			return nil, errors.New("Invalid gasPrice")
		}
	}
	if tx.Value != "" {
		if value, valid = big.NewInt(0).SetString(tx.Value, 10); !valid {
			return nil, errors.New("Invalid value")
		}
	}
	if tx.MaxPriorityFeePerGas != "" {
		if maxFeePerGas, valid = big.NewInt(0).SetString(tx.MaxPriorityFeePerGas, 10); !valid {
			return nil, errors.New("Invalid max priority fee per gas")
		}
	}
	if tx.Nonce != "" {
		if nonce, err = strconv.ParseUint(tx.Nonce, 10, 64); err != nil {
			return nil, errors.New("Invalid Nonce")
		}
	}
	if tx.GasLimit != "" {
		if gasLimit, err = strconv.ParseUint(tx.GasLimit, 10, 64); err != nil {
			return nil, errors.New("Invalid gas limit")
		}
	}
	if tx.To != "" && !common.IsHexAddress(tx.To) {
		return nil, errors.New("Invalid toAddress")
	}
	toAddress = common.HexToAddress(tx.To)
	if tx.Data != "" {
		if data, err = HexType.HexDecodeString(tx.Data); err != nil {
			return nil, errors.New("Invalid data string")
		}
	}

	if maxFeePerGas == nil || maxFeePerGas.Int64() == 0 {
		// is legacy tx
		return types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &toAddress,
			Value:    value,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     data,
		}), nil
	} else {
		// is dynamic fee tx
		return types.NewTx(&types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &toAddress,
			Value:     value,
			Gas:       gasLimit,
			GasFeeCap: gasPrice,
			GasTipCap: maxFeePerGas,
			Data:      data,
		}), nil
	}
}

func (tx *Transaction) TransformToErc20Transaction(contractAddress string) error {
	if len(tx.Data) > 0 && tx.Value == "0" {
		return nil
	}
	data, err := EncodeErc20Transfer(tx.To, tx.Value)
	if err != nil {
		return err
	}

	tx.To = contractAddress
	tx.Value = "0"
	tx.Data = HexType.HexEncodeToString(data)
	return nil
}

// @return gasPrice * gasLimit + value
func (tx *Transaction) TotalAmount() string {
	priceInt, ok := big.NewInt(0).SetString(tx.GasPrice, 10)
	if !ok {
		return "0"
	}
	limitInt, ok := big.NewInt(0).SetString(tx.GasLimit, 10)
	if !ok {
		return "0"
	}
	amount, ok := big.NewInt(0).SetString(tx.Value, 10)
	if !ok {
		return "0"
	}
	return amount.Add(amount, priceInt.Mul(priceInt, limitInt)).String()
}

func (t *Transaction) SignWithAccount(account base.Account) (signedTx *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}
