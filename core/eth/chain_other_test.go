package eth

import (
	"encoding/hex"
	"testing"
)

func TestChain_CallContract(t *testing.T) {
	type fields struct {
		RpcUrl string
	}
	type args struct {
		msg         *CallMsg
		blockNumber string
	}
	dataByte, err := hex.DecodeString("70a082310000000000000000000000008de5ff2eded4d897da535ab0f379ec1b9257ebab")
	if err != nil {
		return
	}
	callmsg := NewCallMsg()
	callmsg.SetTo("0x37088186089c7d6bcd556d9a15087dfae3ba0c32")
	callmsg.SetData(dataByte)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "tt",
			fields: fields{
				RpcUrl: rpcs.sherpaxProd.url,
			},
			args: args{
				msg:         callmsg,
				blockNumber: "1670427",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Chain{
				RpcUrl: tt.fields.RpcUrl,
			}
			got, err := c.CallContract(tt.args.msg, tt.args.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CallContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CallContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
