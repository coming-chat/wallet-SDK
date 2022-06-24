package doge

import (
	"testing"
)

type chainInfo struct {
	net  string
	scan string
}

type chainCfg struct {
	mainnet chainInfo
	testnet chainInfo
}

var chains = &chainCfg{
	mainnet: chainInfo{
		net:  ChainMainnet,
		scan: "https://dogechain.info/",
	},
	testnet: chainInfo{
		net:  ChainTestnet,
		scan: "https://sochain.com/testnet/doge",
	},
}

func (c *chainInfo) Chain() *Chain {
	chain, _ := NewChainWithChainnet(c.net)
	return chain
}

func TestChain_BalanceOfAddress(t *testing.T) {
	tests := []struct {
		name    string
		chain   chainInfo
		address string
		wantErr bool
	}{
		{
			name:    "doge miannet",
			chain:   chains.mainnet,
			address: "DBx1XSBxpSUnEK79nA8VtrKh2qr2LupZ6G",
		},
		{
			name:    "doge testnet (not support now)",
			chain:   chains.testnet,
			address: "nW8tMJ4BxDcc1tKZTBh7uNS8639Aj2Hz6s",
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.chain.Chain()
			got, err := c.BalanceOfAddress(tt.address)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			t.Logf("BalanceOfAddress() got = %v, maybe you can checked at %v/address/%v", got, tt.chain.scan, tt.address)
		})
	}
}
