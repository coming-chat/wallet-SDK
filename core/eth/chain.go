package eth

import "github.com/coming-chat/wallet-SDK/core/base"

type Chain struct {
	RpcUrl string
}

func NewChainWithRpc(rpcUrl string) *Chain {
	return &Chain{
		RpcUrl: rpcUrl,
	}
}

// MARK - Implement the protocol Chain

func (c *Chain) MainToken() base.Token {
	return &Token{chain: c}
}

func (c *Chain) MainEthToken() TokenProtocol {
	return &Token{chain: c}
}

func (c *Chain) Erc20Token(contractAddress string) TokenProtocol {
	return &Erc20Token{
		Token:           &Token{chain: c},
		ContractAddress: contractAddress,
	}
}

func (c *Chain) BalanceOfAddress(address string) (*base.Balance, error) {
	b := base.EmptyBalance()

	eip55Address, err := TransformEIP55Address(address)
	if err != nil {
		return b, err
	}

	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return b, err
	}

	balance, err := chain.Balance(eip55Address)
	if err != nil {
		return b, err
	}
	return &base.Balance{
		Total:  balance,
		Usable: balance,
	}, nil
}

func (c *Chain) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return c.BalanceOfAddress(publicKey)
}

func (c *Chain) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return c.BalanceOfAddress(account.Address())
}

// Send the raw transaction on-chain
// @return the hex hash string
func (c *Chain) SendRawTransaction(signedTx string) (string, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return "", err
	}
	return chain.SendRawTransaction(signedTx)
}

// Fetch transaction details through transaction hash
func (c *Chain) FetchTransactionDetail(hash string) (*base.TransactionDetail, error) {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return nil, err
	}
	return chain.FetchTransactionDetail(hash)
}

// Fetch transaction status through transaction hash
func (c *Chain) FetchTransactionStatus(hash string) base.TransactionStatus {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return base.TransactionStatusNone
	}
	return chain.FetchTransactionStatus(hash)
}

// Batch fetch the transaction status, the hash list and the return value,
// which can only be passed as strings separated by ","
// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
func (c *Chain) BatchFetchTransactionStatus(hashListString string) string {
	chain, err := GetConnection(c.RpcUrl)
	if err != nil {
		return ""
	}
	return chain.SdkBatchTransactionStatus(hashListString)
}
