package wallet

import "errors"

type PolkaChain struct {
	RpcUrl  string
	ScanUrl string
}

// Deprecated: WalletPolkaChain is deprecated. Please Use PolkaChain instead.
func NewPolkaChain(rpcUrl, scanUrl string) *PolkaChain {
	return nil
}

// Deprecated: WalletPolkaChain is deprecated. Please Use PolkaChain() and call LoadCachedMetadataString() instead.
func NewPolkaChainWithRpc(rpcUrl, scanUrl string, metadataString string) (*PolkaChain, error) {
	return nil, errors.New("WalletPolkaChain is deprecated. Please Use PolkaChain instead.")
}
