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
	signatureOptions *types.SignatureOptions
	extrinsic        *types.Extrinsic
	extrinsicChainX  *chainxTypes.Extrinsic
	extrinsicEra     *types.ExtrinsicEra
	payload          *types.ExtrinsicPayloadV4
	signature        *types.Signature
	extSig           *types.ExtrinsicSignatureV4
	SignatureData    []byte
	PublicKey        []byte
}

func (t *Transaction) GetSignData() ([]byte, error) {
	return types.EncodeToBytes(t.payload)
}

func (t *Transaction) GetTx() (tx string, err error) {
	if t.SignatureData == nil {
		return "", errNotSigned
	}

	if t.PublicKey == nil {
		return "", errNoPublicKey
	}
	signatureData := types.NewSignature(t.SignatureData)

	if t.extrinsicChainX != nil {
		extSig := chainxTypes.ExtrinsicSignatureV4{
			Signer:    types.NewAddressFromAccountID(t.PublicKey),
			Signature: types.MultiSignature{IsSr25519: true, AsSr25519: signatureData},
			Era:       *t.extrinsicEra,
			Nonce:     t.signatureOptions.Nonce,
			Tip:       t.signatureOptions.Tip,
		}
		t.extrinsicChainX.Signature = extSig

		// mark the extrinsic as signed
		t.extrinsicChainX.Version |= types.ExtrinsicBitSigned
		tx, err = types.EncodeToHexString(t.extrinsicChainX)
	} else {
		extSig := types.ExtrinsicSignatureV4{
			Signer:    types.NewMultiAddressFromAccountID(t.PublicKey),
			Signature: types.MultiSignature{IsSr25519: true, AsSr25519: signatureData},
			Era:       *t.extrinsicEra,
			Nonce:     t.signatureOptions.Nonce,
			Tip:       t.signatureOptions.Tip,
		}
		t.extrinsic.Signature = extSig

		// mark the extrinsic as signed
		t.extrinsic.Version |= types.ExtrinsicBitSigned
		tx, err = types.EncodeToHexString(t.extrinsic)
	}

	if err != nil {
		return "", err
	}
	return tx, nil
}

func (t *Tx) newTx(isChainX bool, genesisHashString string, nonce uint64, specVersion, transVersion uint32, call string, args ...interface{}) (*Transaction, error) {
	if t.metadata == nil {
		return nil, errNilMetadata
	}
	var (
		transaction = &Transaction{}
		methodBytes []byte
	)

	genesisHash, err := types.NewHashFromHexString(genesisHashString)
	if err != nil {
		return nil, err
	}

	callType, err := types.NewCall(t.metadata, call, args...)
	if err != nil {
		return nil, err
	}

	if isChainX {
		extrinsic := chainxTypes.NewExtrinsic(callType)
		transaction.extrinsicChainX = &extrinsic
		methodBytes, err = types.EncodeToBytes(transaction.extrinsicChainX.Method)
		if err != nil {
			return nil, err
		}
	} else {
		extrinsic := types.NewExtrinsic(callType)
		transaction.extrinsic = &extrinsic
		methodBytes, err = types.EncodeToBytes(transaction.extrinsic.Method)
		if err != nil {
			return nil, err
		}
	}

	transaction.signatureOptions = &types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsImmortalEra: true},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(nonce),
		SpecVersion:        types.NewU32(specVersion),
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: types.NewU32(transVersion),
	}

	transaction.extrinsicEra = &types.ExtrinsicEra{IsImmortalEra: true}

	transaction.payload = &types.ExtrinsicPayloadV4{
		ExtrinsicPayloadV3: types.ExtrinsicPayloadV3{
			Method:      methodBytes,
			Era:         *transaction.extrinsicEra,
			Nonce:       transaction.signatureOptions.Nonce,
			Tip:         transaction.signatureOptions.Tip,
			SpecVersion: transaction.signatureOptions.SpecVersion,
			GenesisHash: transaction.signatureOptions.GenesisHash,
			BlockHash:   transaction.signatureOptions.BlockHash,
		},
		TransactionVersion: transaction.signatureOptions.TransactionVersion,
	}

	return transaction, nil
}

func (t *Tx) NewBalanceTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(false, genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (t *Tx) NewChainXBalanceTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(true, genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (t *Tx) NewComingNftTransferTx(dest, genesisHashString string, cid, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(false, genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "ComingNFT.transfer", types.NewU64(uint64(cid)), destAccountID)
}

func (t *Tx) NewXAssetsTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newTx(true, genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "XAssets.transfer", destAccountID, types.NewUCompactFromUInt(uint64(1)), types.NewUCompactFromUInt(uint64(amount)))
}

func (t *Tx) NewThreshold(thresholdPublicKey, destAddress, aggSignature, aggPublicKey, controlBlock, message, scriptHash, genesisHashString string, transferAmount, nonce int64, blockNumber, specVersion, transVersion int32) (*Transaction, error) {
	thresholdPublicKeyByte, err := types.HexDecodeString(thresholdPublicKey)
	if err != nil {
		return nil, err
	}

	destPublicKeyByte, err := types.HexDecodeString(AddressToPublicKey(destAddress))
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

	return t.newTx(false, genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "Utility.batch_all", arg)
}
