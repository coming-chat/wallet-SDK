package solanaswap

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/blocto/solana-go-sdk/pkg/bincode"
)

func SerializeBytesWithLength(bytes []byte) []byte {
	res := make([]byte, 0, len(bytes))
	res = binary.LittleEndian.AppendUint64(res, uint64(len(bytes)))
	res = append(res, bytes...)
	return res
}

func SerializeU256(num *big.Int) ([]byte, error) {
	cop := big.NewInt(0).Set(num)
	if cop.Sign() < 0 {
		return nil, errors.New("invalid U256: negative number")
	}
	if cop.BitLen() > 256 {
		return nil, errors.New("invalid U256: too large number")
	}

	data := make([]byte, 0)
	for t := 0; t < 256/64; t++ {
		data = binary.LittleEndian.AppendUint64(data, cop.Uint64())
		cop = cop.Rsh(cop, 64)
	}
	return data, nil
}

type SoData struct {
	TransactionId      []byte
	Receiver           []byte
	SourceChainId      uint16
	SendingAssetId     []byte
	DestinationChainId uint16
	ReceivingAssetId   []byte
	Amount             *big.Int
}

func (so *SoData) Serialize() ([]byte, error) {
	data := make([]byte, 0, 1024)

	data = append(data, SerializeBytesWithLength(so.TransactionId)...)
	data = append(data, SerializeBytesWithLength(so.Receiver)...)
	d, err := bincode.SerializeData(so.SourceChainId)
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	data = append(data, SerializeBytesWithLength(so.SendingAssetId)...)
	d, err = bincode.SerializeData(so.DestinationChainId)
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	data = append(data, SerializeBytesWithLength(so.ReceivingAssetId)...)
	d, err = SerializeU256(so.Amount)
	if err != nil {
		return nil, err
	}
	data = append(data, d...)
	return data, nil
}

type SwapData struct {
	CallTo           []byte
	ApproveTo        []byte
	SendingAssetId   []byte
	ReceivingAssetId []byte
	FromAmount       int
	CallData         []byte
	SwapType         string
	SwapFuncName     string
	SwapPath         []any
}

func (swap *SwapData) Serialize() []byte {
	return nil
}

type WormholeData struct {
	DstWormholeChainId            uint16
	DstMaxGasPriceInWeiForRelayer uint64
	WormholeFee                   uint64
	DstSoDiamond                  []byte
}

func (self *WormholeData) Serialize() ([]byte, error) {
	data := make([]byte, 0, 1024)

	d, err := bincode.SerializeData(self.DstWormholeChainId)
	if err != nil {
		return nil, err
	}
	data = append(data, d...)

	d, _ = SerializeU256(big.NewInt(0).SetUint64(self.DstMaxGasPriceInWeiForRelayer))
	data = append(data, d...)

	d, _ = SerializeU256(big.NewInt(0).SetUint64(self.WormholeFee))
	data = append(data, d...)

	data = append(data, SerializeBytesWithLength(self.DstSoDiamond)...)

	return data, nil
}
