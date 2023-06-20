package starknet

import (
	"math/big"

	"github.com/NethermindEth/juno/core"
	"github.com/NethermindEth/juno/core/felt"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func deployAccountTxnForArgentX(pubKey string) (*core.DeployAccountTransaction, error) {
	pubData, err := new(felt.Felt).SetString(pubKey)
	if err != nil {
		return nil, base.ErrInvalidPublicKey
	}
	classHash, _ := new(felt.Felt).SetString("0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918")
	data1, _ := new(felt.Felt).SetString("0x33434ad846cdd5f23eb73ff09fe6fddd568284a0fb7d1be20ee482f044dabe2")
	data2, _ := new(felt.Felt).SetString("0x79dc0da7c54b95f10aa182ad0a46400db63156920adb65eca2654c0945a463") // getSelectorFromName("initialize")
	data3, _ := new(felt.Felt).SetString("0x2")
	data4 := pubData
	data5, _ := new(felt.Felt).SetString("0x0")
	txn := &core.DeployAccountTransaction{
		DeployTransaction: core.DeployTransaction{
			ClassHash:           classHash,
			ContractAddressSalt: pubData,
			ConstructorCallData: []*felt.Felt{
				data1, data2, data3, data4, data5,
			},
			Version: new(felt.Felt).SetUint64(1),
		},
		MaxFee: new(felt.Felt).SetBytes(big.NewInt(1e13).Bytes()),
		Nonce:  new(felt.Felt).SetUint64(0),
	}

	callerAddress, _ := new(felt.Felt).SetString("0x0000000000000000000000000000000000000000")
	address := core.ContractAddress(callerAddress, txn.ClassHash, txn.ContractAddressSalt, txn.ConstructorCallData)
	txn.ContractAddress = address

	return txn, nil
}
