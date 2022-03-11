package wallet

import (
	"encoding/hex"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/client"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"math/big"
	"testing"
)

const (
	keystore1 = "{\"address\":\"5Gc8bR5p9JeCY3dpCvdonRWn79UxhKycDb8aC7xfqQPqWhr8\",\"encoded\":\"jC9MOH7OPYbHdJtiOWFW0lpMUCFO4nASKjzqHvXpEiYAgAAAAQAAAAgAAACm2Dm/CZ98R1uy34lMj7tr9+i3ERCFoeCSdNwOScsyDkvLwhVGv6qxOzmdiR7vzgRgEizMQbq17k0C1Tk59WyDnf9OfaGQTenQQpnFPiXxcmDa6TXQvF7Eq8VYw009ANLmDTIQ125JdQX6edYY85ZFpLiOltXiad44mhS1mC8OSCcOHsViVrk3Lk0eMsClYS1SUzv3QDCoHChFu6Za\",\"encoding\":{\"content\":[\"pkcs8\",\"sr25519\"],\"type\":[\"scrypt\",\"xsalsa20-poly1305\"],\"version\":\"3\"},\"meta\":{\"genesisHash\":\"0x3a10a25727b09cf04a9d143c3ebefb179c3c45613297339d3cbec4e5d4c75242\",\"name\":\"NFT测试2\",\"tags\":[],\"whenCreated\":1623900058655}}"
	password1 = "111"
)

func TestPolDecode(t *testing.T) {
	wallet, err := NewWalletFromKeyStore(keystore1, password1)
	if err != nil {
		t.Fatal(err)
	}
	msg := []byte("8asd8u8qw9ddqu9w8d9wqud89q9wd8uq89uw8u89r893h22")
	sign, err := wallet.Sign(msg, password1)
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hex.EncodeToString(sign))
}

func TestKeystoreSignAndSeedSign(t *testing.T) {
	var (
		metadata       types.Metadata
		latestMetadata string
		callList       []*types.Call
		nonce          int64
	)

	seedWallet, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Error(err)
	}
	keystoreWallet, err := NewWalletFromKeyStore(keystore2, password)
	if err != nil {
		t.Error(err)
	}

	miniXApi, err := gsrpc.NewSubstrateAPI("wss://minichain-testnet.coming.chat/ws")
	if err != nil {
		t.Errorf("connect to MiniX failed: %v", err)
	}
	err = client.CallWithBlockHash(miniXApi.Client, &latestMetadata, "state_getMetadata", nil)
	if err != nil {
		t.Error(err)
	}

	hotWalletAddress, err := seedWallet.GetAddress(44)
	if err != nil {
		t.Error(err)
	}

	err = client.CallWithBlockHash(miniXApi.Client, &nonce, "system_accountNextIndex", nil, hotWalletAddress)
	if err != nil {
		t.Error(err)
	}

	genesisHash, err := miniXApi.RPC.Chain.GetBlockHash(0)
	if err != nil {
		t.Error(err)
	}

	runtimeVersion, err := miniXApi.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		t.Error(err)
	}

	err = types.DecodeFromHexString(latestMetadata, &metadata)
	if err != nil {
		t.Error(err)
	}

	tx, err := NewTx(latestMetadata)
	if err != nil {
		t.Error(err)
	}

	type DbTransfer struct {
		Address string
		Amount  *big.Int
	}

	var dbTransferRecords = []*DbTransfer{{
		Address: address44,
		Amount:  big.NewInt(10000000),
	},
		{
			Address: address44,
			Amount:  big.NewInt(100000000),
		},
		{
			Address: address44,
			Amount:  big.NewInt(501),
		},
		{
			Address: address44,
			Amount:  big.NewInt(501),
		},
		{
			Address: address44,
			Amount:  big.NewInt(501),
		}}

	for _, iterm := range dbTransferRecords {
		publicKey, err := AddressToPublicKey(iterm.Address)
		if err != nil {
			t.Error(err)
		}
		accountID, err := types.NewMultiAddressFromHexAccountID(publicKey)
		if err != nil {
			t.Error(err)
		}
		call, err := types.NewCall(&metadata, "Balances.transfer", accountID, types.NewUCompact(iterm.Amount))
		if err != nil {
			t.Error(err)
		}
		callList = append(callList, &call)
	}

	extrinsicsTxSeed, err := tx.NewExtrinsics("Utility.batch_all", callList)
	if err != nil {
		t.Error(err)
	}

	signDataSeed, err := extrinsicsTxSeed.GetSignData(genesisHash.Hex(), nonce, int32(runtimeVersion.SpecVersion), int32(runtimeVersion.TransactionVersion))
	if err != nil {
		t.Error(err)
	}

	signatureDataSeed, err := seedWallet.Sign(signDataSeed, password)
	if err != nil {
		t.Error(err)
	}

	seedWalletPublicKey, err := seedWallet.GetPublicKey()
	if err != nil {
		t.Error(err)
	}

	sendTxSeed, err := extrinsicsTxSeed.GetTx(seedWalletPublicKey, signatureDataSeed)
	if err != nil {
		t.Error(err)
	}

	t.Log(sendTxSeed)

	extrinsicsTxKeystore, err := tx.NewExtrinsics("Utility.batch_all", callList)
	if err != nil {
		t.Error(err)
	}

	signDataKeystore, err := extrinsicsTxKeystore.GetSignData(genesisHash.Hex(), nonce, int32(runtimeVersion.SpecVersion), int32(runtimeVersion.TransactionVersion))
	if err != nil {
		t.Error(err)
	}

	t.Log(types.HexEncodeToString(signDataKeystore))

	signatureDataKeystore, err := keystoreWallet.Sign(signDataKeystore, password)
	if err != nil {
		t.Error(err)
	}

	t.Log(types.HexEncodeToString(signatureDataKeystore))

	keystoreWalletPublicKey, err := keystoreWallet.GetPublicKey()
	if err != nil {
		t.Error(err)
	}

	sendTxKeystore, err := extrinsicsTxSeed.GetTx(keystoreWalletPublicKey, signatureDataKeystore)
	if err != nil {
		t.Error(err)
	}
	t.Log(sendTxKeystore)
}
