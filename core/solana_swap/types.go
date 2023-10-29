package solanaswap

import (
	"encoding/binary"
	"errors"
	"math/big"
)

func LittleEndianSerializeBytesWithLength(bytes []byte) []byte {
	res := make([]byte, 0, len(bytes))
	res = binary.LittleEndian.AppendUint64(res, uint64(len(bytes)))
	res = append(res, bytes...)
	return res
}

func LittleEndianSerializeU256(num *big.Int) ([]byte, error) {
	cop := big.NewInt(0).Set(num)
	if cop.Sign() < 0 {
		return nil, errors.New("invalid U256: negative number")
	}
	if cop.BitLen() > 256 {
		return nil, errors.New("invalid U256: too large number")
	}

	data := make([]byte, 0, 256/8)
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
	rdata := make([]byte, 0, 1024)
	data := &rdata
	serialize_vector_with_length(data, so.TransactionId)
	serialize_vector_with_length(data, so.Receiver)
	serialize_u16(data, so.SourceChainId)
	serialize_vector_with_length(data, so.SendingAssetId)
	serialize_u16(data, so.DestinationChainId)
	serialize_vector_with_length(data, so.ReceivingAssetId)
	serialize_u256(data, *so.Amount)
	return rdata, nil
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
	rdata := make([]byte, 0, 1024)
	data := &rdata
	serialize_u16(data, self.DstWormholeChainId)
	serialize_u256(data, *big.NewInt(0).SetUint64(self.DstMaxGasPriceInWeiForRelayer))
	serialize_u256(data, *big.NewInt(0).SetUint64(self.WormholeFee))
	serialize_vector_with_length(data, self.DstSoDiamond)
	return rdata, nil
}
