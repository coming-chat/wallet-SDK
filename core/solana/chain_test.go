package solana

import (
	"context"
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/stretchr/testify/require"
)

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
	if false {
		chain := DevnetChain()
		account := M1Account(t)
		txhash, err := chain.client().RequestAirdrop(context.Background(), account.Address(), 1e9)
		require.NoError(t, err)
		t.Log(txhash)
	}
}

func TestEstimateFee(t *testing.T) {
	chain := TestnetChain()
	token := &Token{chain: chain}

	receiver := "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g"
	amount := SOL(0.01).String()
	fee, err := token.EstimateFees(receiver, amount)
	t.Log(fee, err)
}

func TestBuildtxAndSendTransaction(t *testing.T) {
	sender := M1Account(t)
	chain := DevnetChain()
	token := NewToken(chain)

	receiver := "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g"
	amount := SOL(0.01).String()

	signedTx, err := token.BuildTransferTxWithAccount(sender, receiver, amount)
	require.NoError(t, err)

	if false {
		txHash, err := chain.SendRawTransaction(signedTx.Value)
		require.NoError(t, err)
		t.Log(txHash)
	}
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

func Test_FetchTransactionDetail(t *testing.T) {
	hash := "4T5vtEFpn5hpRMSaz1bMAKidfptaL3FZsp5ECv9RpB7WFfzeQfJWtPCBjYjqQv57xnmviqNMj2VMbQ7kzoeC2e4p" // spl token
	// hash := "JMnhxFzF6hceB3zLw6zvHy5RHqVJo2anAVXumaxvZe1FDFGpxEWwbjyNHBGMHMxXc4uX69KGsmVMZgbGWua7YQr" // SOL
	chain := DevnetChain()

	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	t.Log(detail.JsonString())
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
			hash: "4xNjdnHufbsVVgyQxFDzcGsC6ypBntCZoaZgXziwxh2gTXERgUvVvQA1KLGQMGCpJTKctSswZPuCA1DsLii45Jwr",
			want: &base.TransactionDetail{
				HashString:      "4xNjdnHufbsVVgyQxFDzcGsC6ypBntCZoaZgXziwxh2gTXERgUvVvQA1KLGQMGCpJTKctSswZPuCA1DsLii45Jwr",
				FromAddress:     "4MPScMzmKwQfpzQ4MtkSaqKQbTEzGsWqovUMweNz7nFo",
				ToAddress:       "9B5XszUGdMaxCZ7uSQhPzdks5ZQSmWxrmzCSvtJ6Ns6g",
				Amount:          "10000000",
				EstimateFees:    "5000",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1690880972,
			},
		},
		{
			name:    "not found",
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
			hash: "4PtFQcC6WUorchQgxGRzzDpqRSQhkM8ZDtovBWi8nYcAtsckY35XZnAp1rarH1WVMqfMzArbkzXJwBCbhexzRzAJ",
			want: base.TransactionStatusSuccess,
		},
		{
			name: "devnet normal2",
			rpc:  DevnetRPCEndpoint,
			hash: "4xNjdnHufbsVVgyQxFDzcGsC6ypBntCZoaZgXziwxh2gTXERgUvVvQA1KLGQMGCpJTKctSswZPuCA1DsLii45Jwr",
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
