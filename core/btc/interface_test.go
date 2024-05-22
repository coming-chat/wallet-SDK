package btc

import (
	"github.com/coming-chat/wallet-SDK/core/base"
)

var (
	_ base.Account = (*Account)(nil)
	_ base.Chain   = (*Chain)(nil)
	_ base.Token   = (*Chain)(nil)
	// _ base.Transaction = (*Transaction)(nil)

	_ base.Transaction       = (*Brc20MintTransaction)(nil)
	_ base.SignedTransaction = (*Brc20MintTransaction)(nil)
	_ base.Transaction       = (*PsbtTransaction)(nil)
	_ base.SignedTransaction = (*SignedPsbtTransaction)(nil)
	_ base.Transaction       = (*Transaction)(nil)
	_ base.SignedTransaction = (*SignedTransaction)(nil)
)
