package sui

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestAccount(t *testing.T) {
	account := M1Account(t)
	require.Equal(t, account.Address(), "0x7e875ea78ee09f08d72e2676cf84e0f1c8ac61d94fa339cc8e37cace85bebc6e")

	t.Log(account.PrivateKeyHex())
	t.Log(account.PublicKeyHex())
	t.Log(account.Address())
}

func TestPublicKeyToAddress(t *testing.T) {
	account := M1Account(t)

	pub := account.PublicKeyHex()
	addr := account.Address()

	encodedAddr, err := EncodePublicKeyToAddress(pub)
	require.Nil(t, err)
	require.Equal(t, addr, encodedAddr)
}

// Account of os environment M1
func M1Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M1)
	require.Nil(t, err)
	return account
}

func M2Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M2)
	require.Nil(t, err)
	return account
}

func M3Account(t *testing.T) *Account {
	account, err := NewAccountWithMnemonic(testcase.M3)
	require.Nil(t, err)
	return account
}

func TestAccountWithPrivatekey(t *testing.T) {
	mnemonic := testcase.M1
	accountFromMnemonic, err := NewAccountWithMnemonic(mnemonic)
	require.Nil(t, err)
	privateKey, err := accountFromMnemonic.PrivateKeyHex()
	require.Nil(t, err)

	accountFromPrikey, err := AccountWithPrivateKey(privateKey)
	require.Nil(t, err)

	require.Equal(t, accountFromMnemonic.Address(), accountFromPrikey.Address())
}

func Test_IsValidAddress(t *testing.T) {
	addr := "0xd77955e670f42c1bc5e94b9e68e5fe9bdbed9134d784f2a14dfe5fc1b24b5d9f"
	valid := IsValidAddress(addr)
	require.True(t, valid)
}
