package doge

import (
	"github.com/coming-chat/wallet-SDK/core/base"
)

var (
	_ base.Account = (*Account)(nil)
	_ base.Chain   = (*Chain)(nil)
	_ base.Token   = (*Chain)(nil)
	// _ base.Transaction = (*Transaction)(nil)
)
