package wallet

import (
	"encoding/json"
	gsrc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/chainxTypes"
	"github.com/coming-chat/wallet-SDK/customscale"
	"testing"
)

var (
	apiMiniX, _   = gsrc.NewSubstrateAPI("wss://minichain-testnet.coming.chat")
	apiSherpax, _ = gsrc.NewSubstrateAPI("wss://sherpax-testnet.chainx.org")
	apiChainX, _  = gsrc.NewSubstrateAPI("wss://mainnet.chainx.org/ws")
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
	tx, err := txMetadata.NewBalanceTransferTx(address44, "10000000000000000000000")
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
	tx, err := txMetadata.NewBalanceTransferTx(address44, "1000000000000000000")
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
	tx, err := txMetadata.NewChainXBalanceTransferTx(address44, "100000000")
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
	tx, err := txMetadata.NewChainXBalanceTransferTx(address44, "100000000")
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
	tx, err := txMetadata.NewXAssetsTransferTx(address44, "10000000")
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
	tx, err := txMetadata.NewBalanceTransferTx(address44, "100000000")
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
	tx, err := txMetadata.NewBalanceTransferTx(address44, "100000000")
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
	tx, err := txMetadata.NewChainXBalanceTransferTx(address44, "10000000")
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
	tx, err := txMetadata.NewXAssetsTransferTx(address44, "100000000")
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
		"1000", 629560)
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
	var extSuccess chainxTypes.Extrinsic
	//var ext types.Extrinsic
	err := types.DecodeFromHexString("0x0603ff58526e07a368f79e61e4e905c2b88f9ca11eaee26e438ca618d98bb751049e1b0700e8764817", &extSuccess.Method)
	if err != nil {
		t.Fatal(err)
	}
	txMetadata, err := NewTx(ChainX)
	if err != nil {
		t.Fatal(err)
	}
	call, err := customscale.DecodeCall(txMetadata.metadata, &extSuccess.Method)
	if err != nil {
		t.Fatal(err)
	}
	jsonData, err := json.Marshal(call)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonData))
}

func TestDecodeTx2(t *testing.T) {
	var extSuccess chainxTypes.Extrinsic
	var errExtrinsic types.Extrinsic
	//var ext types.Extrinsic
	err := types.DecodeFromHexString("0x31058400043790c6e0b1cd20403f321c0532b5ca254d74eadcf3bcdb962f67c7e77caf42019e9efbbc9a91e9053ac2e0722db633bc2b48ea1f9bb21b59b1f54c8f012b8c565f8c6b49d6d00a5113422122278e95d904391cd22906a59d61780bf084fe4c87009501001e02180a00004416e63c59125cb8e946ba248e52bb6f2c06ff0aff71b9e9ba7532d2123ddc660284d7170a00004416e63c59125cb8e946ba248e52bb6f2c06ff0aff71b9e9ba7532d2123ddc66025a62020a00004416e63c59125cb8e946ba248e52bb6f2c06ff0aff71b9e9ba7532d2123ddc66d1070a00004416e63c59125cb8e946ba248e52bb6f2c06ff0aff71b9e9ba7532d2123ddc66d1070a00004416e63c59125cb8e946ba248e52bb6f2c06ff0aff71b9e9ba7532d2123ddc66d1070a00004416e63c59125cb8e946ba248e52bb6f2c06ff0aff71b9e9ba7532d2123ddc66d107", &errExtrinsic)
	//errExtrinsic.Signature = chainxTypes.ExtrinsicSignatureV4{}
	hexString, err := types.EncodeToHexString(errExtrinsic)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hexString)
	err = types.DecodeFromHexString("0x350284ffc8e8d0473afbe516cb772d504ecb091a139076c9aa4d3e0514aca7837599f861013c7f9eda8ee410ffe4436982303c1b603de6b910f1cd96b8cf143ca11aba7672b12792f59696ee7d35ea759c3bcd7e496c1b79907980696239ab03f8dde544810010000600ff30f9c1fd8d945474d39ef00547d7f13044dcc05b6a4db7ca8aee0a00622578500284d717", &extSuccess)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetadata(t *testing.T) {
	_, err := NewTx("")
	if err != nil {
		t.Fatal(err)
	}
}
