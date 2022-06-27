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

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name     string
		chain    chainInfo
		hash     string
		wantTime int64
		wantErr  bool
	}{
		{
			name:     "doge main",
			chain:    chains.mainnet,
			hash:     "7bc313903372776e1eb81d321e3fe27c9721ce8e71a9bcfee1bde6baea31b5c2",
			wantTime: 1656058561,
		},
		{
			name:     "doge main with 0x",
			chain:    chains.mainnet,
			hash:     "0x7bc313903372776e1eb81d321e3fe27c9721ce8e71a9bcfee1bde6baea31b5c2",
			wantTime: 1656058561,
		},
		{
			name:    "doge main error hash",
			chain:   chains.mainnet,
			hash:    "7bc313903372776e1eb81d321e3fe27c9721ce8e71a9bcfee1bde",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.chain.Chain()
			got, err := chain.FetchTransactionDetail(tt.hash)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.FinishTimestamp != tt.wantTime {
				t.Errorf("FetchTransactionDetail() got = %v, wantTime %v", got, tt.wantTime)
			} else {
				t.Logf("FetchTransactionDetail() got = %v, maybe you can check at %v/tx/%v", got, tt.chain.scan, tt.hash)
			}
		})
	}
}

func TestChain_FetchUtxos(t *testing.T) {
	address := "D8aDCsK4TA9NYhmwiqw1BjZ4CP8LQ814Ea"
	chain := chains.mainnet.Chain()

	jsonString, err := chain.FetchUtxos(address, 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jsonString.Value)
}
