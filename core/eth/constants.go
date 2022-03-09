package eth

var (
	// ERC20  交易method
	ERC20_METHOD_TRANSFER = "transfer"
	ERC20_METHOD_APPROVE  = "approve"
)

// 默认gas limit估算失败后，21000 * 3 = 63000
var (
	DEFAULT_CONTRACT_GAS_LIMIT = "63000"
	DEFAULT_ETH_GAS_LIMIT      = "21000"
	// 当前网络 standard gas price
	DEFAULT_ETH_GAS_PRICE = "20000000000"
)
