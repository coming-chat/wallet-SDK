package aptos

import (
	"strconv"
	"testing"
)

func Test_calcTotal(t *testing.T) {
	type args struct {
		amount   uint64
		feePoint uint64
	}
	tests := []struct {
		name string
		args args
	}{}
	for i := uint64(10000000); i < 10000000+50000; i++ {
		tests = append(tests, struct {
			name string
			args args
		}{
			name: "test" + strconv.FormatUint(i, 10),
			args: args{
				amount:   i,
				feePoint: 250,
			},
		})
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcTotal(tt.args.amount, tt.args.feePoint); got-got/10000*250 != tt.args.amount {
				t.Errorf("calcTotal() = %v", got)
			}
		})
	}
}
