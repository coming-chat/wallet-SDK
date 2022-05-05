package eth

import (
	"errors"

	"github.com/coming-chat/wallet-SDK/core/base"
)

// Deprecated: CoinUtil is deprecated. Please Use EthChain instead.
type CoinUtil struct {
	// 链的 RPC 地址
	RpcUrl string
	// 币种的合约地址，如果为 nil，表示是主网的币
	ContractAddress string
	// 用户的钱包地址
	WalletAddress string
}

// Deprecated: CoinUtil is deprecated. Please Use NewChainWithRpc() instead.
func NewCoinUtilWithRpc(rpcUrl, contractAddress, walletAddress string) *CoinUtil {
	return &CoinUtil{
		RpcUrl:          rpcUrl,
		ContractAddress: contractAddress,
		WalletAddress:   walletAddress,
	}
}

// 是否是主币
func (u *CoinUtil) IsMainCoin() bool {
	return u.ContractAddress == ""
}

// Deprecated: CoinInfo is deprecated.
func (u *CoinUtil) CoinInfo() (*Erc20TokenInfo, error) {
	chain, err := GetConnection(u.RpcUrl)
	if err != nil {
		return nil, err
	}

	if u.IsMainCoin() {
		balance, err := u.QueryBalance()
		if err != nil {
			return nil, err
		}
		return &Erc20TokenInfo{
			Balance: balance,
		}, nil
	} else {
		return chain.Erc20TokenInfo(u.ContractAddress, u.WalletAddress)
	}
}

// Deprecated: QueryBalance is deprecated. Please Use Chain.Token().BalanceOfAddress() instead.
func (u *CoinUtil) QueryBalance() (string, error) {
	balance, err := u.baseToken().BalanceOfAddress(u.WalletAddress)
	return balance.Usable, err
}

// Deprecated: SuggestGasPrice is deprecated. Please Use Chain.SuggestGasPrice() instead.
func (u *CoinUtil) SuggestGasPrice() (string, error) {
	return u.chain().SuggestGasPrice()
}

// Deprecated: Nonce is deprecated. Please Use Chain.NonceOfAddress() instead.
func (u *CoinUtil) Nonce() (string, error) {
	return u.chain().NonceOfAddress(u.WalletAddress)
}

// Deprecated: Nonce is deprecated. Please Use Chain.Token().EstimateGasLimit() instead.
func (u *CoinUtil) EstimateGasLimit(receiverAddress, gasPrice, amount string) (string, error) {
	token, err := u.ethToken()
	if err != nil {
		return "", err
	}
	return token.EstimateGasLimit(u.WalletAddress, receiverAddress, gasPrice, amount)
}

// Deprecated: BuildTransferTx is deprecated. Please Use Chain.Token().BuildTransferTx() instead.
func (u *CoinUtil) BuildTransferTx(privateKey, receiverAddress, nonce, gasPrice, gasLimit, amount string) (string, error) {
	token, err := u.ethToken()
	if err != nil {
		return "", err
	}
	return token.BuildTransferTx(privateKey, receiverAddress, nonce, gasPrice, gasLimit, amount)
}

// Deprecated: SendRawTransaction is deprecated. Please Use Chain.SendRawTransaction() instead.
func (u *CoinUtil) SendRawTransaction(txHex string) (string, error) {
	return u.chain().SendRawTransaction(txHex)
}

// Deprecated: FetchTransactionStatus is deprecated. Please Use Chain.FetchTransactionStatus() instead.
func (u *CoinUtil) FetchTransactionStatus(hashString string) base.TransactionStatus {
	return u.chain().FetchTransactionStatus(hashString)
}

// Deprecated: FetchTransactionDetail is deprecated. Please Use Chain.FetchTransactionDetail() instead.
func (u *CoinUtil) FetchTransactionDetail(hashString string) (*base.TransactionDetail, error) {
	return u.chain().FetchTransactionDetail(hashString)
}

// Deprecated: SdkBatchTransactionStatus is deprecated. Please Use Chain.SdkBatchTransactionStatus() instead.
func (u *CoinUtil) SdkBatchTransactionStatus(hashListString string) (string, error) {
	return u.chain().BatchFetchTransactionStatus(hashListString), nil
}

func (u *CoinUtil) chain() *Chain {
	return NewChainWithRpc(u.RpcUrl)
}

func (u *CoinUtil) baseToken() base.Token {
	if u.IsMainCoin() {
		return u.chain().MainToken()
	} else {
		return u.chain().Erc20Token(u.ContractAddress)
	}
}

func (u *CoinUtil) ethToken() (TokenProtocol, error) {
	token, ok := u.baseToken().(TokenProtocol)
	if !ok {
		return nil, errors.New("golang type cast error") // TODO verify
	}
	return token, nil
}
