package ed25519

import "testing"

func TestIsValidSignature(t *testing.T) {
	type args struct {
		publicKey []byte
		msg       []byte
		signature []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test 1",
			args: args{
				publicKey: []byte{0x36, 0x48, 0x6b, 0xb2, 0xa5, 0xa4, 0x83, 0xd9, 0x91, 0xa0, 0x09, 0x99, 0x45, 0xc4, 0x40, 0x31, 0x7b, 0xa2, 0xc3, 0xa7, 0x85, 0x32, 0x86, 0x22, 0xb9, 0x43, 0x1b, 0x0d, 0x1d, 0x1c, 0xa3, 0xf1},
				msg:       []byte{0x12},
				signature: []byte{85, 105, 58, 89, 193, 218, 236, 85, 11, 70, 73, 13, 77, 224, 117, 53, 73, 187, 124, 125, 91, 107, 186, 177, 206, 82, 253, 216, 11, 227, 118, 38, 67, 220, 73, 54, 8, 125, 219, 63, 168, 223, 102, 161, 39, 108, 175, 97, 155, 54, 16, 29, 134, 3, 3, 220, 95, 51, 235, 205, 252, 105, 71, 13},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidSignature(tt.args.publicKey, tt.args.msg, tt.args.signature); got != tt.want {
				t.Errorf("IsValidSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
