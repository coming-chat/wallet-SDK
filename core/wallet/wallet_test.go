package wallet

import "testing"

const (
	testSecretPhrase = "rookie october miracle crisp invest grace birth exile black attitude bitter napkin"
	testSecretSeed   = "0x167d9a020688544ea246b056799d6a771e97c9da057e4d0b87024537f99177bc"
	testPubKey       = "0xdc64bef918ddda3126a39a11113767741ddfdf91399f055e1d963f2ae1ec2535"
	address44        = "5QUEnWNMDFqsbUGpvvtgWGUgiiojnEpLf7581ELLAQyQ1xnT"
	keystore2        = "{\"encoded\":\"5zmfXmtpiz8sryDmupYcoFDDCRj0ufe1Fx1EfGFLQoMAgAAAAQAAAAgAAADajJFtVRycQELlG4KibfgTOX4zexng/E3oj+I+ND9GYQIcHnIrEfAu1Ptcoi1HLiM8GfKuzcmMg9ZEvhywWF1Hau4XThv8pk8xGQUyMn2iMQtV8JA/5SGL/w5r5bT9vPOsidQEkc4Q5RvEsqjeU0hCkGKQXIui/9DqFR02Dq9pn3KYK3EQNjkNZplBJ59h4pG+E6SNMG8XuKqDMn+b\",\"encoding\":{\"content\":[\"pkcs8\",\"sr25519\"],\"type\":[\"scrypt\",\"xsalsa20-poly1305\"],\"version\":\"3\"},\"address\":\"5UczqUVGsoQpZnBCZkDtxvLxJ42KnUfaGTzPkQmZeAAug4s9\",\"meta\":{\"genesisHash\":\"0x96675ae0e91fe7d102f8eebc4ee4fbb9241b483bc6645ac975864684d1c222ff\",\"isHardware\":false,\"name\":\"wallet test\",\"tags\":[],\"whenCreated\":1645428018341}}"
	password         = "111"
)

func TestNewWallet(t *testing.T) {
	wallet, err := NewWallet(testSecretPhrase)
	if err != nil {
		t.Fatal(err)
	}
	sign, err := wallet.Sign([]byte{}, "")
	if err != nil {
		t.Fatal(sign)
	}
	t.Log(sign)
}

func TestGetPrivateKey(t *testing.T) {
	mnemonic := "rookie october miracle crisp invest grace birth exile black attitude bitter napkin"
	privateHex := "0xba865d03c9f6f27871d4eddd8baffe2b16c444945388b39adb0a0966020bbbbe"
	w, err := NewWallet(mnemonic)
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := w.GetPrivateKeyHex()
	if err != nil {
		t.Fatal(err)
	}
	if privateHex != privateKey {
		t.Fatal(privateKey)
	}
}
