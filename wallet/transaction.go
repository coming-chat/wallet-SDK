package wallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type Transaction struct {
	signatureOptions *types.SignatureOptions
	extrinsic        *types.Extrinsic
	extrinsicEra     *types.ExtrinsicEra
	payload          *types.ExtrinsicPayloadV4
	signature        *types.Signature
	extSig           *types.ExtrinsicSignatureV4
}

func (w *Wallet) newTx(genesisHashString string, nonce uint64, specVersion, transVersion uint32, call string, args ...interface{}) (*Transaction, error) {
	if w.metadata == nil {
		return nil, ErrNilMetadata
	}
	var transaction = &Transaction{}

	genesisHash, err := types.NewHashFromHexString(genesisHashString)
	if err != nil {
		return nil, err
	}

	callType, err := types.NewCall(w.metadata, call, args)
	if err != nil {
		return nil, err
	}

	extrinsic := types.NewExtrinsic(callType)

	methodBytes, err := types.EncodeToBytes(extrinsic.Method)
	if err != nil {
		return nil, err
	}

	transaction.extrinsic = &extrinsic
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

func (w *Wallet) NewBalanceTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return w.newTx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "Balances.transfer", destAccountID, types.NewUCompactFromUInt(uint64(amount)))
}

func (w *Wallet) NewComingNftTransferTx(dest, genesisHashString string, cid, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return w.newTx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "ComingNFT.transfer", types.NewU64(uint64(cid)), destAccountID)
}

func (w *Wallet) NewXAssetsTransferTx(dest, genesisHashString string, amount, nonce int64, specVersion, transVersion int32) (*Transaction, error) {
	//LookupSource
	destAccountID, err := addressStringToAddress(dest)
	if err != nil {
		return nil, err
	}
	return w.newTx(genesisHashString, uint64(nonce), uint32(specVersion), uint32(transVersion), "XAssets.transfer", destAccountID, types.NewUCompactFromUInt(uint64(1)), types.NewUCompactFromUInt(uint64(amount)))
}
