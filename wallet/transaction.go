package wallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"math/big"
	"wallet-SDK/chainxTypes"
)

type Tx struct {
	metadata *types.Metadata
}

func NewTx(metadataString string) (*Tx, error) {
	var metadata types.Metadata
	if err := types.DecodeFromHexString(metadataString, &metadata); err != nil {
		return nil, err
	}
	return &Tx{
		metadata: &metadata,
	}, nil
}

type Transaction struct {
	extrinsic       *types.Extrinsic
	extrinsicChainX *chainxTypes.Extrinsic
}

func (t *Transaction) GetSignData(genesisHashString string, nonce int64, specVersion, transVersion int32) ([]byte, error) {
	var methodBytes []byte
	genesisHash, err := types.NewHashFromHexString(genesisHashString)
	if err != nil {
		return nil, err
	}
	if t.extrinsicChainX != nil {
		t.extrinsicChainX.Signature = chainxTypes.ExtrinsicSignatureV4{
			Nonce: types.NewUCompactFromUInt(uint64(nonce)),
			Era:   types.ExtrinsicEra{IsImmortalEra: true},
			Tip:   types.NewUCompactFromUInt(0),
		}
		methodBytes, err = types.EncodeToBytes(t.extrinsicChainX.Method)
		if err != nil {
			return nil, err
		}
		return types.EncodeToBytes(types.ExtrinsicPayloadV4{
			ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
				Method:      methodBytes,
				Era:         t.extrinsicChainX.Signature.Era,
				Nonce:       t.extrinsicChainX.Signature.Nonce,
				Tip:         t.extrinsicChainX.Signature.Tip,
				SpecVersion: types.NewU32(uint32(specVersion)),
				GenesisHash: genesisHash,
				BlockHash:   genesisHash,
			},
			TransactionVersion: types.NewU32(uint32(transVersion)),
		})
	} else if t.extrinsic != nil {
		t.extrinsic.Signature = types.ExtrinsicSignatureV4{
			Nonce: types.NewUCompactFromUInt(uint64(nonce)),
			Era:   types.ExtrinsicEra{IsImmortalEra: true},
			Tip:   types.NewUCompactFromUInt(0),
		}
		methodBytes, err = types.EncodeToBytes(t.extrinsic.Method)
		if err != nil {
			return nil, err
		}
		return types.EncodeToBytes(types.ExtrinsicPayloadV4{
			ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
				Method:      methodBytes,
				Era:         t.extrinsic.Signature.Era,
				Nonce:       t.extrinsic.Signature.Nonce,
				Tip:         t.extrinsic.Signature.Tip,
				SpecVersion: types.NewU32(uint32(specVersion)),
				GenesisHash: genesisHash,
				BlockHash:   genesisHash,
			},
			TransactionVersion: types.NewU32(uint32(transVersion)),
		})
	}
	return nil, ErrNilExtrinsic
}

func (t *Transaction) GetUnSignTx() (string, error) {
	if t.extrinsicChainX != nil {
		return types.EncodeToHexString(t.extrinsicChainX)
	} else if t.extrinsic != nil {
		return types.EncodeToHexString(t.extrinsic)
	}
	return "", ErrNilExtrinsic
}

func (t *Transaction) GetTxFromHex(signerPublicKeyHex string, signatureDataHex string) (string, error) {
	signerPublicKey, err := types.HexDecodeString(signerPublicKeyHex)
	if err != nil {
		return "", err
	}
	signatureData, err := types.HexDecodeString(signatureDataHex)
	if err != nil {
		return "", err
	}
	return t.GetTx(signerPublicKey, signatureData)
}

func (t *Transaction) GetTx(signerPublicKey []byte, signatureData []byte) (string, error) {
	if signatureData == nil {
		return "", ErrNotSigned
	}

	if signerPublicKey == nil {
		return "", ErrNoPublicKey
	}

	if t.extrinsicChainX != nil {
		t.extrinsicChainX.Signature.Signer = types.NewAddressFromAccountID(signerPublicKey)
		t.extrinsicChainX.Signature.Signature = types.MultiSignature{IsSr25519: true, AsSr25519: types.NewSignature(signatureData)}
		t.extrinsicChainX.Version |= types.ExtrinsicBitSigned
		return types.EncodeToHexString(t.extrinsicChainX)
	} else {
		t.extrinsic.Signature.Signer = types.NewMultiAddressFromAccountID(signerPublicKey)
		t.extrinsic.Signature.Signature = types.MultiSignature{IsSr25519: true, AsSr25519: types.NewSignature(signatureData)}
		t.extrinsic.Version |= types.ExtrinsicBitSigned
		return types.EncodeToHexString(t.extrinsic)
	}
}

func (t *Tx) NewTransactionFromHex(isChainX bool, txHex string) (*Transaction, error) {
	var (
		transaction = &Transaction{}
	)

	if t.metadata == nil {
		return nil, ErrNilMetadata
	}

	if isChainX {
		transaction.extrinsicChainX = &chainxTypes.Extrinsic{}
		err := types.DecodeFromHexString(txHex, &transaction.extrinsicChainX)
		if err != nil {
			return nil, err
		}
	} else {
		transaction.extrinsic = &types.Extrinsic{}
		err := types.DecodeFromHexString(txHex, &transaction.extrinsic)
		if err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

func (t *Tx) newTx(isChainX bool, call string, args ...interface{}) (*Transaction, error) {
	if t.metadata == nil {
		return nil, ErrNilMetadata
	}
	var (
		transaction = &Transaction{}
	)

	callType, err := types.NewCall(t.metadata, call, args...)
	if err != nil {
		return nil, err
	}

	if isChainX {
		extrinsic := chainxTypes.NewExtrinsic(callType)
		transaction.extrinsicChainX = &extrinsic
	} else {
		extrinsic := types.NewExtrinsic(callType)
		transaction.extrinsic = &extrinsic
	}
	return transaction, nil
}

func (t *Tx) NewBalanceTransferTx(dest string, amount int64) (*Transaction, error) {
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(false, "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (t *Tx) NewChainXBalanceTransferTx(dest string, amount int64) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(true, "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (t *Tx) NewComingNftTransferTx(dest string, cid int64) (*Transaction, error) {
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(false, "ComingNFT.transfer", types.NewU64(uint64(cid)), destAccountID)
}

func (t *Tx) NewXAssetsTransferTx(dest string, amount int64) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(true, "XAssets.transfer", destAccountID, types.NewUCompactFromUInt(uint64(1)), types.NewUCompactFromUInt(uint64(amount)))
}

func (t *Tx) NewThreshold(thresholdPublicKey, destAddress, aggSignature, aggPublicKey, controlBlock, message, scriptHash string, transferAmount int64, blockNumber int32) (*Transaction, error) {
	thresholdPublicKeyByte, err := types.HexDecodeString(thresholdPublicKey)
	if err != nil {
		return nil, err
	}

	destPublicKey, err := AddressToPublicKey(destAddress)
	if err != nil {
		return nil, err
	}
	destPublicKeyByte, err := types.HexDecodeString(destPublicKey)
	if err != nil {
		return nil, err
	}

	aggSignatureByte, err := types.HexDecodeString(aggSignature)
	if err != nil {
		return nil, err
	}

	aggPublicKeyByte, err := types.HexDecodeString(aggPublicKey)
	if err != nil {
		return nil, err
	}

	controlBlockByte, err := types.HexDecodeString(controlBlock)
	if err != nil {
		return nil, err
	}

	messageByte, err := types.HexDecodeString(message)
	if err != nil {
		return nil, err
	}

	scriptHashByte, err := types.HexDecodeString(scriptHash)
	if err != nil {
		return nil, err
	}

	passScriptCall, err := types.NewCall(t.metadata, "ThresholdSignature.pass_script", types.NewAccountID(thresholdPublicKeyByte), types.NewBytes(aggSignatureByte), types.NewBytes(aggPublicKeyByte), types.NewBytes(controlBlockByte), types.NewBytes(messageByte), types.NewBytes(scriptHashByte))
	if err != nil {
		return nil, err
	}

	execScriptCall, err := types.NewCall(t.metadata, "ThresholdSignature.exec_script", types.NewAccountID(destPublicKeyByte), types.NewU8(0), types.NewU128(*big.NewInt(transferAmount)), types.NewU32(uint32(blockNumber)), types.NewU32(uint32(blockNumber))+types.NewU32(1000))
	if err != nil {
		return nil, err
	}

	arg := []types.Call{passScriptCall, execScriptCall}

	return t.newTx(false, "Utility.batch_all", arg)
}
