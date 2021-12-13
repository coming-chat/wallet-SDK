package wallet

import "testing"

func TestAddressToPublicKey(t *testing.T) {
	pk := AddressToPublicKey(address44)
	if pk != testPubKey {
		t.Fatal(pk)
	}
}

func TestPublicKeyToAddress(t *testing.T) {
	as := PublicKeyToAddress(testPubKey, 44)
	if as != address44 {
		t.Fatal(as)
	}
}
