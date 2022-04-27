package polka

import (
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	CustomType "github.com/coming-chat/wallet-SDK/core/substrate/types"
	"github.com/decred/base58"
	"github.com/itering/subscan/util/ss58"
)

// MARK - Implement the protocol Chain.Balance

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	ss58Format := base58.Decode(address)
	pubkey, err := hex.DecodeString(ss58.Decode(address, int(ss58Format[0])))
	if err != nil {
		return base.EmptyBalance(), err
	}
	return c.queryBalance(pubkey)
}

func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	publicKey = strings.TrimPrefix(publicKey, "0x")
	data, err := hex.DecodeString(publicKey)
	if err != nil {
		return base.EmptyBalance(), ErrPublicKey
	}
	return c.queryBalance(data)
}

func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfPublicKey(account.PublicKeyHex())
}

// query balance with pubkey data.
func (c *Chain) queryBalance(pubkey []byte) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	b = base.EmptyBalance()

	client, err := getConnectedPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	err = client.LoadMetadataIfNotExists()
	if err != nil {
		return
	}

	call, err := types.CreateStorageKey(client.metadata, "System", "Account", pubkey)
	if err != nil {
		return
	}

	data := struct {
		Nonce       uint32
		Consumers   uint32
		Providers   uint32
		Sufficients uint32
		Data        struct {
			Free       types.U128
			Reserved   types.U128
			MiscFrozen types.U128
			FeeFrozen  types.U128
		}
	}{}

	// Ok is true if the value is not empty.
	ok, err := client.api.RPC.State.GetStorageLatest(call, &data)
	if err != nil {
		return
	}
	if !ok {
		return
	}

	totalInt := big.NewInt(0).Add(data.Data.Free.Int, data.Data.Reserved.Int)
	locked := base.MaxBigInt(data.Data.MiscFrozen.Int, data.Data.FeeFrozen.Int)
	usableInt := big.NewInt(0).Sub(data.Data.Free.Int, locked)

	return &base.Balance{
		Total:  totalInt.String(),
		Usable: usableInt.String(),
	}, nil
}

// 特殊查询 XBTC 的余额
// 只能通过 chainx 链对象来查询，其他链会抛出 error
func (c *Chain) QueryBalanceXBTC(address string) (b *base.Balance, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	b = base.EmptyBalance()

	client, err := getConnectedPolkaClient(c.RpcUrl)
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
