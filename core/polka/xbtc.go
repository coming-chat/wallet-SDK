package polka

import (
	"encoding/hex"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	CustomType "github.com/coming-chat/wallet-SDK/core/substrate/types"
	"github.com/decred/base58"
	"github.com/itering/subscan/util/ss58"
)

type XBTCToken struct {
	chain *Chain
}

// Warning: initial unavailable, You must create based on Chain.XBTCToken()
func NewXBTCToken() (*Token, error) {
	return nil, errors.New("Token initial unavailable, You must create based on Chain.XBTCToken()")
}

// MARK - Implement the protocol Token, Override

func (t *XBTCToken) TokenInfo() (*base.TokenInfo, error) {
	return &base.TokenInfo{
		Name:    "XBTC",
		Symbol:  "XBTC",
		Decimal: 8,
	}, nil
}

func (t *XBTCToken) BalanceOfAddress(address string) (*base.Balance, error) {
	return t.queryBalance(address)
}

func (t *XBTCToken) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	address, err := t.chain.EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return t.queryBalance(address)
}

func (t *XBTCToken) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.queryBalance(account.Address())
}

func (t *XBTCToken) queryBalance(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	b = base.EmptyBalance()

	client, err := getConnectedPolkaClient(t.chain.RpcUrl)
	if err != nil {
		return
	}

	err = client.LoadMetadataIfNotExists()
	if err != nil {
		return
	}

	ss58Format := base58.Decode(address)
	publicKey, err := hex.DecodeString(ss58.Decode(address, int(ss58Format[0])))
	if err != nil {
		return
	}

	assetId, err := types.EncodeToBytes(uint32(1))
	if err != nil {
		return
	}

	metadata := client.metadata
	call, err := types.CreateStorageKey(metadata, "XAssets", "AssetBalance", publicKey, assetId)
	if err != nil {
		return
	}
	entryMetadata, err := metadata.FindStorageEntryMetadata("XAssets", "AssetBalance")
	if err != nil {
		return
	}
	i := entryMetadata.(types.StorageEntryMetadataV14).Type.AsMap.Value
	kIndex := metadata.AsMetadataV14.EfficientLookup[i.Int64()].Params[0].Type.Int64()
	vValue := metadata.AsMetadataV14.EfficientLookup[i.Int64()].Params[1].Type.Int64()
	data := CustomType.NewMap(metadata.AsMetadataV14.EfficientLookup[kIndex], metadata.AsMetadataV14.EfficientLookup[vValue])
	_, err = client.api.RPC.State.GetStorageLatest(call, &data)
	if err != nil {
		return
	}

	usable, ok := data.Data["Usable"]
	if !ok {
		return b, nil
	}
	usableInt := usable.(types.U128).Int

	return &base.Balance{
		Usable: usableInt.String(),
		Total:  usableInt.String(),
	}, nil
}
