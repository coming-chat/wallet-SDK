package cosmos

import (
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type chainInfo struct {
	rpc  string
	rest string
	scan string
}

type chainConfig struct {
	cosmosProd chainInfo
	cosmosTest chainInfo
	terraProd  chainInfo
	terraTest  chainInfo
}

var rpcs = chainConfig{
	cosmosProd: chainInfo{
		rpc:  "https://cosmos-mainnet-rpc.allthatnode.com:26657",
		rest: "https://cosmos-mainnet-rpc.allthatnode.com:1317",
		scan: "https://www.mintscan.io/cosmos",
	},
	cosmosTest: chainInfo{
		rpc:  "https://cosmos-testnet-rpc.allthatnode.com:26657",
		rest: "https://cosmos-testnet-rpc.allthatnode.com:1317",
		scan: "https://cosmoshub-testnet.mintscan.io/cosmoshub-testnet",
	},
	terraProd: chainInfo{
		rpc:  "https://terra-mainnet-rpc.allthatnode.com:26657",
		rest: "https://terra-mainnet-rpc.allthatnode.com:1317",
		scan: "https://finder.terra.money",
	},
	terraTest: chainInfo{
		rpc:  "https://terra-bombay-rpc.allthatnode.com:26657",
		rest: "https://terra-bombay-rpc.allthatnode.com:1317",
		scan: "https://finder.terra.money/testnet",
	},
}

func (i *chainInfo) Chain() *Chain {
	return NewChainWithRpc(i.rpc, i.rest)
}

// $request cosmos1unek4dqvkwxv6sfrakk4903m0gmxkfyeprcqtg  theta

func TestTransfer(t *testing.T) {
	rpcinfo := rpcs.cosmosTest

	from := accountCase1.mnemonic
	account, _ := NewCosmosAccountWithMnemonic(from)

	toAddress := accountCase2.address
	gasPrice := GasPriceLow
	gasLimit := GasLimitDefault
	amount := "1000"

	chain := rpcinfo.Chain()
	token := chain.DenomToken(CosmosPrefix, CosmosCoinDenom)

	signedTx, err := token.BuildTransferTxWithAccount(account, toAddress, gasPrice, gasLimit, amount)
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := chain.SendRawTransaction(signedTx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash, "at", rpcinfo.scan)
}

func TestLunaTransfer(t *testing.T) {
	rpcinfo := rpcs.terraProd

	from := accountTerra.mnemonic
	account, _ := NewTerraAccountWithMnemonic(from)

	toAddress := accountTerra2.address
	gasPrice := TerraGasPrice
	gasLimit := TerraGasLimitDefault
	amount := "10000"

	chain := rpcinfo.Chain()
	token := chain.DenomToken(TerraPrefix, TerraLunaDenom)

	signedTx, err := token.BuildTransferTxWithAccount(account, toAddress, gasPrice, gasLimit, amount)
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := chain.SendRawTransaction(signedTx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash, "at", rpcinfo.scan)
}

func TestUSTTransfer(t *testing.T) {
	rpcinfo := rpcs.terraProd

	from := accountTerra.mnemonic
	account, _ := NewTerraAccountWithMnemonic(from)

	toAddress := accountTerra2.address
	gasPrice := TerraGasPriceUST
	gasLimit := TerraGasLimitDefault
	amount := "10000"

	chain := rpcinfo.Chain()
	token := chain.DenomToken(TerraPrefix, TerraUSTDenom)

	signedTx, err := token.BuildTransferTxWithAccount(account, toAddress, gasPrice, gasLimit, amount)
	if err != nil {
		t.Fatal(err)
	}

	txHash, err := chain.SendRawTransaction(signedTx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txHash, "at", rpcinfo.scan)
}

func TestChain_BalanceOfAddress(t *testing.T) {
	tests := []struct {
		name    string
		rpcinfo chainInfo
		address string
		denom   string
		wantErr bool
	}{
		{
			name:    "cosmos testnet normal",
			rpcinfo: rpcs.cosmosTest,
			address: accountCase1.address,
		},
		{
			name:    "cosmos mainnet normal",
			rpcinfo: rpcs.cosmosProd,
			address: "cosmos1lkw6n8efpj7mk29yvajpn9zue099l359cgzf0t",
			denom:   CosmosCoinDenom,
		},
		{
			name:    "cosmos testnet error address",
			rpcinfo: rpcs.cosmosTest,
			address: accountCase1.address + "s",
			denom:   CosmosCoinDenom,
			wantErr: true,
		},
		{
			name:    "cosmos mainnet error denom",
			rpcinfo: rpcs.cosmosProd,
			address: "cosmos1lkw6n8efpj7mk29yvajpn9zue099l359cgzf0t",
			denom:   "atom",
			wantErr: true,
		},
		{
			name:    "terra testnet luna normal",
			rpcinfo: rpcs.terraTest,
			address: accountTerra.address,
			denom:   TerraLunaDenom,
		},
		{
			name:    "terra mainnet ust normal",
			rpcinfo: rpcs.terraProd,
			address: "terra1dr7ackrxsqwmac2arx26gre6rj6q3sv29fnn7k",
			denom:   TerraUSTDenom,
		},
		{
			name:    "terra mainnet luna normal",
			rpcinfo: rpcs.terraProd,
			address: "terra1ncjg4a59x2pgvqy9qjyqprlj8lrwshm0wleht5",
			denom:   TerraLunaDenom,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.rpcinfo.Chain()
			got, err := chain.BalanceOfAddressAndDenom(tt.address, tt.denom)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t.Logf("balance is %s. We cannost assert the value, You can verify at %s/account/%s",
					got.Total, tt.rpcinfo.scan, tt.address)
			}
		})
	}
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name       string
		rpcInfo    chainInfo
		hash       string
		wantDetail *base.TransactionDetail
		wantErr    bool
	}{
		{
			name:    "cosmos testnet normal succeed tx",
			rpcInfo: rpcs.cosmosTest,
			hash:    "F068275DE4A4CC904D3E6A412A50DFACC235C62770BCD001E54E00BC4C17B1F0",
			wantDetail: &base.TransactionDetail{
				HashString:      "F068275DE4A4CC904D3E6A412A50DFACC235C62770BCD001E54E00BC4C17B1F0",
				FromAddress:     "cosmos19jwusy7lm8v5kqay8qjml79hs6e30t8j7ygm8r",
				ToAddress:       "cosmos10d2wkfl7y8rpgyxkcwa8urwt8muuc9aqcq9vys",
				Amount:          "10000",
				EstimateFees:    "800",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1652356672,
			},
		},
		{
			name:    "cosmos testnet normal failured tx",
			rpcInfo: rpcs.cosmosTest,
			hash:    "915B83334165EC73C5CF7B8FAD50E2B9165C1777C277B29A8D7FF6B2E4D6D96C",
			wantDetail: &base.TransactionDetail{
				HashString:      "915B83334165EC73C5CF7B8FAD50E2B9165C1777C277B29A8D7FF6B2E4D6D96C",
				FromAddress:     "cosmos19jwusy7lm8v5kqay8qjml79hs6e30t8j7ygm8r",
				ToAddress:       "cosmos10d2wkfl7y8rpgyxkcwa8urwt8muuc9aqcq9vys",
				Amount:          "1000",
				EstimateFees:    "500",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1652350345,
				FailureMessage:  "out of gas in location: WriteFlat; gasWanted: 60000, gasUsed: 60683: out of gas",
			},
		},
		{
			name:    "cosmos mainnet normal succeed tx",
			rpcInfo: rpcs.cosmosProd,
			hash:    "56C3314BF9FACA7238AA9C3BDDF622EEE8C2443BA92BC35E443860C7DE3F23AC",
			wantDetail: &base.TransactionDetail{
				HashString:      "56C3314BF9FACA7238AA9C3BDDF622EEE8C2443BA92BC35E443860C7DE3F23AC",
				FromAddress:     "cosmos17muvdgkep4ndptnyg38eufxsssq8jr3wnkysy8",
				ToAddress:       "cosmos1aftc0rwy2zk2ksmlrf3z4llnzlrwn69luwx3gs",
				Amount:          "5520000",
				EstimateFees:    "5000",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1652433492,
			},
		},
		{
			name:    "cosmos testnet normal error hash: odd length hex string",
			rpcInfo: rpcs.cosmosTest,
			hash:    "915B83334165EC73C5CF7B8FAD50E2B9165C1777C277B29A8D7FF6B2E4D6D96",
			wantErr: true,
		},
		{
			name:    "cosmos testnet normal error hash: not found",
			rpcInfo: rpcs.cosmosTest,
			hash:    "915B83334165EC73C5CF7B8FAD50E2B9165C1777C277B29A8D7FF6B2E4D6D9",
			wantErr: true,
		},
		{
			name:    "terra mainnet ust normal succeed tx",
			rpcInfo: rpcs.terraProd,
			hash:    "8C40B11BC4A242A6F0EEA187A37304CF7BB5DEC6E7DD0B6A912EE28E2B6F90B9",
			wantDetail: &base.TransactionDetail{
				HashString:      "8C40B11BC4A242A6F0EEA187A37304CF7BB5DEC6E7DD0B6A912EE28E2B6F90B9",
				FromAddress:     "terra1dr7ackrxsqwmac2arx26gre6rj6q3sv29fnn7k",
				ToAddress:       "terra1ncjg4a59x2pgvqy9qjyqprlj8lrwshm0wleht5",
				Amount:          "62085200000",
				EstimateFees:    "300000",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1652666974,
			},
		},
		{
			name:    "terra mainnet luna normal succeed tx",
			rpcInfo: rpcs.terraProd,
			hash:    "19771A22934641DBD3D347DCCAE939DAC37F39ABD88005AA735B8AAEA78599BA",
			wantDetail: &base.TransactionDetail{
				HashString:      "19771A22934641DBD3D347DCCAE939DAC37F39ABD88005AA735B8AAEA78599BA",
				FromAddress:     "terra1u868n8kekvez2lnrz44ca00ufzk78rux3sn8m2",
				ToAddress:       "terra1ncjg4a59x2pgvqy9qjyqprlj8lrwshm0wleht5",
				Amount:          "1584363636",
				EstimateFees:    "1500000",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1652671784,
			},
		},
		{
			name:    "terra mainnet ust normal failured tx",
			rpcInfo: rpcs.terraProd,
			hash:    "25ACF16526D3A4DEE5FE7C5CCEB597B5691134647829AD30CA1E36EDBEAC32B6",
			wantDetail: &base.TransactionDetail{
				HashString:      "25ACF16526D3A4DEE5FE7C5CCEB597B5691134647829AD30CA1E36EDBEAC32B6",
				FromAddress:     "terra17e2m3cg69klzdt6wpyhvlvmv3wsl9qwa38t7a4",
				ToAddress:       "terra1dr7ackrxsqwmac2arx26gre6rj6q3sv29fnn7k",
				Amount:          "414550000",
				EstimateFees:    "300000",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1652671459,
				FailureMessage:  "failed to execute message; message index: 0: 353650000uusd is smaller than 414550000uusd: insufficient funds",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.rpcInfo.Chain()
			gotDetail, err := chain.FetchTransactionDetail(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(gotDetail, tt.wantDetail) {
				t.Errorf("FetchTransactionDetail() gotDetail = %v, want %v", gotDetail, tt.wantDetail)
			}
		})
	}
}
