package wallet

import (
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

// Deprecated: WalletTx is deprecated. Please Use PolkaTx instead.
type Tx struct {
	metadata *types.Metadata
}

// Deprecated: WalletTx is deprecated. Please Use PolkaTx instead.
func NewTx(metadataString string) (*Tx, error) {
	return nil, errors.New("WalletTx is deprecated. Please Use PolkaTx instead.")
}

// Deprecated: WalletTransaction is deprecated. Please Use PolkaTransaction instead.
type Transaction struct {
	extrinsic *types.Extrinsic
}
