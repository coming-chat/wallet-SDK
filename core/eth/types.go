package eth

import (
	"errors"
	"math/big"
	"strings"

	HexType "github.com/centrifuge/go-substrate-rpc-client/v4/types"
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
func (msg *CallMsg) GetGas() int64       { return int64(msg.msg.Gas) }
func (msg *CallMsg) GetGasPrice() string { return msg.msg.GasPrice.String() }
func (msg *CallMsg) GetValue() string    { return msg.msg.Value.String() }
func (msg *CallMsg) GetData() []byte     { return msg.msg.Data }
func (msg *CallMsg) GetDataHex() string  { return HexType.HexEncodeToString(msg.msg.Data) }
func (msg *CallMsg) GetTo() string       { return msg.msg.To.String() }

func (msg *CallMsg) SetFrom(address string) { msg.msg.From = common.HexToAddress(address) }
func (msg *CallMsg) SetGas(gas int64)       { msg.msg.Gas = uint64(gas) }
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

type Transaction struct {
	Nonce    int64  // nonce of sender account
	GasPrice string // wei per gas
	GasLimit int64  // gas limit
	To       string // receiver
	Value    string // wei amount
	Data     []byte // contract invocation input data

	// EIP1559, Default is ""
	MaxPriorityFeePerGas string
}

func NewTransaction(nonce int64, gasPrice string, gasLimit int64, to string, value string, data []byte) (*Transaction, error) {
	tx := &Transaction{nonce, gasPrice, gasLimit, to, value, data, ""}
	_, err := tx.GetRawTx()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (tx *Transaction) GetRawTx() (*types.Transaction, error) {
	gasPriceInt, valid := big.NewInt(0).SetString(tx.GasPrice, 10)
	if !valid {
		return nil, errors.New("Invalid gasPrice")
	}
	valueInt, valid := big.NewInt(0).SetString(tx.Value, 10)
	if !valid {
		return nil, errors.New("Invalid value")
	}
	if !common.IsHexAddress(tx.To) {
		return nil, errors.New("Invalid toAddress")
	}
	toAddress := common.HexToAddress(tx.To)

	maxPriorityFeePerGasInt, valid := big.NewInt(0).SetString(tx.MaxPriorityFeePerGas, 10)
	if !valid || maxPriorityFeePerGasInt.Int64() == 0 {
		// is legacy tx
		return types.NewTx(&types.LegacyTx{
			Nonce:    uint64(tx.Nonce),
			To:       &toAddress,
			Value:    valueInt,
			Gas:      uint64(tx.GasLimit),
			GasPrice: gasPriceInt,
			Data:     tx.Data,
		}), nil
	} else {
		// is dynamic fee tx
		return types.NewTx(&types.DynamicFeeTx{
			Nonce:     uint64(tx.Nonce),
			To:        &toAddress,
			Value:     valueInt,
			Gas:       uint64(tx.GasLimit),
			GasFeeCap: gasPriceInt,
			GasTipCap: maxPriorityFeePerGasInt,
			Data:      tx.Data,
		}), nil
	}
}
