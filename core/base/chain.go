package base

type Chain interface {
	MainToken() Token

	BalanceOfAddress(address string) (*Balance, error)
	BalanceOfPublicKey(publicKey string) (*Balance, error)
	BalanceOfAccount(account Account) (*Balance, error)

	// Send the raw transaction on-chain
	// @return the hex hash string
	SendRawTransaction(signedTx string) (string, error)

	// Fetch transaction details through transaction hash
	FetchTransactionDetail(hash string) (*TransactionDetail, error)

	// Fetch transaction status through transaction hash
	FetchTransactionStatus(hash string) TransactionStatus

	// Batch fetch the transaction status, the hash list and the return value,
	// which can only be passed as strings separated by ","
	// @param hashListString The hash of the transactions to be queried in batches, a string concatenated with ",": "hash1,hash2,hash3"
	// @return Batch transaction status, its order is consistent with hashListString: "status1,status2,status3"
	BatchFetchTransactionStatus(hashListString string) string

	// Most chains can estimate the fee directly to the transaction object
	// **But two chains don't work: `aptos`, `starcoin`**
	EstimateTransactionFee(transaction Transaction) (fee *OptionalString, err error)

	// All chains can call this method to estimate the gas fee.
	// **Chain  `aptos`, `starcoin` must pass in publickey**
	EstimateTransactionFeeUsePublicKey(transaction Transaction, pubkey string) (fee *OptionalString, err error)

	// -----------------------------
	// polka
	// GetSignDataFromChain(t *Transaction, walletAddress string) ([]byte, error)

	// EstimateFeeForTransaction(transaction *Transaction) (s string, err error)

	// FetchScriptHashForMiniX(transferTo, amount string) (*MiniXScriptHash, error)
}
