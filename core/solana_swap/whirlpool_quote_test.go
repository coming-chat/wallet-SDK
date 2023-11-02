package solanaswap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWhirlpool_quote(t *testing.T) {
	cli := devChain.Client()
	quote, err := GetSwapQuote(cli, SwapQuoteParam{
		poolAddress:       "b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV",
		tokenMint:         "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
		tokenAmount:       1e6 * 1,
		isInput:           true,
		slippageTolerance: *NewPercentage(1, 100),
		refresh:           true,
	})
	require.Nil(t, err)
	t.Log(quote)
	// 10 USDC -> 4102.988965520 TEST
	// 1 USDC -> 449.556945783 TEST
}

// BtoA
// "tick_array_0": "CXmxVvENVutfAmmHUSVNatgcidiu26uSXuCK8ufvqfxp",
// "tick_array_1": "A3hkPb9EgHCTY6QiduwCLojmY9HzMBZW5LXANqSUYmgk",
// "tick_array_2": "A3hkPb9EgHCTY6QiduwCLojmY9HzMBZW5LXANqSUYmgk",

// AtoB
// "tick_array_0": "CXmxVvENVutfAmmHUSVNatgcidiu26uSXuCK8ufvqfxp",
// "tick_array_1": "CXmxVvENVutfAmmHUSVNatgcidiu26uSXuCK8ufvqfxp",
// "tick_array_2": "CXmxVvENVutfAmmHUSVNatgcidiu26uSXuCK8ufvqfxp",

// CXmxVvENVutfAmmHUSVNatgcidiu26uSXuCK8ufvqfxp
// A3hkPb9EgHCTY6QiduwCLojmY9HzMBZW5LXANqSUYmgk
// BdV2dMquVMS7G1fPGMr6pg4CFbPgE1mWXbXzN5s9a5Wv

// ❯ ts-node scripts/test_dex/get_whirlpool_quote_config.ts b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV 4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU 1
// {
//   "whirlpool_program": "whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc",
//   "whirlpool": "b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV",
//   "token_mint_a": "281LhxeKQ2jaFDx9HAHcdrU9CpedSH7hx5PuRrM7e1FS",
//   "token_mint_b": "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU",
//   "token_owner_account_a": "586Vh32drG45YZVRw83pPp6F5ed7tprXg7gL8Kg2hYM9",
//   "token_owner_account_b": "6NKguqRWqC26R4zJxt9qWqoxxxVkM2Rtn6UmpPtZJbYZ",
//   "token_vault_a": "3dycP3pym3q6DgUpZRviaavaScwrrCuC6QyLhiLfSXge",
//   "token_vault_b": "969UqMJSqvgxmNuAWZx91PAnLJU825qJRAAcEVQMWASg",
//   "tick_array_0": "CXmxVvENVutfAmmHUSVNatgcidiu26uSXuCK8ufvqfxp",
//   "tick_array_1": "A3hkPb9EgHCTY6QiduwCLojmY9HzMBZW5LXANqSUYmgk",
//   "tick_array_2": "A3hkPb9EgHCTY6QiduwCLojmY9HzMBZW5LXANqSUYmgk",
//   "oracle": "44xQG1Fgv5k3Us1s5Mcg6MQiQV2oSeocBRwo7hZvKdRo",
//   "is_a_to_b": false,
//   "amount_in": "1000000",
//   "estimated_amount_out": "449556945783",
//   "min_amount_out": "445105886913"
// }
