package btc

import (
	"bytes"
	"encoding/hex"

	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type PsbtTransaction struct {
	Packet psbt.Packet
}

func NewPsbtTransaction(psbtString string) (*PsbtTransaction, error) {
	packet, err := DecodePsbtTxToPacket(psbtString)
	if err != nil {
		return nil, err
	}
	return &PsbtTransaction{*packet}, nil
}

func (t *PsbtTransaction) SignWithAccount(account base.Account) (signedTxn *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}
func (t *PsbtTransaction) SignedTransactionWithAccount(account base.Account) (signedTxn base.SignedTransaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	btcAccount := account.(*Account)
	if btcAccount == nil {
		return nil, base.ErrInvalidAccountType
	}
	if err = SignPSBTTx(&t.Packet, btcAccount); err != nil {
		return
	}
	if err = EnsurePsbtFinalize(&t.Packet); err != nil {
		return
	}
	return &SignedPsbtTransaction{t.Packet}, nil
}

type SignedPsbtTransaction struct {
	Packet psbt.Packet
}

func (t *SignedPsbtTransaction) HexString() (res *base.OptionalString, err error) {
	return nil, base.ErrUnsupportedFunction
}

func (t *SignedPsbtTransaction) PsbtHexString() (*base.OptionalString, error) {
	packet := t.Packet
	if err := EnsurePsbtFinalize(&packet); err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	if err := packet.Serialize(&buff); err != nil {
		return nil, err
	}
	hexString := hex.EncodeToString(buff.Bytes())
	return &base.OptionalString{Value: hexString}, nil
}

func (t *SignedPsbtTransaction) PublishWithChain(c *Chain) (hashs *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	msgTx, err := PsbtPacketToMsgTx(&t.Packet)
	if err != nil {
		return nil, err
	}
	cli, err := rpcClientOf(c.Chainnet)
	if err != nil {
		return nil, err
	}
	hash, err := cli.SendRawTransaction(msgTx, false)
	if err != nil {
		return nil, err
	}
	return &base.OptionalString{Value: hash.String()}, nil
}
