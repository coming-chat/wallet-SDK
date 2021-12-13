package wallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"wallet-SDK/chainxTypes"
)

type TxMetadata struct {
	metadata *types.Metadata
}

func NewTxMetadata(metadataString string) (*TxMetadata, error) {
	var metadata types.Metadata
	if err := types.DecodeFromHexString(metadataString, &metadata); err != nil {
		return nil, err
	}
	return &TxMetadata{
		metadata: &metadata,
	}, nil
}

type Transaction struct {
	signatureOptions *types.SignatureOptions
	extrinsicV4      *types.Extrinsic
	extrinsicV3      *chainxTypes.Extrinsic
	extrinsicEra     *types.ExtrinsicEra
	payload          *types.ExtrinsicPayloadV4
	signature        *types.Signature
	extSig           *types.ExtrinsicSignatureV4
}

func (t *TxMetadata) newV4Tx(genesisHashString string, nonce uint64, specVersion, transVersion uint32, call string, args ...interface{}) (*Transaction, error) {
	if t.metadata == nil {
		return nil, ErrNilMetadata
	}
	var transaction = &Transaction{}

	genesisHash, err := types.NewHashFromHexString(genesisHashString)
	if err != nil {
		return nil, err
	}

	callType, err := types.NewCall(t.metadata, call, args...)
	if err != nil {
		return nil, err
	}

	extrinsic := types.NewExtrinsic(callType)

	methodBytes, err := types.EncodeToBytes(extrinsic.Method)
	if err != nil {
		return nil, err
	}

	transaction.extrinsicV4 = &extrinsic
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

func (t *TxMetadata) newV3Tx(genesisHashString string, nonce uint64, specVersion, transVersion uint32, call string, args ...interface{}) (*Transaction, error) {
	if t.metadata == nil {
		return nil, ErrNilMetadata
	}
	var transaction = &Transaction{}

	genesisHash, err := types.NewHashFromHexString(genesisHashString)
	if err != nil {
		return nil, err
	}

	callType, err := chainxTypes.NewCall(t.metadata, call, args...)
	if err != nil {
		return nil, err
	}

	extrinsic := chainxTypes.NewExtrinsic(callType)

	methodBytes, err := types.EncodeToBytes(extrinsic.Method)
	if err != nil {
		return nil, err
	}

	transaction.extrinsicV3 = &extrinsic
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

func (t *TxMetadata) NewV4BalanceTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newV4Tx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (t *TxMetadata) NewV3BalanceTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newV3Tx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (t *TxMetadata) NewComingNftTransferTx(dest, genesisHashString string, cid, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newV4Tx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "ComingNFT.transfer", types.NewU64(uint64(cid)), destAccountID)
}

func (t *TxMetadata) NewXAssetsTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	//LookupSource
	destAccountID, err := addressStringToMultiAddress(dest)
	if err != nil {
		return nil, err
	}
	return t.newV4Tx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "XAssets.transfer", destAccountID, types.NewUCompactFromUInt(uint64(1)), types.NewUCompactFromUInt(uint64(amount)))
}
