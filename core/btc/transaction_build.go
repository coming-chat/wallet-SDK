package btc

import (
	"errors"
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/util/hexutil"
)

const (
	MaxOpReturnLength = 75
)

type Transaction struct {
	netParams      *chaincfg.Params
	msgTx          *wire.MsgTx
	prevOutFetcher *txscript.MultiPrevOutFetcher
}

func NewTransaction(chainnet string) (*Transaction, error) {
	net, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}
	tx := wire.NewMsgTx(wire.TxVersion)
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)
	return &Transaction{
		netParams:      net,
		msgTx:          tx,
		prevOutFetcher: prevOutFetcher,
	}, nil
}

func (t *Transaction) TotalInputValue() int64 {
	total := int64(0)
	for _, v := range t.msgTx.TxIn {
		out := t.prevOutFetcher.FetchPrevOutput(v.PreviousOutPoint)
		total += out.Value
	}
	return total
}

func (t *Transaction) TotalOutputValue() int64 {
	total := int64(0)
	for _, v := range t.msgTx.TxOut {
		total += v.Value
	}
	return total
}

func (t *Transaction) EstimateTransactionSize() int64 {
	return virtualSize(ensureSignOrFakeSign(t.msgTx, t.prevOutFetcher))
}

func (t *Transaction) AddInput(txId string, index int64, address string, value int64) error {
	outPoint, err := outPoint(txId, uint32(index))
	if err != nil {
		return err
	}
	pkScript, err := addressToPkScript(address, t.netParams)
	if err != nil {
		return err
	}
	txIn := wire.NewTxIn(outPoint, nil, nil)
	prevOut := wire.NewTxOut(value, pkScript)
	t.msgTx.TxIn = append(t.msgTx.TxIn, txIn)
	t.prevOutFetcher.AddPrevOut(*outPoint, prevOut)
	return nil
}

func (t *Transaction) AddInput2(txId string, index int64, prevTx string) error {
	outPoint, err := outPoint(txId, uint32(index))
	if err != nil {
		return err
	}
	prevOut, err := prevTxOut(prevTx, uint32(index))
	if err != nil {
		return err
	}
	txIn := wire.NewTxIn(outPoint, nil, nil)
	t.msgTx.TxIn = append(t.msgTx.TxIn, txIn)
	t.prevOutFetcher.AddPrevOut(*outPoint, prevOut)
	return nil
}

// If the value is 0, `AddOpReturn` will be called.
func (t *Transaction) AddOutput(address string, value int64) error {
	if value == 0 {
		return t.AddOpReturn(address)
	}
	pkScript, err := addressToPkScript(address, t.netParams)
	if err != nil {
		return err
	}
	output := wire.NewTxOut(value, pkScript)
	t.msgTx.TxOut = append(t.msgTx.TxOut, output)
	return nil
}

func (t *Transaction) AddOpReturn(opReturn string) error {
	data := []byte(opReturn)
	script, err := buildOpReturnScript(data)
	if err != nil {
		return err
	}
	output := wire.NewTxOut(0, script)
	t.msgTx.TxOut = append(t.msgTx.TxOut, output)
	return nil
}

func outPoint(txId string, index uint32) (*wire.OutPoint, error) {
	txId = strings.TrimPrefix(txId, "0x")
	txHash, err := chainhash.NewHashFromStr(txId)
	if err != nil {
		return nil, err
	}
	return wire.NewOutPoint(txHash, index), nil
}

func prevTxOut(preTx string, index uint32) (*wire.TxOut, error) {
	txData, err := hexutil.HexDecodeString(preTx)
	if err != nil {
		return nil, err
	}
	tx, err := btcutil.NewTxFromBytes(txData)
	if err != nil {
		return nil, err
	}
	outs := tx.MsgTx().TxOut
	if len(outs) > int(index) {
		return outs[index], nil
	} else {
		return nil, errors.New("invalid output index")
	}
}

func addressToPkScript(addr string, network *chaincfg.Params) ([]byte, error) {
	address, err := btcutil.DecodeAddress(addr, network)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(address)
}

func buildOpReturnScript(data []byte) ([]byte, error) {
	if len(data) > MaxOpReturnLength {
		return nil, errors.New("op return length cannot be greater than " + strconv.FormatInt(MaxOpReturnLength, 10))
	}
	return txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(data).Script()
}
