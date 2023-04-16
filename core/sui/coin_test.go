package sui

import (
	"math/big"
	"testing"

	"github.com/coming-chat/go-sui/types"
	"github.com/stretchr/testify/require"
)

func Test_pickupTransferCoin(t *testing.T) {
	coin1000 := types.Coin{Balance: types.NewSafeSuiBigInt[uint64](1000)}
	coin800 := types.Coin{Balance: types.NewSafeSuiBigInt[uint64](800)}
	coin500 := types.Coin{Balance: types.NewSafeSuiBigInt[uint64](500)}
	coin100 := types.Coin{Balance: types.NewSafeSuiBigInt[uint64](100)}
	coin1 := types.Coin{Balance: types.NewSafeSuiBigInt[uint64](1)}
	coins := []types.Coin{
		coin1000, coin800, coin500, coin100, coin1,
	}
	tests := []struct {
		name    string
		amount  uint64
		want    *PickedCoins
		wantErr bool
	}{
		{
			name:   "can transfer object 1",
			amount: 1000,
			want: &PickedCoins{
				Coins:                []types.Coin{coin1000},
				Total:                *big.NewInt(0).SetUint64(1000),
				Amount:               *big.NewInt(0).SetUint64(1000),
				CanUseTransferObject: true,
			},
		},
		{
			name:   "can transfer object 2",
			amount: 500,
			want: &PickedCoins{
				Coins:                []types.Coin{coin500},
				Total:                *big.NewInt(0).SetUint64(500),
				Amount:               *big.NewInt(0).SetUint64(500),
				CanUseTransferObject: true,
			},
		},
		{
			name:   "can transfer object 3",
			amount: 1,
			want: &PickedCoins{
				Coins:                []types.Coin{coin1},
				Total:                *big.NewInt(0).SetUint64(1),
				Amount:               *big.NewInt(0).SetUint64(1),
				CanUseTransferObject: true,
			},
		},
		{
			name:   "can transfer 1",
			amount: 1001,
			want: &PickedCoins{
				Coins:  []types.Coin{coin1000, coin800},
				Total:  *big.NewInt(0).SetUint64(1800),
				Amount: *big.NewInt(0).SetUint64(1001),
			},
		},
		{
			name:   "can transfer 2",
			amount: 1900,
			want: &PickedCoins{
				Coins:  []types.Coin{coin1000, coin800, coin500},
				Total:  *big.NewInt(0).SetUint64(2300),
				Amount: *big.NewInt(0).SetUint64(1900),
			},
		},
		{
			name:    "insufficient account balance",
			amount:  10000,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pickupTransferCoin(coins, tt.amount, false)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.Equal(t, got, tt.want)
			}
		})
	}
}
