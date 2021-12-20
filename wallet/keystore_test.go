package wallet

import (
	"encoding/hex"
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
