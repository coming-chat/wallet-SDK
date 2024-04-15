package btc

import (
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
	"github.com/tyler-smith/go-bip39"
)

type TestAccountCase struct {
	mnemonic    string
	privateKey  string
	publicKey   string
	addrMainnet string
	addrSignet  string
}

var accountCase = &TestAccountCase{
	mnemonic:    "unaware oxygen allow method allow property predict various slice travel please priority",
	privateKey:  "0xc7fceb75bafba7aa10ffe10315352bfc523ac733f814e6a311bc736873df8923",
	publicKey:   "0x04a721f170043daafde0fa925ab6caf5d2abcdadd2249291b1840e3d99a3f41149e13185ef52451eef2e7cc0c5fe4180b64ca2d17eb886b2328518f6aed684719a",
	addrMainnet: "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg",
	addrSignet:  "tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728",
}
var errorCase = &TestAccountCase{
	mnemonic: "unaware oxygen allow method allow property predict various slice travel please check",
}

func TestAccount(t *testing.T) {
	mn := testcase.M1
	acc, err := NewAccountWithMnemonic(mn, ChainMainnet, AddressTypeComingTaproot)
	require.Nil(t, err)
	t.Log(acc.Address())
}

func TestNewAccountWithMnemonic(t *testing.T) {
	type args struct {
		mnemonic string
		chainnet string
	}
	tests := []struct {
		name        string
		args        args
		wantAddress string // If the generated address can match, there is no problem.
		wantErr     bool
	}{
		{
			name:        "mainnet nomal",
			args:        args{mnemonic: accountCase.mnemonic, chainnet: ChainMainnet},
			wantAddress: accountCase.addrMainnet,
		},
		{
			name:        "coming bitcoin nomal",
			args:        args{mnemonic: accountCase.mnemonic, chainnet: ChainBitcoin},
			wantAddress: accountCase.addrMainnet,
		},
		{
			name:        "signet nomal",
			args:        args{mnemonic: accountCase.mnemonic, chainnet: ChainSignet},
			wantAddress: accountCase.addrSignet,
		},
		{
			name:    "error chainnet",
			args:    args{mnemonic: accountCase.mnemonic, chainnet: "xxxxxxx"},
			wantErr: true,
		},
		{
			name:    "error mnemonic",
			args:    args{mnemonic: errorCase.mnemonic, chainnet: ChainSignet},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.args.mnemonic, tt.args.chainnet, AddressTypeComingTaproot)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.Address() != tt.wantAddress {
				t.Errorf("NewAccountWithMnemonic() got = %v, want %v", got.Address(), tt.wantAddress)
			}
		})
	}
}

func TestAccount_DeriveAccountAt(t *testing.T) {
	baseAccount, err := NewAccountWithMnemonic(accountCase.mnemonic, ChainMainnet, AddressTypeComingTaproot)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		chainnet    string
		wantAddress string
		wantErr     bool
	}{
		{name: "same as mainnet", chainnet: ChainMainnet, wantAddress: accountCase.addrMainnet},
		{name: "change signet", chainnet: ChainSignet, wantAddress: accountCase.addrSignet},
		{name: "error net", chainnet: "signet2", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := baseAccount.DeriveAccountAt(tt.chainnet)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeriveAccountAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && got.Address() != tt.wantAddress {
				t.Errorf("DeriveAccountAt() got = %v, want %v", got.Address(), tt.wantAddress)
			}
		})
	}
}

func TestAccountWithPrivateKey2(t *testing.T) {
	acc, err := NewAccountWithMnemonic(testcase.M1, ChainMainnet, AddressTypeComingTaproot)
	require.Nil(t, err)
	privateKey, err := acc.PrivateKeyHex()
	require.Nil(t, err)

	acc2, err := AccountWithPrivateKey(privateKey, ChainMainnet, AddressTypeComingTaproot)
	require.Nil(t, err)
	require.Equal(t, acc.privateKey, acc2.privateKey)
	require.Equal(t, acc.address, acc2.address)
}

func TestIsValidAddress(t *testing.T) {
	type args struct {
		chainnet string
		address  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "mainnet valid",
			args: args{chainnet: ChainMainnet, address: accountCase.addrMainnet},
			want: true,
		},
		{
			name: "signet valid",
			args: args{chainnet: ChainSignet, address: accountCase.addrSignet},
			want: true,
		},
		{
			name: "mainnet valid check in signet",
			args: args{chainnet: ChainSignet, address: accountCase.addrMainnet},
			want: true,
		},
		{
			name: "signet valid check in mainnet",
			args: args{chainnet: ChainMainnet, address: accountCase.addrSignet},
			want: true,
		},
		{
			name: "error address",
			args: args{chainnet: ChainMainnet, address: "bc1p5uslzuqy8k40mc86jfdtdjh4624umtw"},
			want: false,
		},
		{
			name: "empty address",
			args: args{chainnet: ChainSignet, address: ""},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAddress(tt.args.address, tt.args.chainnet); got != tt.want {
				t.Errorf("IsValidAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBTCWallet_Privatekey_Publickey_Address(t *testing.T) {
	// 从 coming 的 musig 库计算的测试用例
	// private key = 0xc7fceb75bafba7aa10ffe10315352bfc523ac733f814e6a311bc736873df8923
	// public key = 0x04a721f170043daafde0fa925ab6caf5d2abcdadd2249291b1840e3d99a3f41149e13185ef52451eef2e7cc0c5fe4180b64ca2d17eb886b2328518f6aed684719a
	// mainnet address = bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg
	// signet address = tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728

	phrase := "unaware oxygen allow method allow property predict various slice travel please priority"
	data, _ := bip39.NewSeedWithErrorChecking(phrase, "")

	pri, pub := btcec.PrivKeyFromBytes(data)
	priHex := types.HexEncodeToString(pri.Serialize())
	pubHex := types.HexEncodeToString(pub.SerializeUncompressed())
	t.Log("private key = ", priHex)
	t.Log("public key = ", pubHex)

	pubData := pub.SerializeUncompressed()
	addressHash, err := btcutil.NewAddressTaproot(pubData[1:33], &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("mainnet address = ", addressHash.EncodeAddress())

	addressHash, err = btcutil.NewAddressTaproot(pubData[1:33], &chaincfg.SigNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("signet address = ", addressHash.EncodeAddress())
}

func TestAccountWithPrivateKey(t *testing.T) {
	type args struct {
		prikey  string
		chain   string
		addType AddressType
	}
	tests := []struct {
		name     string
		args     args
		wantAddr string
		wantErr  bool
	}{
		{
			name: "Native Segwit(P2WPKH)",
			args: args{
				prikey:  "cPLpgDV8njCYGWCrXtvfSXo8fBkiCuoDjXfYbawNNaQkF3RyT2Km",
				chain:   ChainSignet,
				addType: AddressTypeNativeSegwit,
			},
			wantAddr: "tb1qcal96xxt64xtl0hp55erejn4awnmyx9c88nnmh",
		},
		{
			name: "Nested Segwit(P2SH-P2WPKH)",
			args: args{
				prikey:  "cMkXm38MtiLpeorUNtmMt5rrvfUZXkmyYtEtirEFsLFGVmWRThWq",
				chain:   ChainSignet,
				addType: AddressTypeNestedSegwit,
			},
			wantAddr: "2N489AZCJpazr2xLEygsGwUKbxixvUZaV6P",
		},
		{
			name: "Taproot (P2TR)",
			args: args{
				prikey:  "cSyGeGDKpaw6Y6vqJMDzVaN73YYZT64koA2JBuiifckAnhGS6SHZ",
				chain:   ChainSignet,
				addType: AddressTypeTaproot,
			},
			wantAddr: "tb1pdq423fm5dv00sl2uckmcve8y3w7guev8ka6qfweljlu23mmsw63qk6w2v3",
		},
		{
			name: "Legacy (P2PKH)",
			args: args{
				prikey:  "cTkZaPpb1pDdor36V5VY4uu5LE6tgzrjRADvrEXimEqWqvwRbfXY",
				chain:   ChainSignet,
				addType: AddressTypeLegacy,
			},
			wantAddr: "mxZX45K9oFMdJBpJXSVieMT3Wof3sCWUB6",
		},
		{
			name: "Taproot (P2TR) 2 use legacy pri",
			args: args{
				prikey:  "cTkZaPpb1pDdor36V5VY4uu5LE6tgzrjRADvrEXimEqWqvwRbfXY",
				chain:   ChainSignet,
				addType: AddressTypeTaproot,
			},
			wantAddr: "tb1p4re2ndazkwe7dmmpz3jhfk0rf9gv3y5z4kf3392d889ytqg5q67s3yuh2u",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := AccountWithPrivateKey(tt.args.prikey, tt.args.chain, tt.args.addType)
			if (err != nil) != tt.wantErr {
				t.Errorf("AccountWithPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err == nil) && acc.Address() != tt.wantAddr {
				t.Errorf("AccountWithPrivateKey() got = %v, want %v", acc.Address(), tt.wantAddr)
			}
		})
	}
}

func TestPublicKeyTransform(t *testing.T) {
	pubkeyCompressed := "0x02cfb7f626025d6826253f8fc5858e7a5c4b853350b4385ae909cf66138b71ec77"
	pubkeyUncompressed := "0x04cfb7f626025d6826253f8fc5858e7a5c4b853350b4385ae909cf66138b71ec772bac329b29f6eca2e1d25170a346c041e3dca5ebd68a91ac36cc53cd3be4908c"

	var pubkey string
	var err error

	pubkey, err = PublicKeyTransform(pubkeyCompressed, true)
	require.NoError(t, err)
	require.Equal(t, pubkey, pubkeyCompressed)
	pubkey, err = PublicKeyTransform(pubkeyCompressed, false)
	require.NoError(t, err)
	require.Equal(t, pubkey, pubkeyUncompressed)
	pubkey, err = PublicKeyTransform(pubkeyUncompressed, true)
	require.NoError(t, err)
	require.Equal(t, pubkey, pubkeyCompressed)
	pubkey, err = PublicKeyTransform(pubkeyUncompressed, false)
	require.NoError(t, err)
	require.Equal(t, pubkey, pubkeyUncompressed)

	pubkeyErr := "0x02cfb7f626025d6826253f8fc5858e7a5c4b853350b4385ae909cf66138b71ec71"
	pubkey, err = PublicKeyTransform(pubkeyErr, true)
	require.Error(t, err)
}

func TestIsValidPrivateKey(t *testing.T) {
	valid := IsValidPrivateKey("cTkZaPpb1pDdor36V5VY4uu5LE6tgzrjRADvrEXimEqWqvwRbfXY")
	t.Log(valid)
}

func TestAccount_SignMessage(t *testing.T) {
	acc, err := NewAccountWithMnemonic(testcase.M1, ChainMainnet, AddressTypeTaproot)
	require.Nil(t, err)
	t.Log(acc.Address())

	// sign message
	message := "hello world~"
	signature, err := acc.SignMessage(message)
	require.Nil(t, err)
	t.Log("sign message result: ", signature.Value)

	// check signature
	valid := VerifySignature(acc.PublicKeyHex(), message, signature.Value)
	t.Log("signature is valid: ", valid)
}

func TestAccount_SignPsbt(t *testing.T) {
	old, err := NewAccountWithMnemonic(testcase.M1, ChainMainnet, AddressTypeComingTaproot)
	require.NoError(t, err)
	priHex, _ := old.PrivateKeyHex()
	acc, err := AccountWithPrivateKey(priHex, ChainMainnet, AddressTypeTaproot)
	require.NoError(t, err)
	t.Log(acc.Address())

	psbtHex := "70736274ff01007d010000000158d38c41272f90e01d8dd3c5f9f6c4d37892ddcbaee155ecd826113c163003700000000000fdffffff02ea0500000000000022512073c9f2168a01fb0f0caa8a3fb7889ce4ab2cec67bfb16d272af4eea91fbd83e011cc240000000000160014f23a59a174bf387281535e2c0f6a395b77eb04a3000000000001011f2dd3240000000000160014f23a59a174bf387281535e2c0f6a395b77eb04a3000000"
	txn, err := acc.SignPsbt(psbtHex)
	require.NoError(t, err)
	require.True(t, txn.Packet.IsComplete())

	res, err := txn.PsbtHexString()
	require.NoError(t, err)
	t.Log("signed result: ", res.Value)
}

func TestChain_PushPsbt(t *testing.T) {
	chain, err := NewChainWithChainnet(ChainTestnet)
	require.NoError(t, err)

	psbtHex := "70736274ff01007d010000000158d38c41272f90e01d8dd3c5f9f6c4d37892ddcbaee155ecd826113c163003700000000000fdffffff02ea0500000000000022512073c9f2168a01fb0f0caa8a3fb7889ce4ab2cec67bfb16d272af4eea91fbd83e011cc240000000000160014f23a59a174bf387281535e2c0f6a395b77eb04a3000000000001011f2dd3240000000000160014f23a59a174bf387281535e2c0f6a395b77eb04a3000000"
	hash, err := chain.PushPsbt(psbtHex)
	require.NoError(t, err)
	t.Log("txn hash = ", hash.Value)
}
