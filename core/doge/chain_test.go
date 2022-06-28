package doge

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/ignored"
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
	address = ignored.Accounts.DogeMainnet.Address
	chain := chains.mainnet.Chain()

	jsonString, err := chain.FetchUtxos(address, 20)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(jsonString.Value)
}

func TestChain_SendRawTransaction(t *testing.T) {
	txHex := "01000000014a892af99a2b4b1a2a4bf157832c3870b26576e844ae1be2dfa962c02d8d7255000000008b4830450221009e5e3554f6ea7a7f90e00ef307ee7bef37ed7435d6f610f842cb7f2db64dfdc1022053e594c4d4a185d8b3f03d081be667bca252d985ce4a3b1c84ea978677148aa8014104a09e8182977710bab64472c0ecaf9e52255a890554a00a62facd05c0b13817f8995bf590851c19914bfc939d53365b90cc2f0fcfddaca184f0c1e7ce1736f0b80000000003a0860100000000001976a9144afe03f863d27be1cfb7ec0859c4ff89569bb23988ac0000000000000000326a3035516a706f3772516e7751657479736167477a6334526a376f737758534c6d4d7141754332416255364c464646476a3800350c00000000001976a9144da9bb5dea4c42219a2a120523d1a0ce6c268f3788ac00000000"
	txHex = "0100000001d0f37086e51a233afe7ad4e36c6341c1a7762a42d3bd72ea40bbf79d247512f0010000008a47304402205a0aa6a9f92544617f6ecfa74b34dcf6cda4ccde8562e03e52329f61d33676d902206f1f36b034d3f8443b41a37c4841c19a017d8ba305718fa2baff94fa5d57ed2e014104cfb7f626025d6826253f8fc5858e7a5c4b853350b4385ae909cf66138b71ec772bac329b29f6eca2e1d25170a346c041e3dca5ebd68a91ac36cc53cd3be4908c000000000200e1f505000000001976a914a605e6f1e47a3dcad69e555732b3aa416dab4a8b88acc06982900b0000001976a9148a5c58f6c47d8cfd467be2a96df606b7ed7174b288ac00000000"
	chain := chains.mainnet.Chain()

	hash, err := chain.SendRawTransaction(txHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}
