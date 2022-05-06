package polka

import (
	"testing"
)

type rpcInfo struct {
	url  string
	scan string
	net  int

	realScan string
}
type rpcConfigs struct {
	chainxProd  rpcInfo
	chainxTest  rpcInfo
	minixProd   rpcInfo
	minixTest   rpcInfo
	sherpaxProd rpcInfo
	sherpaxTest rpcInfo
	polkadot    rpcInfo
	kusama      rpcInfo
}

var rpcs = rpcConfigs{
	chainxProd:  rpcInfo{"https://mainnet.chainx.org/rpc", "https://multiscan-api.coming.chat/chainx/extrinsics", 44, "https://scan.chainx.org"},
	chainxTest:  rpcInfo{"https://testnet3.chainx.org/rpc", "https://multiscan-api-pre.coming.chat/chainx/extrinsics", 44, "https://scan-pre.chainx.cc"},
	minixProd:   rpcInfo{"https://minichain-mainnet.coming.chat/rpc", "https://multiscan-api.coming.chat/minix/extrinsics", 44, "https://mini-scan.coming.chat"},
	minixTest:   rpcInfo{"https://rpc-minichain.coming.chat", "https://multiscan-api-pre.coming.chat/minix/extrinsics", 44, "https://mini-scan-pre.coming.chat"},
	sherpaxProd: rpcInfo{"https://mainnet.sherpax.io/rpc", "https://multiscan-api.coming.chat/sherpax/extrinsics", 44, "https://scan.sherpax.io"},
	sherpaxTest: rpcInfo{"https://sherpax-testnet.chainx.org/rpc", "https://multiscan-api-pre.coming.chat/sherpax/extrinsics", 44, "https://scan-pre.sherpax.io/"},
	polkadot:    rpcInfo{"https://polkadot.api.onfinality.io/public-rpc", "", 0, "https://polkadot.subscan.io"},
	kusama:      rpcInfo{"https://kusama.api.onfinality.io/public-rpc", "", 2, "https://kusama.subscan.io"},
}

func (n *rpcInfo) Chain() (*Chain, error) {
	return NewChainWithRpc(n.url, n.scan, n.net)
}

func (n *rpcInfo) MetadataString() (string, error) {
	c, err := NewChainWithRpc(n.url, n.scan, n.net)
	if err != nil {
		return "", err
	}
	return c.GetMetadataString()
}

// MARK - Test Start

func TestChain_BalanceOfAddress(t *testing.T) {
	tests := []struct {
		name    string
		rpcInfo rpcInfo
		address string
		wantErr bool
	}{
		{name: "chainx prod empty account", rpcInfo: rpcs.chainxProd, address: accountCase.address44},
		{name: "chainx pre normal account", rpcInfo: rpcs.chainxTest, address: "5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1"},
		{name: "minix prod normal account", rpcInfo: rpcs.minixProd, address: "5RR4TW6b1cA3iqPUZDmQFPT6c5onmXQuvwLWkCmdAz3q8h38"},
		{name: "sherpax prod normal account", rpcInfo: rpcs.sherpaxProd, address: "5SRWCZMfBawyCviWciDmAQ6bNoeH4w3yDRkz7L73vNbWGLrF"},
		{name: "polkadot prod empty account", rpcInfo: rpcs.polkadot, address: accountCase.address0},
		{name: "polkadot prod normal account", rpcInfo: rpcs.polkadot, address: "13wNbioJt44NKrcQ5ZUrshJqP7TKzQbzZt5nhkeL4joa3PAX"},
		{name: "kusama prod normal account", rpcInfo: rpcs.kusama, address: "E8ky1h79157Fpcm8cXFyYvRHhN5pb67z6MPcbQ9Z7kgiTV9"},
		{name: "chainx prod error address", rpcInfo: rpcs.chainxTest, address: accountCase.address44 + "s",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := tt.rpcInfo.Chain()
			got, err := c.BalanceOfAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceOfAddress() error = %v, wantErrs %v", err, tt.wantErr)
				return
			}
			if err == nil {
				url := tt.rpcInfo.realScan + "/account/" + tt.address
				t.Log("BalanceOfAddress() result: ", got, ", Maybe you should verify via the link: ", url)
			}
		})
	}
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name     string
		rpcInfo  rpcInfo
		hash     string
		wantTime int64
		wantErr  bool
	}{
		{
			name:     "chainx prod normal",
			rpcInfo:  rpcs.chainxProd,
			hash:     "0x1b38149fe30b7b8aa9cd240b4a871c72ae3e4baf333beaa061a68b5d54234385",
			wantTime: 1651819326,
		},
		{
			name:     "chainx prod xbtc transfer",
			rpcInfo:  rpcs.chainxProd,
			hash:     "0xed0bec9a427557bee16f6d85e505107151a8c2e41e87a9ce6e1020ba542cb4fb",
			wantTime: 1630142208,
		},
		{
			name:     "minix prod normal",
			rpcInfo:  rpcs.minixProd,
			hash:     "0xb490e8d7ec951a009d45571020e07552efe43db8267ea3e462417a7e6ac41dd4",
			wantTime: 1651584702,
		},
		{
			name:     "minix prod failed transfer", // Balance too low to send value
			rpcInfo:  rpcs.minixProd,
			hash:     "0xd4dea2519c34e9e38a45331063188560abae770adf4216dff66c43ecf05f4584",
			wantTime: 1651647450,
		},
		{
			name:     "sherpax prod normal",
			rpcInfo:  rpcs.sherpaxProd,
			hash:     "0x53add48fada48d32da326094f50a069a51406fb5a2526676bdf796b4a592a872",
			wantTime: 1644213546,
		},
		{
			name:    "polkadot prod normal, but we not support yet",
			rpcInfo: rpcs.polkadot,
			hash:    "0x3d8c531fe52dfda116c2988fcc7728cbc19543a7ea713812c0a7f5501f815cbe",
			wantErr: true, // not support polka & kusama yet.
		},
		{
			name:    "chainx prod error hash",
			rpcInfo: rpcs.chainxProd,
			hash:    "0x3d8c531fe52dfda116c2988fcc7728cbc19543a7ea713812c0a7f5501f815cbe",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := tt.rpcInfo.Chain()
			got, err := c.FetchTransactionDetail(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && tt.wantTime != got.FinishTimestamp {
				t.Errorf("FetchTransactionDetail() got = %v, want %v", got, tt.wantTime)
			}
		})
	}
}
