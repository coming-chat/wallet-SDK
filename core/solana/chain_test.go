package solana

import (
	"context"
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
)

func newChainAndAccount() (*Chain, *Account) {
	// chain := NewChainWithRpc(rpc.MainnetRPCEndpoint)
	chain := NewChainWithRpc(DevnetRPCEndpoint)
	// c := client.NewClient(rpc.LocalnetRPCEndpoint)
	account, _ := NewAccountWithMnemonic(testcase.M1)
	return chain, account
}

func DevnetChain() *Chain {
	return NewChainWithRpc(DevnetRPCEndpoint)
}
func TestnetChain() *Chain {
	return NewChainWithRpc(TestnetRPCEndpoint)
}
func MainnetChain() *Chain {
	return NewChainWithRpc(MainnetRPCEndpoint)
}

func TestAirdrop(t *testing.T) {
	chain, account := newChainAndAccount()
	txhash, err := chain.client().RequestAirdrop(context.Background(), account.Address(), 1e9)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(txhash)
}

func TestEstimateFee(t *testing.T) {
	chain := TestnetChain()
	token := &Token{chain: chain}

	receiver := "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g"
	amount := "10000000"
	fee, err := token.EstimateFees(receiver, amount)
	t.Log(fee, err)
}

func TestBuildtxAndSendTransaction(t *testing.T) {
	chain, acc := newChainAndAccount()
	token := &Token{chain: chain}

	receiver := "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g"
	amount := "10000000"
	// amount := "1879985000"

	signedTx, err := token.BuildTransferTxWithAccount(acc, receiver, amount)
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := chain.SendRawTransaction(signedTx.Value)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash)
}

func TestChain_BalanceOfAddress(t *testing.T) {
	// https://explorer.solana.com/address/AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY
	// https://explorer.solana.com/address/AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY?cluster=devnet
	scan := "https://explorer.solana.com"
	tests := []struct {
		name    string
		rpc     string
		address string
		wantErr bool
	}{
		{
			name:    "mainnet empty account",
			rpc:     MainnetRPCEndpoint,
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY",
		},
		{
			name:    "devnet normal account",
			rpc:     DevnetRPCEndpoint,
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY",
		},
		{
			name:    "invalid address",
			rpc:     DevnetRPCEndpoint,
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqP",
			wantErr: false, //
		},
		{
			name:    "invalid address base58 char 0OIl",
			rpc:     DevnetRPCEndpoint,
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXq0O",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChainWithRpc(tt.rpc)
			got, err := c.BalanceOfAddress(tt.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chain.BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("Chain.BalanceOfAddress() We cannot assert the balance %v, maybe you can check at %v/address/%v", got.Total, scan, tt.address)
			}
		})
	}
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name    string
		rpc     string
		hash    string
		want    *base.TransactionDetail
		wantErr bool
	}{
		{
			name: "mainnet succeed transfer",
			rpc:  MainnetRPCEndpoint,
			hash: "2RFRkcx8biPpfrSeZvtiWVihGWLrK5GF9J2GnD8xJvR5sGrF7jsVZHQcXpqCLSZEmP2zi7PqUngz5W6mfeDKNy9w",
			want: &base.TransactionDetail{
				HashString:      "2RFRkcx8biPpfrSeZvtiWVihGWLrK5GF9J2GnD8xJvR5sGrF7jsVZHQcXpqCLSZEmP2zi7PqUngz5W6mfeDKNy9w",
				FromAddress:     "AXUChvpRwUUPMJhA4d23WcoyAL7W8zgAeo7KoH57c75F",
				ToAddress:       "AXUChvpRwUUPMJhA4d23WcoyAL7W8zgAeo7KoH57c75F",
				Amount:          "7360",
				EstimateFees:    "5000",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1657249797,
			},
		},
		{
			name: "devnet succeed transfer",
			rpc:  DevnetRPCEndpoint,
			hash: "2CLc8VSHsS67JT6a4UQZMMwFzcgnriHL4d1rpwu9WBwAqzucj6XPBcL2AYJy7n6xvmrnXTgGRKoThHizN8E8NTFN",
			want: &base.TransactionDetail{
				HashString:      "2CLc8VSHsS67JT6a4UQZMMwFzcgnriHL4d1rpwu9WBwAqzucj6XPBcL2AYJy7n6xvmrnXTgGRKoThHizN8E8NTFN",
				FromAddress:     "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY",
				ToAddress:       "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g",
				Amount:          "100000000",
				EstimateFees:    "5000",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1657177658,
			},
		},
		{
			name:    "devnet succeed but not contain an transfer",
			rpc:     DevnetRPCEndpoint,
			hash:    "3VjZjLrinNbHnkoTcvFi37nZgBcLdpCeHtmuXxqWS21stBbfKMCNhqmtG46BpPnWav16zPjNoSgM2eDn6w9k6bDN",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChainWithRpc(tt.rpc)
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

func TestChain_FetchTransactionStatus(t *testing.T) {
	tests := []struct {
		name string
		rpc  string
		hash string
		want base.TransactionStatus
	}{
		{
			name: "devnet normal1",
			rpc:  DevnetRPCEndpoint,
			hash: "3L8qPd9cTFj3KSe8bBhjdopik2xb7m2gE1ji6oWpZWgu6RgAKxA7Z8b6o11uBSYc8MmwFpeE4EoRtd2q14d4ePWe",
			want: base.TransactionStatusSuccess,
		},
		{
			name: "devnet normal2",
			rpc:  DevnetRPCEndpoint,
			hash: "xV56yXtPnxnzgvf5uNp3es5znqNt2FerXiVxWvEDZxLzsNkK2zy6vUYZhcLLQYbswfpbUBXnmdVnDZ3dAR5YSYm",
			want: base.TransactionStatusSuccess,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewChainWithRpc(tt.rpc)
			if got := c.FetchTransactionStatus(tt.hash); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chain.FetchTransactionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
