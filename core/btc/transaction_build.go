package btc

import (
	"errors"
	"strings"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/util/hexutil"
)

type Transaction struct {
	inputs    []input
	outputs   []output
	netParams *chaincfg.Params
}

type input struct {
	outPoint *wire.OutPoint
	prevOut  *wire.TxOut
}

type output *wire.TxOut

func NewTransaction(chainnet string) (*Transaction, error) {
	net, err := netParamsOf(chainnet)
	if err != nil {
		return nil, err
	}
	return &Transaction{netParams: net}, nil
}

func (t *Transaction) TotalInputValue() int64 {
	total := int64(0)
	for _, v := range t.inputs {
		total += v.prevOut.Value
	}
	return total
}

func (t *Transaction) TotalOutputValue() int64 {
	total := int64(0)
	for _, v := range t.outputs {
		total += v.Value
	}
	return total
}

func (t *Transaction) AddInput(txId string, index int64, address string, value int64) error {
	outPoint, err := outPoint(txId, uint32(index))
	if err != nil {
		return err
	}
	pkScript, err := addrToPkScript(address, t.netParams)
	if err != nil {
		return err
	}
	input := input{
		outPoint: outPoint,
		prevOut:  wire.NewTxOut(value, pkScript),
	}
	t.inputs = append(t.inputs, input)
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
	input := input{
		outPoint: outPoint,
		prevOut:  prevOut,
	}
	t.inputs = append(t.inputs, input)
	return nil
}

func (t *Transaction) AddOutput(address string, value int64) error {
	pkScript, err := addrToPkScript(address, t.netParams)
	if err != nil {
		return err
	}
	output := wire.NewTxOut(value, pkScript)
	t.outputs = append(t.outputs, output)
	return nil
}

func (t *Transaction) AddOpReturn(opReturn string) error {
	data := []byte(opReturn)
	script, err := buildOpReturnScript(data)
	if err != nil {
		return err
	}
	output := wire.NewTxOut(0, script)
	t.outputs = append(t.outputs, output)
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

func addrToPkScript(addr string, network *chaincfg.Params) ([]byte, error) {
	address, err := btcutil.DecodeAddress(addr, network)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(address)
}

func buildOpReturnScript(data []byte) ([]byte, error) {
	return txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData(data).Script()
}
