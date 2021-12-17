package wallet

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestPolDecode(t *testing.T) {
	decodeByte, err := base64.RawStdEncoding.DecodeString("jC9MOH7OPYbHdJtiOWFW0lpMUCFO4nASKjzqHvXpEiYAgAAAAQAAAAgAAACm2Dm/CZ98R1uy34lMj7tr9+i3ERCFoeCSdNwOScsyDkvLwhVGv6qxOzmdiR7vzgRgEizMQbq17k0C1Tk59WyDnf9OfaGQTenQQpnFPiXxcmDa6TXQvF7Eq8VYw009ANLmDTIQ125JdQX6edYY85ZFpLiOltXiad44mhS1mC8OSCcOHsViVrk3Lk0eMsClYS1SUzv3QDCoHChFu6Za")
	var s = "111"
	publicKey, secretKeys, err := decodePolkaKeystore(&s, decodeByte)
	if err != nil {
		t.Error(err)
	}
	sr, _ := hex.DecodeString(testSecretSeed)
	t.Log(hex.EncodeToString(publicKey), hex.EncodeToString(secretKeys), sr)
	msg := []byte("8asd8u8qw9ddqu9w8d9wqud89q9wd8uq89uw8u89r893h22")
	sign, err := Sign(secretKeys, signingContext(msg), publicKey)
	if err != nil {
		t.Error(err)
	}
	signedData := sign.Encode()

	t.Log(hex.EncodeToString(signedData[:]))
}
