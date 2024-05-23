package btc

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/stretchr/testify/require"
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
	chains := []string{
		ChainMainnet,
		ChainSignet,
		ChainTestnet,
	}
	for _, chainnet := range chains {
		t.Run(chainnet, func(t *testing.T) {
			chain, err := NewChainWithChainnet(chainnet)
			require.Nil(t, err)
			rate, err := chain.SuggestFeeRate()
			require.Nil(t, err)
			t.Logf(`
	===== %v =====
	high: %v
	average: %v
	low: %v
			`, chainnet, rate.High, rate.Average, rate.Low)
		})
	}
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

	chain, err := NewChainWithChainnet(ChainSignet)
	require.NoError(t, err)
	hashString, err := chain.SendRawTransaction(txHex)
	require.NoError(t, err)
	t.Log("send raw transaction success: ", hashString)
}

func TestBtcClient(t *testing.T) {
	c, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         "k8s-chainpro-bitcoint-b4cc02de2a-782f891c88532615.elb.us-east-2.amazonaws.com:18332",
		User:         "auth",
		Pass:         "bitcoin-b2dd077",
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	require.NoError(t, err)
	info, err := c.GetBlockChainInfo()
	require.NoError(t, err)
	t.Log(info)
	hash := chainhash.Hash{}
	err = chainhash.Decode(&hash, "86398781cc7cdc5e736b88984d43eb41e70179dd8f7c5a2639012f40af7704a9")
	require.NoError(t, err)
	transaction, err := c.GetRawTransaction(&hash)
	require.NoError(t, err)
	t.Log(transaction)
	require.NoError(t, err)
	b := bytes.NewBuffer([]byte{})
	err = transaction.MsgTx().Serialize(b)
	require.NoError(t, err)
	t.Logf("%x", b.Bytes())
	psbtPackge, err := hex.DecodeString("70736274ff01007d0100000001026daad8c0c609877bc93d71494444038b7e5f5b7e9b063208220ccba27b2abe0000000000fdffffff0284660000000000002251200f46df23ae8a7725afd1b3f68699ea605c881ec606f6d75f5866af3dee8400233d41190000000000160014d31f560e89a32a2252544ed8e2878c234122b748000000000001011f0dbf190000000000160014d31f560e89a32a2252544ed8e2878c234122b74801086c02483045022100849d91c702cbf19c63a5d2ca9eb1c6cb0175c1e215fecdcc6136794a88cf52f5022078cb31b0f2170c21b4f85e3a0453949d718885e9249ab4b789d350e194b5c8d50121031bed49dd1a78ed1ff5cffef83d79d6936fc50624cc73d99da6c766297342b52f000000")
	require.NoError(t, err)
	psbtPack, err := psbt.NewFromRawBytes(bytes.NewReader(psbtPackge), false)
	require.NoError(t, err)
	extract, err := psbt.Extract(psbtPack)
	b1 := bytes.NewBuffer([]byte{})
	err = extract.Serialize(b1)
	require.NoError(t, err)
	require.Equal(t, b1.Bytes(), b.Bytes())
	rv1, err := DecodeTx("")
	require.NoError(t, err)
	rv2, err := DecodeTx("010000000001013622e96374e1ff2e1ef1a0ccb2d7237bba4ffb83ca556668d9989f91fd38b7c50100000000fdffffff022202000000000000160014d31f560e89a32a2252544ed8e2878c234122b748c6180000000000002251200f46df23ae8a7725afd1b3f68699ea605c881ec606f6d75f5866af3dee84002303402afe6303364ccd9dcd2d1e7d35f8652f27b09482deb4c670c16a4feef614c25fbbfc75bdad2e7ad9e08581a68405579f48e651f8bd052e4e7ebf8bed9566b42776208651980c1703275f50edf7e024e8ec5acffe21003250704dfe8cfa07333ea545ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800307b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2232323232222c22616d74223a317d6821c08651980c1703275f50edf7e024e8ec5acffe21003250704dfe8cfa07333ea54500000000")
	require.NoError(t, err)
	rv3, err := DecodeTx("0100000000010101bb7ac806b1a480cc60831d102e9fe51e68419ea6eb06826566d8c98e48670c0100000000fdffffff012202000000000000160014d31f560e89a32a2252544ed8e2878c234122b748034069bffa1de44309fcd657ce9ecfc61742b657cffaa07a477af17ff196a5fce68a16c35c8738d01c1cf5e64f082338bef54d3367f7cf93ad443604723ce7b88d7976208651980c1703275f50edf7e024e8ec5acffe21003250704dfe8cfa07333ea545ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800307b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2232323232222c22616d74223a317d6821c08651980c1703275f50edf7e024e8ec5acffe21003250704dfe8cfa07333ea54500000000")
	require.NoError(t, err)
	t.Log(rv1.TxHash().String())
	t.Logf("%x", rv1.TxIn[0].Witness[1])
	_, addrs, _, err := txscript.ExtractPkScriptAddrs(rv1.TxOut[0].PkScript, &chaincfg.SigNetParams)
	require.NoError(t, err)
	//txscript.ExtractWitnessProgramInfo()
	t.Log(addrs)
	//_, info, err := txscript.ExtractWitnessProgramInfo(rv1.TxIn[0].Witness[1])
	//require.NoError(t, err)
	//t.Log(info)
	t.Log(rv2.TxHash().String())
	t.Log(rv3.TxHash().String())
}
