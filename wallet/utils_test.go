package wallet

import "testing"

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
