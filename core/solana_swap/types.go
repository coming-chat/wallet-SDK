package solanaswap

import (
	"math/big"
)

type SoData struct {
	TransactionId      []byte
	Receiver           []byte
	SourceChainId      uint16
	SendingAssetId     []byte
	DestinationChainId uint16
	ReceivingAssetId   []byte
	Amount             *big.Int
}

func (so *SoData) Serialize() []byte {
	sl := NewEthSerializer(1024)
	sl.AppendBytesWithLenth(so.TransactionId)
	sl.AppendBytesWithLenth(so.Receiver)
	sl.AppendU16(so.DestinationChainId)
	sl.AppendBytesWithLenth(so.SendingAssetId)
	sl.AppendU16(so.DestinationChainId)
	sl.AppendBytesWithLenth(so.ReceivingAssetId)
	sl.AppendU256(*so.Amount)
	return sl.Bytes()
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

func (self *WormholeData) Serialize() []byte {
	s := NewEthSerializer(1024)
	s.AppendU16(self.DstWormholeChainId)
	s.AppendU256(*big.NewInt(0).SetUint64(self.DstMaxGasPriceInWeiForRelayer))
	s.AppendU256(*big.NewInt(0).SetUint64(self.WormholeFee))
	s.AppendBytesWithLenth(self.DstSoDiamond)
	return s.Bytes()
}
