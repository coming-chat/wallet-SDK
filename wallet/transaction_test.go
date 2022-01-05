package wallet

import (
	gsrc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"math/big"
	"testing"
	"wallet-SDK/chainxTypes"
)

var (
	apiMiniX, _   = gsrc.NewSubstrateAPI("wss://minichain-testnet.coming.chat")
	apiChainX, _  = gsrc.NewSubstrateAPI("wss://testnet.chainx.org")
	apiSherpax, _ = gsrc.NewSubstrateAPI("wss://sherpax-testnet.chainx.org")
	Minix         = ""
	ChainX        = ""
	Sherpax       = ""
	_             = client.CallWithBlockHash(apiSherpax.Client, &Sherpax, "state_getMetadata", nil)
	_             = client.CallWithBlockHash(apiMiniX.Client, &Minix, "state_getMetadata", nil)
	_             = client.CallWithBlockHash(apiChainX.Client, &ChainX, "state_getMetadata", nil)
)

func TestTransactionSherpax(t *testing.T) {
	txMetadata, err := NewTx(Sherpax)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewBalanceTransferTx(address44, 1000000000000000000)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0xbcffcb56cf05eb71e5f59eaf35de2bbe330f925a065d852859b1737ce02342a0", 1, 12, 1)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, "")
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext chainxTypes.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Sherpax sendTx: %v", sendTx)
}

func TestTransactionSherpaxGetUnSign(t *testing.T) {
	txMetadata, err := NewTx(Sherpax)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewBalanceTransferTx(address44, 1000000000000000000)
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetUnSignTx()
	if err != nil {
		t.Fatal(err)
	}
	var ext chainxTypes.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Sherpax unSign sendTx: %v", sendTx)
}

func TestTransactionPCX(t *testing.T) {
	txMetadata, err := NewTx(ChainX)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewChainXBalanceTransferTx(address44, 100000000)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0x2fd9e861564c428cf16c3d6e0ec82c5a07ddcd9ec44f37ff4627ab385d1cb597", 1, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, "")
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext chainxTypes.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Chainx sendTx: %v", sendTx)
}

func TestTransactionPCXByKeystore(t *testing.T) {
	txMetadata, err := NewTx(ChainX)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewChainXBalanceTransferTx(address44, 100000000)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0x2fd9e861564c428cf16c3d6e0ec82c5a07ddcd9ec44f37ff4627ab385d1cb597", 1, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWalletFromKeyStore(keystore1, password1)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, password1)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext chainxTypes.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Chainx sendTx: %v", sendTx)
}

func TestTransactionXBTCByKeystore(t *testing.T) {
	txMetadata, err := NewTx(ChainX)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewXAssetsTransferTx(address44, 10000000)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWalletFromKeyStore(keystore1, password1)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0x2fd9e861564c428cf16c3d6e0ec82c5a07ddcd9ec44f37ff4627ab385d1cb597", 1, 2, 1)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, password1)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext chainxTypes.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Chainx sendTx: %v", sendTx)
}

func TestTransactionMini(t *testing.T) {
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewBalanceTransferTx(address44, 100000000)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 2, 115, 1)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, "")
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext types.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Mini sendTx: %v", sendTx)
}

func TestTransactionMiniByKeystore(t *testing.T) {
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewBalanceTransferTx(address44, 100000000)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWalletFromKeyStore(keystore1, password1)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 0, 115, 1)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, password1)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext types.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Mini sendTx: %v", sendTx)
}

func TestTransactionNFTByKeystore(t *testing.T) {
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewComingNftTransferTx("5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1", 289262366)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWalletFromKeyStore(keystore1, password1)
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 1, 115, 1)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, password1)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext types.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Mini sendTx: %v", sendTx)
}

func TestGetUnSignTxNFT(t *testing.T) {
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewComingNftTransferTx("5PjZ58jF72pCz6Y3FkB3jtyWbhhEbWxBz8CkDD7NG3yjL6s1", 289262366)
	if err != nil {
		t.Fatal(err)
	}
	signTx, err := tx.GetUnSignTx()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("UnSign tx: %s", signTx)
}

func TestGetUnSignTxPCX(t *testing.T) {
	txMetadata, err := NewTx(ChainX)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewChainXBalanceTransferTx(address44, 10000000)
	if err != nil {
		t.Fatal(err)
	}
	signTx, err := tx.GetUnSignTx()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("UnSign tx: %s", signTx)
}

func TestGetUnSignTxXBTC(t *testing.T) {
	txMetadata, err := NewTx(ChainX)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewXAssetsTransferTx(address44, 100000000)
	if err != nil {
		t.Fatal(err)
	}
	signTx, err := tx.GetUnSignTx()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("UnSign tx: %s", signTx)
}

func TestThreshold(t *testing.T) {
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := txMetadata.NewThreshold("0x24c5d4ad9a2052bf8535f98b46815fab02ace5eb286951459686229319c49556", "5QUEsBTRCB5GNGVHo67DrDpNY5y9Y12RpjNpSzrK56fGeS5H", "f08e6ce7b72b2fb256b1bf1e9186920a8b10d251c38bec9ae167f4964aeefe01b4d77d08f9006900c924756dfb04472ddf21b121d8f6e8f92932649cbb4f6582", "aa68dced52cfe04e3b7a0457bdcfda00e463044eadac12bada22c192e4f6af5d", "44a39dcf13ec8b9427375f3cd6c3552f5941b633092f7bfaee5bc6d8d8b0d03a898d4079480f8326122d60ac1b8747514a8ae6adeaea8dbb758597a2834e27f6c39da7eb29fbc714d0190bf5a29be2da0523ba8f726a8b1c213173d9d568e626", "576520617265206c6567696f6e21", "6ed03482e88c37d015cc44b7fc581209c37caf0a74fc25479ef3a4630eb34b58",
		1000, 629560)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWallet("boss mind sauce seek clutch busy boil screen room timber shop same")
	if err != nil {
		t.Fatal(err)
	}
	signData, err := tx.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 0, 115, 1)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := wallet.Sign(signData, "")
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := tx.GetTx(publicKey, signed)
	if err != nil {
		t.Fatal(err)
	}
	var ext types.Extrinsic
	err = types.DecodeFromHexString(sendTx, &ext)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Mini Threshold sendTx: %v", sendTx)
}

func TestNewTransactionFromHex(t *testing.T) {
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}
	transaction, err := txMetadata.NewTransactionFromHex(false, "0x6c042a013fb26e2800000000000e2707000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}
	signData, err := transaction.GetSignData("0xfb58f83706a065ced8f658fafaba97e6e49b772287e332077c499784184eda9f", 0, 115, 1)
	if err != nil {
		t.Fatal(err)
	}
	wallet, err := NewWallet("boss mind sauce seek clutch busy boil screen room timber shop same")
	if err != nil {
		t.Fatal(err)
	}
	signedData, err := wallet.Sign(signData, "")
	if err != nil {
		t.Fatal(err)
	}
	pubkey, err := wallet.GetPublicKey()
	if err != nil {
		t.Fatal(err)
	}
	sendTx, err := transaction.GetTx(pubkey, signedData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sendTx)
}

func TestDecodeTx(t *testing.T) {
	var extSuccess types.Extrinsic
	//var ext types.Extrinsic
	err := types.DecodeFromHexString("0x6c042a013fb26e2800000000000e2707000000000000000000000000", &extSuccess)
	if err != nil {
		t.Fatal(err)
	}
	txMetadata, err := NewTx(Minix)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := txMetadata.newTx(false, "ComingAuction.bid", types.NewU64(123456789), types.NewU128(*big.NewInt(9999994000)))
	if err != nil {
		return
	}

	t.Log(tx.extrinsic)
	t.Log(tx.GetUnSignTx())
	t.Log(types.EncodeToHexString(tx.extrinsic.Method))

	var (
		variant types.Si1Variant
		mod     types.PalletMetadataV14
	)

	for _, mod = range txMetadata.metadata.AsMetadataV14.Pallets {
		if mod.Index == types.NewU8(extSuccess.Method.CallIndex.SectionIndex) {
			callType := mod.Calls.Type.Int64()
			if typ, ok := txMetadata.metadata.AsMetadataV14.EfficientLookup[callType]; ok {
				if len(typ.Def.Variant.Variants) > 0 {
					for _, variant = range typ.Def.Variant.Variants {
						if variant.Index == types.NewU8(extSuccess.Method.CallIndex.MethodIndex) {
							break
						}
					}
				}
			}
		}
	}
	t.Logf("%s.%s", mod.Name, variant.Name)
	var arg []interface{}
	for _, field := range variant.Fields {
		if field.HasName {
			t.Log(field.Name)
		}
		if field.HasTypeName {
			t.Log(field.TypeName)
		}
		arg = append(arg, field.Type)
	}
	err = types.DecodeFromBytes(extSuccess.Method.Args, &arg)
	if err != nil {
		t.Fatal(err)
	}
	//err = types.DecodeFromHexString("0xed068400c42129c6bed8c7fc85776a7687b250908fde45a3b25365ce5d3432d36d29d459018a677a90cb4c77ea2d2114203763b710c7b7192957d1b4eff43df754b6c24a043a9b6cbd3a2a9727619aeff13055d64bbbe5567d69ac9ddbf474c666c9de2b890000001e00081f0024c5d4ad9a2052bf8535f98b46815fab02ace5eb286951459686229319c495560101f08e6ce7b72b2fb256b1bf1e9186920a8b10d251c38bec9ae167f4964aeefe01b4d77d08f9006900c924756dfb04472ddf21b121d8f6e8f92932649cbb4f658280aa68dced52cfe04e3b7a0457bdcfda00e463044eadac12bada22c192e4f6af5d810144a39dcf13ec8b9427375f3cd6c3552f5941b633092f7bfaee5bc6d8d8b0d03a898d4079480f8326122d60ac1b8747514a8ae6adeaea8dbb758597a2834e27f6c39da7eb29fbc714d0190bf5a29be2da0523ba8f726a8b1c213173d9d568e62638576520617265206c6567696f6e21806ed03482e88c37d015cc44b7fc581209c37caf0a74fc25479ef3a4630eb34b581f0124c5d4ad9a2052bf8535f98b46815fab02ace5eb286951459686229319c4955600e8030000000000000000000000000000389b0900209f0900", &ext)
	//
	//sh := extSuccess.Method.Args[248:280]
	//shS := hex.EncodeToString(sh)
	//t.Log(shS)
	////err := types.DecodeFromHexString("0x410284ffdc64bef918ddda3126a39a11113767741ddfdf91399f055e1d963f2ae1ec2535018a3c93d07de30fe4e4fa58af5cd844d83b01f12bf34aa5ca3057db3b2d65e0791f26e4f3d73c2287279a4a57ef19273e1c333b20924240ac218359a670de36870000000600ffdc64bef918ddda3126a39a11113767741ddfdf91399f055e1d963f2ae1ec25350b00a0724e1809", &ext)
	////err := types.DecodeFromHexString("0x350284004826ef96bd5e88fc3e93af37dfea53635e84fed316ccebc7d1bd2f4f259af949018ad516499042d3b8fbada690b8dc78af645770688a6b56b3624eaf19bf10050ad07ee9884fe569980bc13f125711fe1a283a27cd9f222b456470440c9a9809850000000a00004826ef96bd5e88fc3e93af37dfea53635e84fed316ccebc7d1bd2f4f259af9490208af2f", &ext)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(ext)
}

func TestDecodeTx2(t *testing.T) {
	var extSuccess chainxTypes.Extrinsic
	var errExtrinsic chainxTypes.Extrinsic
	//var ext types.Extrinsic
	err := types.DecodeFromHexString("0x350284ffc8e8d0473afbe516cb772d504ecb091a139076c9aa4d3e0514aca7837599f861010600ff30f9c1fd8d945474d39ef00547d7f13044dcc05b6a4db7ca8aee0a0062257850025a620200100002000000010000002fd9e861564c428cf16c3d6e0ec80010000600ff30f9c1fd8d945474d39ef00547d7f13044dcc05b6a4db7ca8aee0a0062257850025a6202", &errExtrinsic)
	err = types.DecodeFromHexString("0x350284ffc8e8d0473afbe516cb772d504ecb091a139076c9aa4d3e0514aca7837599f861013c7f9eda8ee410ffe4436982303c1b603de6b910f1cd96b8cf143ca11aba7672b12792f59696ee7d35ea759c3bcd7e496c1b79907980696239ab03f8dde544810010000600ff30f9c1fd8d945474d39ef00547d7f13044dcc05b6a4db7ca8aee0a00622578500284d717", &extSuccess)
	pub1, err := types.EncodeToHexString(extSuccess.Signature.Signer.AsAccountID)
	pub2, err := types.EncodeToHexString(errExtrinsic.Signature.Signer.AsAccountID)
	t.Log(pub1)
	t.Log(pub2)
	if err != nil {
		t.Fatal(err)
	}
}
