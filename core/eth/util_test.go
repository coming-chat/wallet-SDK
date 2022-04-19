package eth

import (
	"crypto/ecdsa"
	"strconv"
	"testing"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestValidAddress(t *testing.T) {
	// 检测地址需要用 common.IsHexAddress
	// 结论：地址检测只检测了字符串长度、是否 16 进制
	addresses := []string{
		"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2", // 正确的
		"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c1", // 这个地址修改了末尾一个数，大概率应该是错误的，但是暂时无法检测

		// 错误的
		"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c",
		"0x7161ada3EA6e53E5652A",
		"",
	}

	for _, a := range addresses {
		address := common.HexToAddress(a)
		t.Log(common.IsHexAddress(a), address)
	}
}

type nnn int

func TestMapConcurrent(t *testing.T) {
	nums := []interface{}{1, 2, 3, 4, 5, 6}
	// nums := []interface{}{"1", "2", "3", "4"}
	res, _ := MapListConcurrent(nums, func(i interface{}) (interface{}, error) {
		return strconv.Itoa(i.(int) * 100), nil
	})
	t.Log(res)
}

func TestETHWallet_Privatekey_Publickey_Address(t *testing.T) {
	// 从 coming 的 trust wallet 库计算的测试用例
	// private key = 0x8c3083c24062f065ff2ee71b21f665375b266cebffa920e8909ec7c48006725d
	// public key  = 0xc66cbe3908fda67d2fb229b13a63aa1a2d8428acef2ff67bc31f6a79f2e2085f // Curve25519
	// public key  = 0xb34ec4ec2ebc84b04d9170bed91f65306c7045863efb9175d721104a8ecc17f2 // Ed25519
	// public key  = 0x011e56a004e205db53ae3cc7291ffb8a28181aed4b4e95813c17b9a96db2d769 // Ed25519Blake2b
	// public key  = 0x04bd6d7af856d20188fcfdb8ff38b978bc7c72fd028b67a6fab3d2120dd9bd1db61c5d44e242001dce224188a8b88150e16e9748438703bbf2dc417135c4f9377e // Secp256k1 compressed false
	// public key  = 0x02bd6d7af856d20188fcfdb8ff38b978bc7c72fd028b67a6fab3d2120dd9bd1db6 // Secp256k1 compressed true
	// public key  = 0x027bcb5a6edf262eca9602b8343baa1cd5dd7811e540e850b05661b6524e504222 // Nist256p1
	// address     = 0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2

	phrase := "unaware oxygen allow method allow property predict various slice travel please priority"
	seed, _ := bip39.NewSeedWithErrorChecking(phrase, "")

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}

	path, _ := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	key := masterKey
	for _, n := range path {
		key, err = key.DeriveNonStandard(n)
		if err != nil {
			t.Fatal(err)
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		t.Fatal(err)
	}
	privateKeyECDSA := privateKey.ToECDSA()
	privateKeyHex := types.HexEncodeToString(privateKey.Serialize())
	t.Log("private key = ", privateKeyHex)

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		t.Fatal(".....")
	}

	data := crypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHex := types.HexEncodeToString(data)
	t.Log("public key = ", publicKeyHex)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	t.Log("address = ", address.Hex())

}
