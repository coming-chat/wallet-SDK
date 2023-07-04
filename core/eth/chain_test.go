package eth

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestZK(t *testing.T) {

	// rpc := "https://testnet.era.zksync.dev"
	// hash := "0x590fae06853d27495f1c5d4d23ea35ad910454a30de012231cb390485bc8e268"

	rpc := "https://mainnet.era.zksync.io"
	hash := "0xc9ac881431a8df988bdcc00693e2c286848ce12f48e6af0323588cb68d822417"

	chain := NewChainWithRpc(rpc)
	detail, err := chain.FetchTransactionDetail(hash)
	require.Nil(t, err)
	t.Log(detail.Status)
}

func TestChainDetail(t *testing.T) {
	hash := "0xb5f9b89667d854555ce1b178a025c44d869dfcdef24ae4c5a1e790c5f41bbafb"
	chain := rpcs.ethereumProd.Chain()

	// hash := "0xa2b05fc52748ee21845c8596263d75d85f54c522575be2e399235d144cd5a153"
	// chain := rpcs.optimismProd.Chain()

	// hash := "0x740cc93867cd3b84928ba9ae65a6511368f856e180e82fbf2649cca540d6db2c"
	// chain := rpcs.optimismTest.Chain()

	detail, err := chain.FetchTransactionDetail(hash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(detail)
}

func TestChain_BalanceOfAddress(t *testing.T) {
	addressBlackHole := "0x0000000000000000000000000000000000000000"
	tests := []struct {
		name    string
		net     rpcInfo
		address string
		wantErr bool
	}{
		{
			name:    "eth normal",
			net:     rpcs.ethereumProd,
			address: "0x62c3aF16954fba6D920835ec56f7b63139daAa6e",
		},
		{
			name:    "eth black hole",
			net:     rpcs.ethereumProd,
			address: addressBlackHole,
		},
		{
			name:    "binance-prod normal",
			net:     rpcs.binanceProd,
			address: "0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2",
		},
		{
			name:    "binance-prod error address altered one char", // but is can queryed
			net:     rpcs.binanceProd,
			address: "0x62c3aF16954fba6D920835ec56f7b63139daAa6d",
		},
		{
			name:    "eth error eip55 address", // but is can queryed
			net:     rpcs.ethereumProd,
			address: "0x62c3aF16954fba6D920835ec56f7b63139daAA6E",
		},
		{
			name:    "eth error address short",
			net:     rpcs.ethereumProd,
			address: "0x62c3aF16954fba6D920835ec56f",
			wantErr: true,
		},
		{
			name:    "optmism prod",
			net:     rpcs.optimismProd,
			address: addressBlackHole,
		},
		{
			name:    "optimism test",
			net:     rpcs.optimismTest,
			address: addressBlackHole,
		},
		{
			name:    "arbitrum prod",
			net:     rpcs.arbitrumProd,
			address: addressBlackHole,
		},
		{
			name:    "kcc prod",
			net:     rpcs.kccProd,
			address: addressBlackHole,
		},
		{
			name:    "avax prod",
			net:     rpcs.avaxProd,
			address: addressBlackHole,
		},
		{
			name:    "avax test",
			net:     rpcs.avaxTest,
			address: addressBlackHole,
		},
		{
			name:    "polygon prod",
			net:     rpcs.polygonProd,
			address: addressBlackHole,
		},
		{
			name:    "polygon test",
			net:     rpcs.polygonTest,
			address: addressBlackHole,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.net.Chain().BalanceOfAddress(tt.address)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("BalanceOfAddress() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			totalFloat, _ := big.NewFloat(0).SetString(got.Total)
			totalFloat.Quo(totalFloat, big.NewFloat(1e18))
			t.Logf("BalanceOfAddress() balance ≈ %v, full = %v", totalFloat.String(), got.Total)
			t.Log("Unable to verify balance, maybe you should check with this address which may be useful: " + tt.net.scan + "/address/" + tt.address)
		})
	}
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	type args struct {
		on   rpcInfo
		hash string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.TransactionDetail
		wantErr bool
	}{
		{
			name: "sherpax-prod contract token USB succeed",
			args: args{rpcs.sherpaxProd, "0x004eaae28f7a7f947c6e8a159f4b74a3aa719248ca4a47d9e5bbf32c394b460f"},
			want: &base.TransactionDetail{
				HashString:      "0x004eaae28f7a7f947c6e8a159f4b74a3aa719248ca4a47d9e5bbf32c394b460f",
				Amount:          "100000000000000000",
				EstimateFees:    "22895100000000000",
				FromAddress:     "0x3AA9c65C4393920E46B2B022841f3EaB7f49f7BC",
				ToAddress:       "0x5fD3d526A946DdB67C810f1F1C4A8c9214da17ef",
				Status:          base.TransactionStatusSuccess,
				FinishTimestamp: 1647598998,
			},
		},
		{
			name: "optimism prod erc20 token failured execution reverted",
			args: args{rpcs.optimismProd, "0x13dfd70e710e8451cf88cf8bd55b02a525a45efe028309a019defe5ffc9d5e83"},
			want: &base.TransactionDetail{
				HashString:      "0x13dfd70e710e8451cf88cf8bd55b02a525a45efe028309a019defe5ffc9d5e83",
				Amount:          "38919826",
				EstimateFees:    "202882647584248",
				FromAddress:     "0x8F1c69De5E086BA1E441707B9cbD94860529beE4",
				ToAddress:       "0xE56BD3FfC787942F8aB9cf20D2D650E3184aCCc3",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1654140210,
				FailureMessage:  "execution reverted",
			},
		},
		{
			name: "binance-test BNB failured out of gas",
			args: args{rpcs.binanceTest, "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c7"},
			want: &base.TransactionDetail{
				HashString:      "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c7",
				Amount:          "1000000000000000000",
				EstimateFees:    "378000000000000",
				FromAddress:     "0xaa25Aa7a19f9c426E07dee59b12f944f4d9f1DD3",
				ToAddress:       "0x6cd2Bf22B3CeaDfF6B8C226487265d81164396C5",
				Status:          base.TransactionStatusFailure,
				FinishTimestamp: 1649840877,
				FailureMessage:  "err: insufficient funds for gas * price + value: address 0xaa25Aa7a19f9c426E07dee59b12f944f4d9f1DD3 have 0 want 378000000000000 (supplied gas 21000)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.args.on.Chain()
			got, err := chain.FetchTransactionDetail(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchTransactionDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain_FetchTransactionDetail_Cover_Multi_Rpcs(t *testing.T) {
	type args struct {
		on   rpcInfo
		hash string
	}
	tests := []struct {
		name     string
		args     args
		wantTime int64 // Pass if timestamp is correct
		wantErr  bool
	}{
		{
			name:     "ethereum contract USDT succeed",
			args:     args{rpcs.ethereumProd, "0xdde5debd0856f295f85f6ee9ae13c84519fb01ef02cb85fa0e5194b618857898"},
			wantTime: 1650816611,
		},
		{
			name:     "sherpax-test KSX success",
			args:     args{rpcs.sherpaxTest, "0xf8ce72bbdfaa51a94bfe34a8bb8f973e35a0e4927711d9b247fb519ad5fb36ed"},
			wantTime: 1650872358,
		},
		{
			name:     "sherpax-prod contract token USB succeed",
			args:     args{rpcs.sherpaxProd, "0x004eaae28f7a7f947c6e8a159f4b74a3aa719248ca4a47d9e5bbf32c394b460f"},
			wantTime: 1647598998,
		},
		{
			name:     "binance-test BNB failured out of gas",
			args:     args{rpcs.binanceTest, "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c7"},
			wantTime: 1649840877,
		},
		{
			name:     "binance-prod contract USDT succeed",
			args:     args{rpcs.binanceProd, "0xdb3fd18da286f1a2fa90ec9b224b0eb0163a2ab98b9e701da66dc3d8d1d38c14"},
			wantTime: 1651217162,
		},
		{
			name:     "sherpax test EIP1559 succeed",
			args:     args{rpcs.sherpaxTest, "0x32afbd65fe73dda3734cb9419a87fccb861633db086f7ca487c044a8112dcbe7"},
			wantTime: 1652150106,
		},
		{
			name:     "ethereum EIP1559 transfer succeed",
			args:     args{rpcs.ethereumProd, "0xbba7921fda55dae8423d105d74acdbfea346f15969f6c924552d0b062128f271"},
			wantTime: 1651215357,
		},
		{
			name:    "binance-prod error hash",
			args:    args{rpcs.binanceProd, "0x5841f924fd76434f7f17ef8faf192142dbb5a363b9295eda0cc9f385e22399c6"},
			wantErr: true,
		},
		{
			name:    "ethereum-prod error hash",
			args:    args{rpcs.binanceProd, "0x5841f924fd76434f7f17ef8faf19214"},
			wantErr: true,
		},
		{
			name:     "optimism prod succeed",
			args:     args{rpcs.optimismProd, "0xda38aaaaa858fb65f62a41455308e71a57cc6c5a1c647d7f80ba316362a5a31c"},
			wantTime: 1654133507,
		},
		{
			name:     "arbitrum prod erc20 succeed",
			args:     args{rpcs.arbitrumProd, "0xc9b7e00273af851237f4cd76570da81942aca9b163044c6b7b9d09a46e17338b"},
			wantTime: 1653897408,
		},
		{
			name:     "kcc prod succeed",
			args:     args{rpcs.kccProd, "0xb118c7957aacf4c63c8b723776ade76fd77d5411ea799741ce9edf80d6a5739f"},
			wantTime: 1654119340,
		},
		{
			name:     "avax prod succeed transfer",
			args:     args{rpcs.avaxProd, "0x4e7ec461ab759c1b6f042c6676ef4ec928cc1f8adabb302df8aa44ebd314a6d0"},
			wantTime: 1656820690,
		},
		{
			name:     "polygon prod succeed transfer",
			args:     args{rpcs.polygonProd, "0xe26d28e25bdb82f67e993ac8491249707361f044f6de5ba3d47d5b41ab8feb3e"},
			wantTime: 1656922893,
		},
		{
			name:     "polygon prod erc20 transfer",
			args:     args{rpcs.polygonProd, "0x02b1d1e48584dfffa89d17e87548d8e2dd0b882942a909c6955ac2961262ac11"},
			wantTime: 1656932999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.args.on.Chain()
			got, err := chain.FetchTransactionDetail(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.FinishTimestamp != tt.wantTime {
				t.Errorf("FetchTransactionDetail() got = %v, want %v", got, tt.wantTime)
			}
		})
	}
}

func TestChain_CallContract(t *testing.T) {
	msg := NewCallMsg()
	msg.SetTo("0x37088186089c7d6bcd556d9a15087dfae3ba0c32")
	msg.SetDataHex("0x70a082310000000000000000000000008de5ff2eded4d897da535ab0f379ec1b9257ebab")

	tests := []struct {
		name    string
		on      rpcInfo
		msg     *CallMsg
		block   int64
		wantErr bool
	}{
		{
			name:  "binance prod pending",
			on:    rpcs.binanceProd,
			msg:   msg,
			block: -1,
		},
		{
			name:  "binance prod latest",
			on:    rpcs.binanceProd,
			msg:   msg,
			block: -2,
		},
		{
			name:    "sherpax prod pending",
			on:      rpcs.sherpaxProd,
			msg:     msg,
			block:   -1,
			wantErr: true, // sherpax not support pending call
		},
		{
			name:  "sherpax prod latest",
			on:    rpcs.sherpaxProd,
			msg:   msg,
			block: -2,
		},
		{
			name:  "sherpax prod normal block",
			on:    rpcs.sherpaxProd,
			msg:   msg,
			block: 1670400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := tt.on.Chain()
			got, err := chain.CallContract(tt.msg, tt.block)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t.Logf("CallContract() got = %v", got)
			}
		})
	}
}

func TestEthCall(t *testing.T) {
	chain := rpcs.sherpaxProd.Chain()

	msg := &CallMsg{}
	msg.SetFrom("0x8de5ff2eded4d897da535ab0f379ec1b9257ebab")
	msg.SetTo("0xf4ffbd250415d12bb5aa498cce28ecbe85fb7141")
	msg.SetValueHex("0x16345785d8a0000")
	msg.SetDataHex("0x7ff36ab5000000000000000000000000000000000000000000000000000000000121887000000000000000000000000000000000000000000000000000000000000000800000000000000000000000008de5ff2eded4d897da535ab0f379ec1b9257ebab0000000000000000000000000000000000000000000000000000000062820dbe0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000900639cc5a37c519c9e32bfa7eadf747c53d9b0a00000000000000000000000091aac463b5473bde2fdd910f09101ed687bcb97a")

	res, err := chain.LatestCallContract(msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}

func TestTransferOptimism(t *testing.T) {
	from, _ := NewAccountWithMnemonic(accountCase1.mnemonic)
	to := "0x8De5fF2edeD4d897da535ab0F379Ec1b9257eBaB"
	amount := "10000000000000" // 0.00001
	t.Log(TransferCoin(&rpcs.optimismTest, from, to, amount))
}

func TestTransferArbitrum(t *testing.T) {
	from, _ := NewAccountWithMnemonic(accountCase1.mnemonic)
	to := "0x8De5fF2edeD4d897da535ab0F379Ec1b9257eBaB"
	amount := "10000000000000" // 0.00001
	t.Log(TransferCoin(&rpcs.kccTest, from, to, amount))
}

func TestTransferAvax(t *testing.T) {
	from, _ := NewAccountWithMnemonic(testcase.M1)
	to := "0x2352D20fC81225c8ECD8f6FaA1B37F24FEd450c9"
	amount := "1000000000000000" // 0.001
	t.Log(TransferCoin(&rpcs.avaxTest, from, to, amount))
}

func TestTransferPolygon(t *testing.T) {
	from, _ := NewAccountWithMnemonic(testcase.M1)
	to := accountCase1.address
	amount := "100000000000000" // 0.0001
	t.Log(TransferCoin(&rpcs.polygonTest, from, to, amount))
}

func TransferCoin(at *rpcInfo, from *Account, toAddress string, amount string) (string, error) {
	chain := at.Chain()

	gasPrice, err := chain.SuggestGasPrice()
	if err != nil {
		return "", err
	}

	token := chain.MainEthToken()
	gasLimit, err := token.EstimateGasLimit(from.Address(), toAddress, gasPrice.Value, amount)
	if err != nil {
		return "", err
	}

	transaction := NewTransaction("", gasPrice.Value, gasLimit, toAddress, amount, "")
	signedTx, err := token.BuildTransferTxWithAccount(from, transaction)
	if err != nil {
		return "", err
	}

	return chain.SendRawTransaction(signedTx.Value)
}
