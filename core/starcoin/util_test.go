package starcoin

import (
	"math/big"
	"testing"

	"github.com/starcoinorg/starcoin-go/client"
	"github.com/stretchr/testify/require"
)

func TestUint128(t *testing.T) {
	number := "1234567890987654321"
	numInt, _ := big.NewInt(0).SetString(number, 10)

	sdkU128, _ := client.BigIntToU128(numInt)
	sdkRestoreInt := client.U128ToBigInt(sdkU128)
	t.Log(sdkRestoreInt.String())

	myU128, err := NewU128FromString(number)
	myRestoreInt := client.U128ToBigInt(myU128)
	t.Log(myRestoreInt.String())
	if err != nil {
		t.Log(err)
		return
	}
	require.Equal(t, numInt, myRestoreInt)
}
