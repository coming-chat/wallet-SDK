package solana

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func SOL(amount float64) testcase.Amount {
	return testcase.Amount{Amount: amount, Multiple: 1e9}
}

func TestToken_BuildTransfer_SignedTransaction(t *testing.T) {
	account := M1Account(t)
	chain := TestnetChain()
	token := chain.MainToken()

	balance, err := token.BalanceOfAddress(account.Address())
	require.Nil(t, err)
	t.Log("sender address = ", account.Address())
	t.Log("balance = ", balance.Usable)

	txn, err := token.BuildTransfer(account.Address(), account.Address(), "100")
	require.Nil(t, err)

	gasfee, err := chain.EstimateTransactionFeeUsePublicKey(txn, account.PublicKeyHex())
	require.Nil(t, err)
	t.Log("Estimate fee = ", gasfee.Value)

	signedTxn, err := txn.SignedTransactionWithAccount(account)
	require.Nil(t, err)

	if false {
		hash, err := chain.SendSignedTransaction(signedTxn)
		require.Nil(t, err)
		t.Log("Transaction hash = ", hash.Value)
	}
}

func TestTransaction_New(t *testing.T) {
	raw := "GuxUPtW1ioB6PShC7aBtjuVqHXzHZzE9Xo3j81FmDSNs23C4tv6tGXih6sgaVdzCnKfv5sLL5LcFn19Y5CHTnZSuDD1byvEvisjqvCNoh4APExH6baTNZatEy8AQMBpXPA74n71fXwKAGQCkcYeqDTnoWkfvkpM7iWe16rFkFkKC241cMhNs6hZAjVJpeHeBVMJSLmvVmkJgaSGdia6wMdcosAfTC84CgbAEknuvpSERBEuTz5MwPw3dVirQiShp42xvMNkYodjQkgdf6qLywBteSAaq6zHHbScQi5eQL991QBgdHy6e49NEyK7E1RVhNqM3jHNWu6gQM4qf9VMtdzdoG3VRDSKoGYq7mxJjozfxk3orNuTjPASo6ixQX163W11zR3fHw28hmc3RxxwKZ1jjhQgbAFZmNoJUVqxxjWLYKjEjx7sizWoDA9UzryZQ7BVkdB4j2VhUzu78T5yYZmLzPkgKE84yRhxMqcHC6BnKTRoD4zyz8S4cdRM4Rz2aLUx339S9NFcdonj8wW1EByPHuFsPamvoLKgA2ecMxkaaJxsSMDpz73ppBfFjCFphnEHQwKgkY5KJYZHu26T1qMr6a9Q3koumYmNKy2gHDNnHvHGKk5xnfXQtPPmrur41tBHbRuH9BnWUYniNLxjZVVQU9MvGJdvPicphjnVSrc2onEYHYC29b5QXgQZGSUjAVFp5L3t7qnZ86SwtjmPris96h5NzTLTaCq3fo7Qr4aVGYF8MSbiyhduTeLhyAZczMrWbYtLpQBkbZNxs2BtHjHSLuxmbAcxEbniJQ3Gjf5LFDHvFPbKpJeoLYXhA8mwQizP7xTSi6v33J5nYkF2BmhneCXh8aME92NZVrZ8BDaEWHfGvhfmmqAdxwg3zwkRzpzGVZUx4ixiJMA3fFhziFs4VkXpdWug4hsRVuhRT1uMGh3yn5jcvhE7jgUjXwusD53tkzq7faotpkA1UWJeRNrzYJprHfstfnR6Kh8Nk1FgwjZqkj7bzDn6ftVGN2abLQsaMbi4XN3iYnJULYDERMAHwkobJW8apt5pobHApSLwJg6pQZbw2EpBN6N2zDrGSRetMiKjdq1zhU3jU7Aws1EjnSdbWKqxJ8j5KGJ8AUeTU9L2x9mLVBBdr4c18WkXfzJ173p6WUkpbeTtKvrS2XCZsaMN3wKV3zdUT9eWcQ5QQHoeDiKVeQjtiWctFZrLs5ZAshHUEJgbLorqw1zusL5vmwAZLonzQUc5C683P5v9AcSDG6D3Hjn235fNmfXaxYEeAc9NsXzB7RLFqk469TYteG1gHrBY9N5yGf5bTkZmiN4o12UssMvHGvJHefJbtH3iXQjWaPhYUzVPGPgq23CS3FVKDpk"
	txn, err := NewTransaction(raw)
	require.Nil(t, err)
	require.Equal(t, txn.Message.RecentBlockHash, "BXhGCzgqK32G6wZJtzBbMUn3HxZy5nAQyJ8CvuQ9jq8x")
	t.Log(txn)
}
