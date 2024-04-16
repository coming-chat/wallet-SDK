package btc

import (
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/blockchain"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
)

type TxOut struct {
	Hash  string `json:"hash,omitempty"`
	Index int64  `json:"index,omitempty"`

	Value   int64  `json:"value,omitempty"`
	Address string `json:"address,omitempty"`
}

type TxOutArray struct {
	inter.AnyArray[*TxOut]
}

func (ta *TxOutArray) detailDesc() string {
	desc := ""
	for idx, out := range ta.AnyArray {
		if idx != 0 {
			desc += "\n"
		}
		if out.Address != "" {
			addr := out.Address
			if len(addr) > 13 {
				addr = addr[:5] + "..." + addr[len(addr)-5:]
			}
			valueBTC := btcutil.Amount(out.Value).String()
			desc += "\t" + addr + "\t" + valueBTC
		} else {
			desc += "\thash: " + out.Hash + "\n" + "\tindex: " + strconv.FormatInt(out.Index, 10)
		}
	}
	return desc
}

type TransactionDetail struct {
	TotalCost  int64       `json:"totalCost"`
	NetworkFee int64       `json:"networkFee"`
	FeeRate    float64     `json:"feeRate"`
	Inputs     *TxOutArray `json:"inputs"`
	Outputs    *TxOutArray `json:"outputs"`
}

func (td *TransactionDetail) JsonString() (*base.OptionalString, error) {
	return base.JsonString(td)
}

func (td *TransactionDetail) Desc() string {
	totalBTC := "Unknow"
	feeBTC := "Unknow"
	feeRate := "Unknow"
	if td.NetworkFee != 0 {
		totalBTC = btcutil.Amount(td.TotalCost).String()
		feeBTC = btcutil.Amount(td.NetworkFee).String()
		feeRate = fmt.Sprintf("%.2f sat/vB", td.FeeRate)
	}

	return fmt.Sprintf(`Total Cost:
	%v
Network Fee:
	%v
Network FeeRate:
	%v
Inputs:
%v
Outputs:
%v`, totalBTC, feeBTC, feeRate, td.Inputs.detailDesc(), td.Outputs.detailDesc())
}

func DecodePsbtTransactionDetail(psbtHex string, chainnet string) (d *TransactionDetail, err error) {
	chainParams, err := netParamsOf(chainnet)
	if err != nil {
		return
	}
	packet, err := DecodePsbtTxToPacket(psbtHex)
	if err != nil {
		return
	}
	txFee, err := packet.GetTxFee()
	if err != nil {
		return
	}
	feeFloat := txFee.ToUnit(btcutil.AmountSatoshi)
	vSize := virtualSize(packet.UnsignedTx)
	feeRate := feeFloat / float64(vSize)

	inputs := make([]*TxOut, len(packet.Inputs))
	for idx, in := range packet.Inputs {
		switch {
		case in.WitnessUtxo != nil:
			inputs[idx], err = txOutFromWireTxOut(in.WitnessUtxo, chainParams)
			if err != nil {
				return
			}
		case in.NonWitnessUtxo != nil:
			utxOuts := in.NonWitnessUtxo.TxOut
			txIn := packet.UnsignedTx.TxIn[idx]
			opIdx := txIn.PreviousOutPoint.Index
			txOut := utxOuts[opIdx]
			inputs[idx], err = txOutFromWireTxOut(txOut, chainParams)
			if err != nil {
				return
			}
		default:
			return nil, fmt.Errorf("input %d has no UTXO information",
				idx)
		}
	}

	outputs := make([]*TxOut, len(packet.UnsignedTx.TxOut))
	for idx, out := range packet.UnsignedTx.TxOut {
		outputs[idx], err = txOutFromWireTxOut(out, chainParams)
		if err != nil {
			return
		}
	}

	totalCost := int64(0)
	for _, input := range inputs {
		totalCost += input.Value
	}
	for _, output := range outputs {
		if output.Address == inputs[0].Address {
			totalCost -= output.Value
		}
	}

	return &TransactionDetail{
		TotalCost:  totalCost,
		NetworkFee: int64(feeFloat),
		FeeRate:    float64(feeRate),
		Inputs:     &TxOutArray{inputs},
		Outputs:    &TxOutArray{outputs},
	}, nil
}

func DecodeTxHexTransactionDetail(txHex string, chainnet string) (detail *TransactionDetail, err error) {
	chainParams, err := netParamsOf(chainnet)
	if err != nil {
		return
	}
	msgTx, err := DecodeTx(txHex)
	if err != nil {
		return nil, err
	}
	// msgTx cannot get `fee`, `feeRate`

	inputs := make([]*TxOut, len(msgTx.TxIn))
	for idx, in := range msgTx.TxIn {
		inputs[idx] = &TxOut{
			Hash:  in.PreviousOutPoint.Hash.String(),
			Index: int64(in.PreviousOutPoint.Index),
		}
	}
	outputs := make([]*TxOut, len(msgTx.TxOut))
	for idx, out := range msgTx.TxOut {
		outputs[idx], err = txOutFromWireTxOut(out, chainParams)
		if err != nil {
			return
		}
	}

	return &TransactionDetail{
		TotalCost:  0,
		NetworkFee: 0,
		FeeRate:    0,
		Inputs:     &TxOutArray{inputs},
		Outputs:    &TxOutArray{outputs},
	}, nil
}

func txOutFromWireTxOut(txout *wire.TxOut, params *chaincfg.Params) (*TxOut, error) {
	pkobj, err := txscript.ParsePkScript(txout.PkScript)
	if err != nil {
		return nil, err
	}
	addr, err := pkobj.Address(params)
	if err != nil {
		return nil, err
	}
	return &TxOut{
		Value:   txout.Value,
		Address: addr.EncodeAddress(),
	}, nil
}

// calculation reference:
// https://github.com/btcsuite/btcd/blob/569155bc6a502f45b4a514bc6b9d5f814a980b6c/mempool/policy.go#L382
func virtualSize(tx *wire.MsgTx) int64 {
	// vSize := (((baseSize * 3) + totalSize) + 3) / 4
	baseSize := int64(tx.SerializeSizeStripped())
	totalSize := int64(tx.SerializeSize())
	weight := (baseSize * (blockchain.WitnessScaleFactor - 1)) + totalSize
	return (weight + blockchain.WitnessScaleFactor - 1) / blockchain.WitnessScaleFactor
}
