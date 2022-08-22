package aptos

import (
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
)

const testnetRestUrl = "https://fullnode.devnet.aptoslabs.com"

func TestFaucet(t *testing.T) {
	account, _ := NewAccountWithMnemonic(testcase.M1)
	hashs, err := FaucetFundAccount(account.Address(), 2000, "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hashs.Value)
}

func TestTransafer(t *testing.T) {
	account, _ := NewAccountWithMnemonic(testcase.M1)
	toAddress := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
	amount := "100"

	chain := NewChainWithRestUrl(testnetRestUrl)
	token := NewToken(chain)

	balance, err := token.BalanceOfAccount(account)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(balance)

	signedTx, err := token.BuildTransferTxWithAccount(account, toAddress, amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signedTx.Value)

	txHash, err := chain.SendRawTransaction(signedTx.Value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestEstimateFee(t *testing.T) {
	account, _ := NewAccountWithMnemonic(testcase.M1)
	toAddress := "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015"
	amount := "100"

	chain := NewChainWithRestUrl(testnetRestUrl)
	token := NewToken(chain)

	fee, err := token.EstimateFees(account, toAddress, amount)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)
}

func TestFetchTransactionDetail(t *testing.T) {
	chain := NewChainWithRestUrl(testnetRestUrl)
	hash := "0xbe86d90b331d9db8002f7e435647bd108074cd4bd81ef0933ea3144ce82fc566"
	detail, err := chain.FetchTransactionDetail(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(detail)
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name    string
		restUrl string
		hash    string
		want    *base.TransactionDetail
		wantErr bool
	}{
		{
			name:    "testnet normal",
			restUrl: testnetRestUrl,
			hash:    "0xfd496b3dccae000096d4bf4aef581863ce2600c8867be9c2fe5b82a3408441af",
			want: &base.TransactionDetail{
				HashString:      "0xfd496b3dccae000096d4bf4aef581863ce2600c8867be9c2fe5b82a3408441af",
				FromAddress:     "0xcf4ddd208bbbbefb3227cafa5c917fc6541d26b1869276ea80d99ee0505fc6f8",
				ToAddress:       "0xcdbe33da8d218e97a9bec6443ba4a1b1858494f29142976d357f4770c384e015",
				Amount:          "100",
				EstimateFees:    "4",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1660556054,
			},
		},
		{
			name:    "testnet failed tx",
			restUrl: testnetRestUrl,
			hash:    "0x0a5720b278708820eb46c24af485858da8668e183a27ee57f3eed402cdda7436",
			want: &base.TransactionDetail{
				HashString:      "0x0a5720b278708820eb46c24af485858da8668e183a27ee57f3eed402cdda7436",
				FromAddress:     "0x56252ac5625573224bcaf89119b046f4a35a5c23bbf3d0f3eaa43311fbd2b2b3",
				ToAddress:       "0x903056ed3ddd9c7b9d5231ac96c8e6a218fe2a7cf26f17f04a96edb2cb832566",
				Amount:          "100",
				EstimateFees:    "4",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1660617399,
				FailureMessage:  "Move abort by ECOIN_STORE_NOT_PUBLISHED\n When an account hasn't registered `CoinStore` for `CoinType`.",
			},
		},
		{
			name:    "set address to hash",
			restUrl: testnetRestUrl,
			hash:    "0xcf4ddd208bbbbefb3227cafa5c917fc6541d26b1869276ea80d99ee0505fc6f8",
			wantErr: true,
		},
		{
			name:    "not transfer",
			restUrl: testnetRestUrl,
			hash:    "0xb7c06248b83bb7854d75f8d09a56ce4f5d7f799445fdb8781fccc536a01cd971",
			wantErr: true,
		},
		{
			name:    "mint tx",
			restUrl: testnetRestUrl,
			hash:    "0x6934afd26b2e371f69ed2095dab30961b4c5c4b40fca2351966cbcd6add96a69",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChainWithRestUrl(tt.restUrl)
			got, err := c.FetchTransactionDetail(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chain.FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain.FetchTransactionDetail() = %v, want %v", got, tt.want)
			}
		})
	}
}
