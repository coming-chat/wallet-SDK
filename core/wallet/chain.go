package wallet

import (
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/itering/subscan/util/ss58"
)

type PolkaBalance struct {
	Total  string
	Usable string
}

func emptyBalance() *PolkaBalance {
	return &PolkaBalance{
		Total:  "0",
		Usable: "0",
	}
}

type PolkaChain struct {
	RpcUrl string
}

func NewPolkaChain(rpc string) *PolkaChain {
	return &PolkaChain{RpcUrl: rpc}
}

// 刷新最新的 metadata
func (c *PolkaChain) ReloadMetadata() error {
	client, err := getPolkaClient(c.RpcUrl)
	if err != nil {
		return err
	}
	return client.ReloadMetadata()
}

// 通过 address 查询余额
func (c *PolkaChain) QueryBalance(address string) (*PolkaBalance, error) {
	ss58Format := base58.Decode(address)
	pubkey, err := hex.DecodeString(ss58.Decode(address, int(ss58Format[0])))
	if err != nil {
		return emptyBalance(), err
	}
	return c.queryBalance(pubkey)
}

// 通过 public key 查询余额
func (c *PolkaChain) QueryBalancePubkey(pubkey string) (*PolkaBalance, error) {
	pubkey = strings.TrimPrefix(pubkey, "0x")
	data, err := hex.DecodeString(pubkey)
	if err != nil {
		return emptyBalance(), ErrPublicKey
	}
	return c.queryBalance(data)
}

func (c *PolkaChain) queryBalance(pubkey []byte) (b *PolkaBalance, err error) {
	b = emptyBalance()

	client, err := getPolkaClient(c.RpcUrl)
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

	freeInt := data.Data.Free.Int
	total := freeInt.Add(freeInt, data.Data.Reserved.Int)

	locked := data.Data.MiscFrozen.Int
	if data.Data.MiscFrozen.Cmp(data.Data.FeeFrozen.Int) <= 0 {
		locked = data.Data.FeeFrozen.Int
	}
	usable := freeInt.Sub(freeInt, locked)

	return &PolkaBalance{
		Total:  total.String(),
		Usable: usable.String(),
	}, nil
}

// 发起交易
func (c *PolkaChain) SendRawTransaction(txHex string) (s string, err error) {
	client, err := getPolkaClient(c.RpcUrl)
	if err != nil {
		return
	}

	var hashString string
	err = client.api.Client.Call(&hashString, "author_submitExtrinsic", txHex)
	if err != nil {
		return
	}

	return hashString, nil
}
