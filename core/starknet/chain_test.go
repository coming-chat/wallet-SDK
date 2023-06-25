package starknet

import (
	"context"
	"testing"

	"github.com/dontpanicdao/caigo/gateway"
	"github.com/stretchr/testify/require"
)

func MainnetChain() *Chain {
	c, _ := NewChainWithRpc(BaseRpcUrlMainnet, NetworkMainnet)
	return c
}
func GoerliChain() *Chain {
	c, _ := NewChainWithRpc(BaseRpcUrlGoerli, NetworkGoerli)
	return c
}

func TestBalance(t *testing.T) {
	owner := "0x0023C4475F2f2355580f5994294997d3A18237ef62223D20C41876556327A05E"
	chain := GoerliChain()

	balance, err := chain.BalanceOf(owner, ETHTokenAddress)
	require.Nil(t, err)
	t.Log(balance.Total)
}

func TestDeployAccount(t *testing.T) {
	acc := M1Account(t)
	chain := GoerliChain()

	txn, err := chain.BuildDeployAccountTransaction(acc.PublicKeyHex())
	require.Nil(t, err)

	signedTxn, err := txn.SignedTransactionWithAccount(acc)
	require.Nil(t, err)

	hash, err := chain.SendSignedTransaction(signedTxn)
	require.Nil(t, err)

	t.Log(hash.Value)
}

func TestTransfer(t *testing.T) {
	acc := M1Account(t)
	chain := GoerliChain()

	token, err := chain.NewToken(ETHTokenAddress)
	require.Nil(t, err)

	txn, err := token.BuildTransfer(acc.Address(), acc.Address(), "10000000")
	require.Nil(t, err)

	gasFee, err := chain.EstimateTransactionFeeUseAccount(txn, acc)
	require.Nil(t, err)
	t.Log(gasFee.Value)

	// signedTxn, err := txn.SignedTransactionWithAccount(acc)
	// require.Nil(t, err)
	// hash, err := chain.SendSignedTransaction(signedTxn)
	// require.Nil(t, err)
	// t.Log(hash.Value)
}

func TestNonce(t *testing.T) {
	account := M1Account(t)
	chain := GoerliChain()

	address := account.Address()

	nonce, err := chain.gw.Nonce(context.Background(), address, "latest")
	require.Nil(t, err)
	t.Log(nonce.String())
}

func TestFetchTransactionDetail(t *testing.T) {
	chain := GoerliChain()
	hash := "0x01de50b64326c02a9830ea7bf825224103dd3ea4426309514039a01eaadcf5a4"

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	t.Log(detail)
}

func TestTransactionInfo(t *testing.T) {
	chain := GoerliChain()

	hash := "0x47ba4a447e929987094289c27ecc3d37b5a02e580835083d87247c6a97a4e00"

	txn, err := chain.gw.Transaction(context.Background(), gateway.TransactionOptions{
		TransactionHash: hash,
	})
	require.Nil(t, err)

	block, err := chain.gw.Block(context.Background(), &gateway.BlockOptions{
		BlockHash: txn.BlockHash,
	})
	require.Nil(t, err)

	receipt, err := chain.gw.TransactionReceipt(context.Background(), hash)
	require.Nil(t, err)

	status, err := chain.gw.TransactionStatus(context.Background(), gateway.TransactionStatusOptions{
		TransactionHash: hash,
	})
	require.Nil(t, err)

	t.Log(txn)
	t.Log(block)
	t.Log(receipt)
	t.Log(status)
}

func TestFetchTransactionStatus(t *testing.T) {
	chain := GoerliChain()
	hash := "0x03ae12fb58a3f4a6dcd7d04ad10c4d3b2ab97d23ee167a6109db719ba703eed9"

	status := chain.FetchTransactionStatus(hash)
	t.Log(status)
}

// func TestGraphql(t *testing.T) {
// 	graphqlUrl := "https://api-testnet.starkscan.co/graphql"
// 	query := map[string]any{
// 		"query": "query TransactionPageTabs_TransactionQuery(\\n  $input: TransactionInput\u0021\\n) {\\n  transaction(input: $input) {\\n    transaction_hash\\n    transaction_status\\n    number_of_events\\n    number_of_message_logs\\n    ...TransactionPageOverviewTab\\n    id\\n  }\\n}\\n\\nfragment TransactionActualFeesItem_transaction on Transaction {\\n  actual_fee\\n  actual_fee_display\\n  erc20_transfer_events {\\n    id\\n    call_invocation_type\\n  }\\n  ...TransactionActualFeesTransferredItems_transaction\\n}\\n\\nfragment TransactionActualFeesTransferredItem on ERC20TransferEvent {\\n  id\\n  from_address\\n  from_erc20_identifier\\n  transfer_amount_display\\n  transfer_to_address\\n  transfer_to_identifier\\n  call_invocation_type\\n}\\n\\nfragment TransactionActualFeesTransferredItems_transaction on Transaction {\\n  erc20_transfer_events {\\n    id\\n    ...TransactionActualFeesTransferredItem\\n  }\\n}\\n\\nfragment TransactionCalldataItem_transaction on Transaction {\\n  entry_point_selector_name\\n  calldata\\n  calldata_decoded\\n  entry_point_selector\\n  initiator_address\\n  initiator_identifier\\n  main_calls {\\n    selector_name\\n    calldata_decoded\\n    selector\\n    calldata\\n    contract_address\\n    contract_identifier\\n    id\\n  }\\n}\\n\\nfragment TransactionConstructorCalldataItem_transaction on Transaction {\\n  entry_point_selector_name\\n  calldata_decoded\\n  entry_point_selector\\n  constructor_calldata\\n  initiator_address\\n  initiator_identifier\\n}\\n\\nfragment TransactionDeployedContractsItem_transaction on Transaction {\\n  deployed_contracts {\\n    id\\n    contract_address\\n    contract_identifier\\n  }\\n}\\n\\nfragment TransactionExecutionResourcesItem_transaction on Transaction {\\n  execution_resources {\\n    execution_resources_n_steps\\n    execution_resources_n_memory_holes\\n    execution_resources_builtin_instance_counter {\\n      name\\n      value\\n    }\\n  }\\n}\\n\\nfragment TransactionNFTEventsItem_transaction on Transaction {\\n  nft_events {\\n    id\\n    type\\n    nft_contract_address\\n    nft_contract_nft_identifier\\n    nft_token_id\\n    from_address\\n    from_identifier\\n    to_address\\n    to_identifier\\n    nft {\\n      image_small_url\\n      name\\n      id\\n    }\\n  }\\n}\\n\\nfragment TransactionPageOverviewTab on Transaction {\\n  transaction_hash\\n  block_hash\\n  block_number\\n  transaction_status\\n  timestamp\\n  transaction_type\\n  contract_address\\n  contract_identifier\\n  sender_address\\n  sender_identifier\\n  class_hash\\n  entry_point_selector\\n  max_fee\\n  max_fee_display\\n  nonce\\n  ...TransactionDeployedContractsItem_transaction\\n  ...TransactionActualFeesItem_transaction\\n  ...TransactionTokensTransferredItem_transaction\\n  ...TransactionCalldataItem_transaction\\n  ...TransactionConstructorCalldataItem_transaction\\n  ...TransactionSignatureItem_transaction\\n  ...TransactionExecutionResourcesItem_transaction\\n  ...TransactionNFTEventsItem_transaction\\n}\\n\\nfragment TransactionSignatureItem_transaction on Transaction {\\n  signature\\n}\\n\\nfragment TransactionTokensTransferredItem_transaction on Transaction {\\n  erc20_transfer_events {\\n    id\\n    from_address\\n    from_erc20_identifier\\n    transfer_amount_display\\n    transfer_to_address\\n    transfer_to_identifier\\n    transfer_from_address\\n    transfer_from_identifier\\n    call_invocation_type\\n  }\\n}\\n",
// 		"variables": map[string]any{
// 			"input": map[string]string{
// 				"transaction_hash": "0x01982805b13bb3d661c0015df15c210a6166cbad720e554dbf62f2106102c849"},
// 		},
// 	}
// 	body, err := json.Marshal(query)
// 	require.Nil(t, err)
// 	params := httpUtil.RequestParams{
// 		Header: map[string]string{
// 			"Content-Type": "application/json",
// 			"Authority":    "api-testnet.starkscan.co",
// 			"Accept":       "application/json",
// 			"Origin":       "https://testnet.starkscan.co",
// 			"Referer":      "https://testnet.starkscan.co/",
// 			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
// 		},
// 		Body:    body,
// 		Timeout: time.Duration(20 * int64(time.Second)),
// 	}
// 	response, err := httpUtil.Post(graphqlUrl, params)
// 	require.Nil(t, err)
// 	t.Log(string(response)) // 403 forbidden
// }
