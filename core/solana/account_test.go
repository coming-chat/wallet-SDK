package solana

import (
	"testing"

	"github.com/coming-chat/wallet-SDK/core/testcase"
	"github.com/stretchr/testify/require"
)

func TestNewAccountWithMnemonic(t *testing.T) {
	tests := []struct {
		name    string
		acase   testcase.AccountCase
		wantErr bool
	}{
		{
			name:  "normal1",
			acase: testcase.Accounts.Solana,
		},
		{
			name:  "normal2",
			acase: testcase.Accounts2.Solana,
		},
		{
			name:    "empty mnemonic",
			acase:   testcase.EmptyMnemonic,
			wantErr: true,
		},
		{
			name:    "error mnemonic",
			acase:   testcase.ErrorMnemonic,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAccountWithMnemonic(tt.acase.Mnemonic)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("NewAccountWithMnemonic() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if got.Address() != tt.acase.Address {
				t.Errorf("NewAccountWithMnemonic() want address %v, got %v", tt.acase.Address, got.Address())
			}
		})
	}
}

func TestAccount_IsValidAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{
			name:    "normal",
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPY",
			want:    true,
		},
		{
			name:    "miss char",
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqP",
			want:    false,
		},
		{
			name:    "redundant char",
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqPYD",
			want:    false,
		},
		{
			name:    "invalid base58 char 0OIl",
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vC0OIl",
			want:    false,
		},
		{
			name:    "alter chars",
			address: "AfBfH4ehvcXx66Y5YZozgTYPC1nieL9A3r2yT3vCXqQQ",
			want:    true, // Maybe it's wrong, but we can't tell
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidAddress(tt.address)
			if got != tt.want {
				t.Errorf("Account.IsValidAddress() %v got = %v, want %v", tt.address, got, tt.want)
			}
		})
	}
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
