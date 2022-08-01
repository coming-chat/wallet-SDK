package eth

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
)

func TestRedPacketProcess(t *testing.T) {
	chain := rpcs.sherpaxProd.Chain()
	account, _ := NewAccountWithMnemonic(testcase.M1)
	// account, _ := EthAccountWithPrivateKey("0x......")
	contractAddress := "0x4777abDEc6D52C25b4bc55a361da495011ccDBC3"

	balance, err := chain.BalanceOfAddress(account.Address())
	t.Logf("the account's balance %v %v", balance, err)

	// 1.1 Create red packet
	erc20Token := "0xa10AF02fD7eD3B5FF107B57bB1068a3f54BcAE92" // erc20 PCX
	count := 3
	amount := "100000000"
	action, err := NewRedPacketActionCreate(erc20Token, count, amount)

	// 1.2 Open red packet
	// packetId := int64(0)
	// addresses := []string{"0x99e5f4759fC07ee8F4f1B5a017ba100EFFC0C9C0"}
	// amounts := []string{"50000000"}
	// action, err := NewRedPacketActionOpen(packetId, addresses, amounts)

	// 1.3 Close red packet
	// packetId := int64(0)
	// creator := "0x6334d64D5167F726d8A44f3fbCA66613708E59E7"
	// action, err := NewRedPacketActionClose(packetId, creator)

	if err != nil {
		t.Fatal(err)
	}

	// 2. ensure erc20 coin approved
	// 如果是发红包，那这一步是必须的
	txhash, err := action.EnsureApprovedTokens(account, chain, contractAddress, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Are there new approve? %v", txhash)

	// 3. make transaction
	transaction, err := action.TransactionFrom(account.Address(), contractAddress, chain)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Service Fee = %v", transaction.TotalAmount())

	// 4.1 signing raw tx with privatekey
	// rawTx, err := chain.BuildTransferTx(account.PrivateKeyHex(), transaction)
	// 4.2 or signing raw tx with account
	rawTx, err := chain.BuildTransferTxWithAccount(account, transaction)
	if err != nil {
		t.Fatal(err)
	}

	// 5. send transaction
	txHash, err := chain.SendRawTransaction(rawTx.Value)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Red packet process success! txHash = %v", txHash)
}

func TestFetchRedPacketCreationDetail(t *testing.T) {
	hash := "0x598ed72d6ddcc1a4b378acd9b6d1917dc0eea0eb905de4aca27ce50e61b1539c"
	chain := rpcs.sherpaxProd.Chain()
	detail, err := chain.FetchRedPacketCreationDetail(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(detail, detail.TransactionDetail)

	jsonString := detail.JsonString()
	t.Log("json string = ", jsonString)

	model2, err := NewRedPacketDetailWithJsonString(jsonString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(model2, model2.TransactionDetail)
}
