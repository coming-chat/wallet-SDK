package polka

import (
	"testing"

	gsrc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

var (
	// apiMiniX, _   = gsrc.NewSubstrateAPI("wss://minichain-testnet.coming.chat")
	// apiSherpax, _ = gsrc.NewSubstrateAPI("wss://sherpax-testnet.chainx.org")
	// apiChainX, _  = gsrc.NewSubstrateAPI("wss://testnet3.chainx.org")
	apiMiniX   gsrc.SubstrateAPI
	apiSherpax gsrc.SubstrateAPI
	apiChainX  gsrc.SubstrateAPI
	Minix      = ""
	ChainX     = ""
	Sherpax    = ""
	// _          = client.CallWithBlockHash(apiSherpax.Client, &Sherpax, "state_getMetadata", nil)
	// _          = client.CallWithBlockHash(apiMiniX.Client, &Minix, "state_getMetadata", nil)
	// _          = client.CallWithBlockHash(apiChainX.Client, &ChainX, "state_getMetadata", nil)
)

const (
	keystore1 = "{\"address\":\"5Gc8bR5p9JeCY3dpCvdonRWn79UxhKycDb8aC7xfqQPqWhr8\",\"encoded\":\"jC9MOH7OPYbHdJtiOWFW0lpMUCFO4nASKjzqHvXpEiYAgAAAAQAAAAgAAACm2Dm/CZ98R1uy34lMj7tr9+i3ERCFoeCSdNwOScsyDkvLwhVGv6qxOzmdiR7vzgRgEizMQbq17k0C1Tk59WyDnf9OfaGQTenQQpnFPiXxcmDa6TXQvF7Eq8VYw009ANLmDTIQ125JdQX6edYY85ZFpLiOltXiad44mhS1mC8OSCcOHsViVrk3Lk0eMsClYS1SUzv3QDCoHChFu6Za\",\"encoding\":{\"content\":[\"pkcs8\",\"sr25519\"],\"type\":[\"scrypt\",\"xsalsa20-poly1305\"],\"version\":\"3\"},\"meta\":{\"genesisHash\":\"0x3a10a25727b09cf04a9d143c3ebefb179c3c45613297339d3cbec4e5d4c75242\",\"name\":\"NFT测试2\",\"tags\":[],\"whenCreated\":1623900058655}}"
	password1 = "111"
)

func GenerateTransferSignData(on rpcInfo, from, to, amount string) ([]byte, error) {
	chain, err := on.Chain()
	if err != nil {
		return nil, err
	}
	metadata, err := chain.GetMetadataString()
	if err != nil {
		return nil, err
	}
	tx, err := NewTx(metadata)
	if err != nil {
		return nil, err
	}
	transaction, err := tx.NewBalanceTransferTx(to, amount)
	if err != nil {
		return nil, err
	}
	return chain.GetSignData(transaction, from)
}

func TestTransaction_GetSignData(t1 *testing.T) {
	type args struct {
		on     rpcInfo
		from   string
		to     string
		amount string
	}
	tests := []struct {
		name    string
		args    args
		wantHex string
		wantErr bool
	}{
		{
			name: "chainx-pre trannsfer 1 PCX",
			args: args{
				on:     rpcs.chainxTest,
				from:   accountCase.address44,
				to:     "5TE1T7Znw5eaDqzpxCm8MBocbZLbbeeZ4GcnsRN37s2amd5s",
				amount: "100000000",
			},
			wantHex: "0x0600009e9fefcd4c4d6702b7599f093a860de30a99b114072e3cb3aed14172c2433c240284d71700000015000000040000002e25d2145e9ecf2d1c185b052e085e3c39340edf3dba74f702653afcdd0a9c372e25d2145e9ecf2d1c185b052e085e3c39340edf3dba74f702653afcdd0a9c37",
		},
		{
			name: "chainx trannsfer from error address",
			args: args{
				on:     rpcs.chainxTest,
				from:   "5TE1T7Znw5eaDqzpxCm8MBocbZLbbeeZ4GcnsRN37s2amd55",
				to:     accountCase.address44,
				amount: "100000000",
			},
			wantErr: true,
		},
		{
			name: "chainx trannsfer to error address",
			args: args{
				on:     rpcs.chainxTest,
				from:   accountCase.address44,
				to:     "5TE1T7Znw5eaDqzpxCm8MBocbZLbbeeZ4GcnsRN37s2amd5S",
				amount: "100000000",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			signData, err := GenerateTransferSignData(tt.args.on, tt.args.from, tt.args.to, tt.args.amount)
			if err != nil {
				if !tt.wantErr {
					t1.Errorf("GetSignData() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			signedHex := types.HexEncodeToString(signData)
			if signedHex != tt.wantHex {
				t1.Errorf("GetSignData() got = %v, want %v, Maybe the metadata version updated ???", signedHex, tt.wantHex)
			}
		})
	}
}

// func TestTransactionSherpax(t *testing.T) {
// 	txMetadata, err := NewTx(Sherpax)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "10000000000000000000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0xbcffcb56cf05eb71e5f59eaf35de2bbe330f925a065d852859b1737ce02342a0", 1, 12, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWallet(testSecretPhrase)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Sherpax sendTx: %v", sendTx)
// }

// func TestTransactionSherpaxGetUnSign(t *testing.T) {
// 	txMetadata, err := NewTx(Sherpax)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "1000000000000000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetUnSignTx()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Sherpax unSign sendTx: %v", sendTx)
// }

// func TestTransactionPCX(t *testing.T) {
// 	txMetadata, err := NewTx(ChainX)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "100000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0x2e25d2145e9ecf2d1c185b052e085e3c39340edf3dba74f702653afcdd0a9c37", 2, 13, 4)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWallet(testSecretPhrase)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Chainx sendTx: %v", sendTx)
// }

// func TestTransactionPCXByKeystore(t *testing.T) {
// 	txMetadata, err := NewTx(ChainX)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "100000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0x2fd9e861564c428cf16c3d6e0ec82c5a07ddcd9ec44f37ff4627ab385d1cb597", 1, 2, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWalletFromKeyStore(keystore1, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Chainx sendTx: %v", sendTx)
// }

// func TestTransactionXBTCByKeystore(t *testing.T) {
// 	txMetadata, err := NewTx(ChainX)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewXAssetsTransferTx(address44, "10000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWalletFromKeyStore(keystore1, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0x2fd9e861564c428cf16c3d6e0ec82c5a07ddcd9ec44f37ff4627ab385d1cb597", 1, 2, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Chainx sendTx: %v", sendTx)
// }

// func TestTransactionMini(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "100000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWallet(testSecretPhrase)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 2, 115, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var ext types.Extrinsic
// 	err = types.DecodeFromHexString(sendTx, &ext)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Mini sendTx: %v", sendTx)
// }

// func TestTransactionMiniByKeystore(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "100000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWalletFromKeyStore(keystore1, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 0, 115, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var ext types.Extrinsic
// 	err = types.DecodeFromHexString(sendTx, &ext)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Mini sendTx: %v", sendTx)
// }

// func TestTransactionNFTByKeystore(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewComingNftTransferTx("5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1", 289262366)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWalletFromKeyStore(keystore1, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 1, 115, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, password1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var ext types.Extrinsic
// 	err = types.DecodeFromHexString(sendTx, &ext)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Mini sendTx: %v", sendTx)
// }

// func TestGetUnSignTxMINI(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx("5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1", "10000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signTx, err := tx.GetUnSignTx()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("UnSign tx: %s", signTx)
// }

// func TestGetUnSignTxNFT(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewComingNftTransferTx("5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1", 289262366)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signTx, err := tx.GetUnSignTx()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("UnSign tx: %s", signTx)
// }

// func TestGetUnSignTxPCX(t *testing.T) {
// 	txMetadata, err := NewTx(ChainX)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewBalanceTransferTx(address44, "10000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signTx, err := tx.GetUnSignTx()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("UnSign tx: %s", signTx)
// }

// func TestGetUnSignTxXBTC(t *testing.T) {
// 	txMetadata, err := NewTx(ChainX)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewXAssetsTransferTx(address44, "100000000")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signTx, err := tx.GetUnSignTx()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("UnSign tx: %s", signTx)
// }

// func TestThreshold(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	tx, err := txMetadata.NewThreshold("0x24c5d4ad9a2052bf8535f98b46815fab02ace5eb286951459686229319c49556", "5QUEsBTRCB5GNGVHo67DrDpNY5y9Y12RpjNpSzrK56fGeS5H", "f08e6ce7b72b2fb256b1bf1e9186920a8b10d251c38bec9ae167f4964aeefe01b4d77d08f9006900c924756dfb04472ddf21b121d8f6e8f92932649cbb4f6582", "aa68dced52cfe04e3b7a0457bdcfda00e463044eadac12bada22c192e4f6af5d", "44a39dcf13ec8b9427375f3cd6c3552f5941b633092f7bfaee5bc6d8d8b0d03a898d4079480f8326122d60ac1b8747514a8ae6adeaea8dbb758597a2834e27f6c39da7eb29fbc714d0190bf5a29be2da0523ba8f726a8b1c213173d9d568e626", "576520617265206c6567696f6e21", "6ed03482e88c37d015cc44b7fc581209c37caf0a74fc25479ef3a4630eb34b58",
// 		"1000", 629560)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWallet("boss mind sauce seek clutch busy boil screen room timber shop same")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 0, 115, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signed, err := wallet.Sign(signData, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	publicKey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := tx.GetTx(publicKey, signed)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	var ext types.Extrinsic
// 	err = types.DecodeFromHexString(sendTx, &ext)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("Mini Threshold sendTx: %v", sendTx)
// }

// func TestNewTransactionFromHex(t *testing.T) {
// 	txMetadata, err := NewTx(Minix)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	transaction, err := txMetadata.NewTransactionFromHex("0xb00429021ecb3d110000000000043790c6e0b1cd20403f321c0532b5ca254d74eadcf3bcdb962f67c7e77caf42")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signData, err := transaction.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 0, 115, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	wallet, err := NewWallet("boss mind sauce seek clutch busy boil screen room timber shop same")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signedData, err := wallet.Sign(signData, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	pubkey, err := wallet.GetPublicKey()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	sendTx, err := transaction.GetTx(pubkey, signedData)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(sendTx)
// }

// func TestDecodeTx(t *testing.T) {
// 	var ext types.Extrinsic
// 	err := types.DecodeFromHexString("0x0603ff58526e07a368f79e61e4e905c2b88f9ca11eaee26e438ca618d98bb751049e1b0700e8764817", &ext.Method)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	txMetadata, err := NewTx(ChainX)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	call, err := customscale.DecodeCall(txMetadata.metadata, &ext.Method)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	jsonData, err := json.Marshal(call)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(string(jsonData))
// }

// func TestDecodeTx2(t *testing.T) {
// 	var errExtrinsic types.Extrinsic
// 	//var ext types.Extrinsic
// 	err := types.DecodeFromHexString("0x4d0284003e3543edb7c9d4c2246b841a14cfe1b9076043de3a219be5a76c97556297463f01d85fc2c1979e9ceaa639082545f041225400a75df50c4a54261913b30a84f83a1725f03e0dc2a425751e98b27f19a32e2716f92f067ab8908fc947dc7283ce88002503001e02040a0000760bc02040cd949016216c067331ee0d333056773f63a3ad5b3d7365e2b3c32f07c0f0172f03", &errExtrinsic)
// 	//errExtrinsic.Signature = chainxTypes.ExtrinsicSignatureV4{}
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

// func TestMetadata(t *testing.T) {
// 	_, err := NewTx("")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestTx_NewExtrinsics(t1 *testing.T) {
	type fields struct {
		metadata *types.Metadata
	}
	type args struct {
		call string
		args []interface{}
	}
	chain, err := rpcs.chainxTest.Chain()
	if err != nil {
		t1.Error(err)
	}
	metadata, err := chain.GetMetadataString()
	if err != nil {
		t1.Error(err)
	}
	t, err := NewTx(metadata)
	if err != nil {
		t1.Error(err)
	}
	accountID, err := addressStringToMultiAddress("5Gc8bR5p9JeCY3dpCvdonRWn79UxhKycDb8aC7xfqQPqWhr8")
	if err != nil {
		t1.Error(err)
	}
	newCall, err := types.NewCall(t.metadata, "Balances.transfer", accountID, types.NewUCompactFromUInt(100000000))
	if err != nil {
		return
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Transaction
		wantErr bool
	}{
		{
			name: "transfer1",
			args: args{
				call: "Utility.batch_all",
				args: []interface{}{
					[]types.Call{newCall, newCall, newCall, newCall, newCall},
				},
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			got, err := t.NewExtrinsics(tt.args.call, tt.args.args...)
			data, err := chain.GetSignData(got, "5Gc8bR5p9JeCY3dpCvdonRWn79UxhKycDb8aC7xfqQPqWhr8")
			if err != nil {
				t1.Error(err)
			}
			t1.Log(ByteToHex(data))
			withKeystore, err := NewAccountWithKeystore(keystore1, password1, 44)
			if err != nil {
				t1.Error(err)
			}
			signedData, err := withKeystore.Sign(data, password1)
			if err != nil {
				t1.Error(err)
			}
			t1.Log(ByteToHex(signedData))
			tx, err := got.GetTx(withKeystore.PublicKey(), signedData)
			if err != nil {
				t1.Error(err)
			}
			t1.Log(tx)
		})
	}
}
