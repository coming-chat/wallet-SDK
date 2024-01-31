package eth

import (
	"github.com/coming-chat/wallet-SDK/core/base"
)

var (
	_ base.Account     = (*Account)(nil)
	_ base.Chain       = (*Chain)(nil)
	_ base.Token       = (*Token)(nil)
	_ base.Token       = (*Erc20Token)(nil)
	_ base.Transaction = (*Transaction)(nil)

	_ base.Token = (*Src20Token)(nil)
)
