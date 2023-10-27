package solanaswap

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/coming-chat/wallet-SDK/core/solana"
	"github.com/coming-chat/wallet-SDK/core/testcase"
)

type ConfigProgram struct {
	SoDiamond   string
	Wormhole    string
	TokenBridge string
	Whirlpools  string
}

type ConfigWormhole struct {
	Chainid          int
	Actual_reserve   int
	Estimate_reserve int
	Dst_chain        ConfigDSTChain
}

type ConfigDSTChain struct {
	OmnibtcChainid       uint16
	Chainid              uint16
	Base_gas             int
	Per_byte_gas         int
	Price_manager        string
	Token_bridge_emitter string
	Omniswap_emitter     string
}

type ConfigTokenss struct {
	TEST       ConfigToken
	USDC       ConfigToken
	WrappedBSC ConfigToken
}

type ConfigToken struct {
	Mint     string
	Decimals uint8
}

type ConfigPools struct {
	Whirlpool_TEST_USDC string
}

type ConfigLookupTable struct {
	Key       string
	Addresses []string
}

type Config struct {
	OmnibtcChainid uint16
	Program        ConfigProgram
	Wormhole       ConfigWormhole
	Token          ConfigTokenss
	Pools          ConfigPools
	Lookup_Table   ConfigLookupTable

	Beneficiary    string
	Redeemer_proxy string
	So_fee_by_ray  int
}

var DevnetConfig = Config{
	OmnibtcChainid: 30006,
	Program: ConfigProgram{
		SoDiamond:   "5DncnqicaHDZTMfkcfzKaYP5XzD5D9jg3PGNTT5J1Qg7",
		Wormhole:    "3u8hJUVTA4jH1wYAyUur7FFZVQ8H635K3tSHHF4ssjQ5",
		TokenBridge: "DZnkkTmCiFWfYTfT41X3Rd1kDgozqzxWaHqsw6W4x2oe",
		Whirlpools:  "whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc",
	},
	Wormhole: ConfigWormhole{
		Chainid:          1,
		Actual_reserve:   110000000,
		Estimate_reserve: 120000000,
		Dst_chain: ConfigDSTChain{
			OmnibtcChainid:       30003,
			Chainid:              4,
			Base_gas:             700000,
			Per_byte_gas:         68,
			Price_manager:        "4q2wPZMys1zCoAVpNmhgmofb6YM9MqLXmV25LdtEMAf9",
			Token_bridge_emitter: "0x0000000000000000000000009dcf9d205c9de35334d646bee44b2d2859712a09",
			Omniswap_emitter:     "0x00000000000000000000000084b7ca95ac91f8903acb08b27f5b41a4de2dc0fc",
		},
	},
	Token: ConfigTokenss{
		TEST: ConfigToken{
			Mint:     "281LhxeKQ2jaFDx9HAHcdrU9CpedSH7hx5PuRrM7e1FS",
			Decimals: 9,
		},
		USDC: ConfigToken{
			Mint:     "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
			Decimals: 6,
		},
		WrappedBSC: ConfigToken{
			Mint:     "xxtdhpCgop5gZSeCkRRHqiVu7hqEC9MKkd1xMRUZqrz",
			Decimals: 8,
		},
	},
	Pools: ConfigPools{
		Whirlpool_TEST_USDC: "b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV",
	},
	Lookup_Table: ConfigLookupTable{
		Key: "ESxWFjHVo2oes1eAQiwkAUHNTTUT9Xm5zsSrE7QStYX8",
		Addresses: []string{
			// token_bridge_config
			"8PFZNjn19BBYVHNp4H31bEW7eAmu78Yf2RKV8EeA461K",
			// token_bridge_authority_signer
			"3VFdJkFuzrcwCwdxhKRETGxrDtUVAipNmYcLvRBDcQeH",
			// token_bridge_custody_signer
			"H9pUTqZoRyFdaedRezhykA1aTMq7vbqRHYVhpHZK2QbC",
			// token_bridge_mint_authority
			"rRsXLHe7sBHdyKU3KY3wbcgWvoT1Ntqudf6e9PKusgb",
			// wormhole_bridge
			"6bi4JGDoRwUs9TYBuvoA7dUVyikTJDrJsJU1ew6KVLiu",
			// token_bridge_emitter
			"4yttKWzRoNYS2HekxDfcZYmfQqnVWpKiJ8eydYRuFRgs",
			// wormhole_fee_collector
			"7s3a1ycs16d6SNDumaRtjcoyMaTDZPavzgsmS3uUZYWX",
			// token_bridge_sequence
			"9QzqZZvhxoHzXbNY9y2hyAUfJUzDwyDb7fbDs9RXwH3",
			// Rent
			"SysvarRent111111111111111111111111111111111",
			// Clock
			"SysvarC1ock11111111111111111111111111111111",
			// ComputeBudget
			"ComputeBudget111111111111111111111111111111",
			// System
			"11111111111111111111111111111111",
			// whirlpool_program
			"whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc",
			// test_usdc_pool
			"b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV",
			// test token
			"281LhxeKQ2jaFDx9HAHcdrU9CpedSH7hx5PuRrM7e1FS",
			// usdc token
			"4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
			// test_usdc_pool vault_a
			"3dycP3pym3q6DgUpZRviaavaScwrrCuC6QyLhiLfSXge",
			// test_usdc_pool vault_b
			"969UqMJSqvgxmNuAWZx91PAnLJU825qJRAAcEVQMWASg",
			// omniswap sender config
			"GR7xDWrbWcEYsnz1e5WDfy3iXvPw5tmWjeV8MY1sVHCp",
			// omniswap fee config
			"EcZK7hAyxzjeCL1zM9FKWeWcdziF4pFHiUCJ5r2886TP",
			// omniswap bsc(4) price manage
			"EofptCXfgVxRk1vxBLNP1Zk6SSPBiPdkYWVPgTLzbzGR",
			// omniswap bsc(4) foreign_contract
			"FV2SB6pUGWABHxmnoVUxxdTVctzY7puAQon38sJ8oNm",
		},
	},

	Beneficiary:    "vQkE51MXJiwqtbwf562XWChNKZTgh6L2jHPpupoCKjS",
	Redeemer_proxy: "4q2wPZMys1zCoAVpNmhgmofb6YM9MqLXmV25LdtEMAf9",
	// ray = 100_000_000
	So_fee_by_ray: 0,
}

var config = DevnetConfig
var account, _ = solana.NewAccountWithMnemonic(testcase.M3)
var devChain = solana.NewChainWithRpc(solana.DevnetRPCEndpoint)

var PROGRAM_ID = common.PublicKeyFromString("5DncnqicaHDZTMfkcfzKaYP5XzD5D9jg3PGNTT5J1Qg7")
