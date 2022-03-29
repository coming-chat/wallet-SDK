package wallet

import (
	"testing"
)

func TestAddressToPublicKey(t *testing.T) {
	pk, err := AddressToPublicKey(address44)
	if err != nil {
		t.Fatal(err)
	}
	if pk != testPubKey {
		t.Fatal(pk)
	}
}

func TestPublicKeyToAddress(t *testing.T) {
	as, err := PublicKeyToAddress(testPubKey, 44)
	if err != nil {
		t.Fatal(err)
	}
	if as != address44 {
		t.Fatal(as)
	}
}

func TestValidAddress(t *testing.T) {
	addArray := []string{
		"5T8E3ZgvtHdZfEwsfZ9bZE5VCixUfCYS764We7xVuvQDbVrU",
		"5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnH",

		"5RNt3DACYRhwHyy9esTZXVvffkFL3pQHv4qoEMFVfDqeDEnB", // 换了一个字符，错误的地址
		"5RNt3DACYRhwHyy9esTZXVvffkFL3pQH",
	}

	for _, item := range addArray {
		pubkey, err := AddressToPublicKey(item)
		if err != nil {
			t.Log("false address, ", err)
		} else {
			t.Log("true address, pubkey=", pubkey)
		}
	}
}
