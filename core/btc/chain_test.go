package btc

import (
	"strings"
	"testing"

	"github.com/coming-chat/wallet-SDK/core/base"
)

func TestChain_BalanceOfAddress(t *testing.T) {
	tests := []struct {
		name     string
		chainnet string // We need not test invalid chainnet
		address  string
		wantErr  bool
	}{
		{
			name:     "mainnet normal",
			chainnet: ChainMainnet,
			address:  accountCase.addrMainnet,
		},
		{
			name:     "signet normal",
			chainnet: ChainSignet,
			address:  accountCase.addrSignet,
		},
		{
			name:     "signet have balance",
			chainnet: ChainSignet,
			address:  "tb1pqtguh4mt0206qr7t3pze5zf4st4v3xtvqfhgv7q7j6ymnv7gtutqy4nrud",
		},
		{
			name:     "signet multi sign wallet address",
			chainnet: ChainSignet,
			address:  "tb1pesh6vwvq8xqfff9frs47sfyd83m6h8297ngn7qhg35xzsuuac67q8tgepm",
		},
		{
			name:     "mainnet error address",
			chainnet: ChainMainnet,
			address:  accountCase.addrMainnet + "s",
			wantErr:  true,
		},
		{
			name:     "signet chain to query miannet address",
			chainnet: ChainSignet,
			address:  accountCase.addrMainnet,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, _ := NewChainWithChainnet(tt.chainnet)
			got, err := chain.BalanceOfAddress(tt.address)

			host, _ := scanHostOf(tt.chainnet)
			url := host + "/address/" + tt.address
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceOfAddress() error = %v, wantErr %v, url = %v", err, tt.wantErr, url)
				return
			}
			if err == nil {
				t.Log("result: ", got.Total, ", Maybe you should verify via the link: ", url)
			}
		})
	}
}

func TestChain_BalanceOfPublicKey(t *testing.T) {
	tests := []struct {
		name      string
		chainnet  string // We need not test invalid chainnet
		publicKey string
		wantErr   bool
	}{
		{
			name:      "mainnet normal",
			chainnet:  ChainMainnet,
			publicKey: accountCase.publicKey,
		},
		{
			name:      "signet normal",
			chainnet:  ChainSignet,
			publicKey: accountCase.publicKey,
		},
		{
			name:      "mainnet error public key",
			chainnet:  ChainMainnet,
			publicKey: accountCase.publicKey + "s",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, _ := NewChainWithChainnet(tt.chainnet)
			got, err := chain.BalanceOfPublicKey(tt.publicKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("BalanceOfPublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				host, _ := scanHostOf(tt.chainnet)
				url := host + "/pubkey/" + strings.TrimPrefix(tt.publicKey, "0x")
				t.Log("result: ", got.Total, ", Maybe you should verify via the link: ", url)
			}
		})
	}
}

func TestChain_FetchTransactionDetail(t *testing.T) {
	tests := []struct {
		name     string
		chainnet string // We need not test invalid chainnet
		hash     string
		want     *base.TransactionDetail
		wantErr  bool
	}{
		{
			name:     "mainnet normal",
			chainnet: ChainMainnet,
			hash:     "182218b286c78aae63aac2f72fe44f7f35206500cb0bdb96eda20449c482b698",
			want:     &base.TransactionDetail{Status: 2, FinishTimestamp: 1649234531},
		},
		{
			name:     "signet normal",
			chainnet: ChainSignet,
			hash:     "efb7849f8f5a76da41faaa100977d189b025f1d01dee0fade87ffca4515af23a",
			want:     &base.TransactionDetail{Status: 2, FinishTimestamp: 1649836154},
		},
		{
			name:     "signet error hash",
			chainnet: ChainSignet,
			hash:     "efb7849f8f5a76da41faaa100977d189b025f1d01dee0fade87ffca4515af23",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain, _ := NewChainWithChainnet(tt.chainnet)
			got, err := chain.FetchTransactionDetail(tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchTransactionDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (got.Status != tt.want.Status || got.FinishTimestamp != tt.want.FinishTimestamp) {
				t.Errorf("FetchTransactionDetail() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChain_SuggestFeeRate(t *testing.T) {
	t.Log(SuggestFeeRate())
}

func TestChain_BatchFetchTransactionStatus(t *testing.T) {
	hashStrings := "xxx,182218b286c78aae63aac2f72fe44f7f35206500cb0bdb96eda20449c482b698,31244281753a3934060f6258cae6f87de7d96d8fc3c2f42d128dd3e0f72679b9"
	statuses := SdkBatchTransactionStatus(hashStrings, ChainMainnet)
	want := "0,2,2"
	if statuses != want {
		t.Errorf("BatchFetchTransactionStatus() %v, want %v", statuses, want)
	}
}

func TestSendRawTransaction(t *testing.T) {
	txHex := "02000000000101a1cea694898bbc9b288b61b9fc6f90d23433b8d58c3aa4fd0e99acc5f75b6a360100000000000000000158020000000000002251209435c165dfe71fd8d490176ff8085ed48862cddcb67de407e5b90fbe0d51cae5014022585aa4798a50ce3e156350fb93e7c476af3dac9bb9cd7dccea2721475c02483a397ab89f587d3cf37404568a6049b4b395c4b5afa39d14591a82c0d22173870000000002000000000101ddaae4634088b39ce763b7537798a0d0c502fb0e328ca0fe280919b70a2267a00000000000000000000920a1070000000000225120dc3e36639d0195d51b4a688315d0e718885266a31d6d2a132e5a70d1749e0f3720a10700000000002251202e2c3ecef16cec17551f485f5709aa51dc1af2702151619a664d07542a48ea4ec06e8f0000000000225120b5213fdd8cbe0ad7386098cb18d7c69540c0903352cb2d1fbecaeab33b38a5f7e0930400000000002251202de4b45e7b2042e6269fff60847d37b7c95b3ff1a50f8b0957d78dfa5c5854a360070b01000000002251207425f77de38e9b5fce71a7cc280a6883b1bf2bdf8dc66bdd8a3a673e43774463801a060000000000225120154b226b634f287d5b425b81376e433a71022fc5342c167aa3793ff898ce3bb860f590000000000022512063a5728682d8794f0b62f48c2bd103dc8e7eb86cd0df017f5c7584af0332a38ce09304000000000022512022557dec20b323f764edc523ae928dc6b6f610e62594ff051011ad83ff3ffe14dae1dd5e00000000225120cb72accaf99f5243ba473b14bfeb7372f204315ca0b9dcd7475565fc5a639d750340f1528a6b790ed9c9fb371823be025e61947233f6e155b4c4cfdefb4f9a199493a40fd19fbb9ca32b89baede1814a724c03689de66130db5af413e485c3cfd9b12220e46227fd1e926d6437986783b75526326192bc43f8e5fdb4e55d26024eaac1d1acfd0101c01b62c6adc3fdb99fce4e7771e3fd93dbefa66b4413e301fcb0c7e1db69f99f1169120ea48ec28d79dbf53cab6473466b9c8e4bd6436f71fc9f173dd2ee45e2aba2abec0c4fbebfd98abec73d60ca2e7a29e37cfd7ef1021231d8931a4e770ba97961b74f3ca918f63ac5cf26b2f1b4403815caaf102f45d0344f15d787eafb151804b5adb1d5c8fd2476055ff9efdfd747cfaa8d90a9333fc4f859966dada281dcc296d90bd725bf24a147d2c4cef9faa68b9e40828577b3590a4bc061b37a15da1322f0c173111a2e73c2e13a11445a70338ced5e8fb443cb44cab758c04d2316eafb1ab7a945c46a6c9bc1b6019cf54153d9878b84662280e2f3fd8ce80cf500000000"

	hashString, err := SendRawTransaction(txHex, "signet")
	if err != nil {
		t.Fatal("send raw transaction error: ", err)
	}
	t.Log("send raw transaction success: ", hashString)
}
