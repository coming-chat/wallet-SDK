package solanaswap

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/stretchr/testify/require"
)

func Test_Drive_Address(t *testing.T) {
	programId := "whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"
	whirlpoolConfigKey := "FcrweFY1G9HJAHG5inkGB6pKg1HZ6x9UC2WioAfWrGkR"
	mintA := "281LhxeKQ2jaFDx9HAHcdrU9CpedSH7hx5PuRrM7e1FS" // test
	mintB := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU" // usdc
	tickSpacing := uint16(128)
	pub, err := getWhirlpoolPDA(programId, whirlpoolConfigKey, mintA, mintB, tickSpacing)
	require.Nil(t, err)
	addr := pub.ToBase58()
	want := "b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV"
	require.Equal(t, addr, want)
}

func TestGetPoolData(t *testing.T) {
	cli := devChain.Client()
	poolAddr := "b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV"
	info, err := getPoolData(cli, poolAddr)
	require.Nil(t, err)
	t.Log(info)
}

func Test_GetPoolsData(t *testing.T) {
	cli := devChain.Client()
	addresses := []string{
		common.PublicKey{0x11, 0x32}.ToBase58(),
		"b3D36rfrihrvLmwfvAzbnX9qF1aJ4hVguZFmjqsxVbV",
	}
	res, err := getPoolsData(cli, addresses)
	require.Nil(t, err)
	t.Log(res)
}

func TestDecode(t *testing.T) {
	dataHex := "3f95d10ce1806309d9336a3df48f361e5706e69c3cb6b6d91774e47935c8526de5a0f59f215a236afe80008000a00f2c01fc339a930e0000000000000000000000a734ab6d4a585e000000000000000000c500feff00fa641900000000df0100000000000010a72302b3fed346d77240c165c64c7aafa5012ada611aad6ddd14829c9bd02d27304ac7f1c702065761f500bf4ffce3be605b35ea03c24e5307dcb02fa0079bc4715e427008403f00000000000000003b442cb3912157f13a933d0134282d032b5ffecd01a2dbf1b7790608df002ea7782efcd9cd735079635206f82abbb54ccb9189268aaced48f64801830fd661bda0e99df70205000000000000000000004a523f65000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029ba11d131591598051342b3d56d475488fc9cd209d220b8758c531d4deed01a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029ba11d131591598051342b3d56d475488fc9cd209d220b8758c531d4deed01a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000029ba11d131591598051342b3d56d475488fc9cd209d220b8758c531d4deed01a0000000000000000000000000000000000000000000000000000000000000000"
	data, _ := hex.DecodeString(dataHex)

	var d WhirlpoolData
	err := d.Deserializer(data)
	require.Nil(t, err)

	bbb, err := json.Marshal(d)
	t.Log(string(bbb))
}

func Test_findWhirlpoolData(t *testing.T) {
	programId := "whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"
	whirlpoolConfigKey := "FcrweFY1G9HJAHG5inkGB6pKg1HZ6x9UC2WioAfWrGkR"
	mintA := "281LhxeKQ2jaFDx9HAHcdrU9CpedSH7hx5PuRrM7e1FS" // test
	mintB := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU" // usdc

	cli := devChain.Client()

	data, err := SearchWhirlpool(cli, programId, whirlpoolConfigKey, mintA, mintB)
	require.Nil(t, err)
	t.Log(data)
}
