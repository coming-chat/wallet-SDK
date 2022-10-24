package sr25519

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
				publicKey: []byte{0x24, 0xc5, 0x8f, 0x3b, 0x78, 0xbe, 0xe0, 0x7f, 0x93, 0xce, 0xd6, 0x8b, 0x1f, 0x74,
					0x86, 0x5a, 0xbf, 0x59, 0x84, 0xd0, 0x32, 0x14, 0xbd, 0xe8, 0xcc, 0x41, 0x16, 0x5e,
					0x2a, 0xc8, 0x61, 0x25},
				msg:       []byte{0x12},
				signature: []byte{132, 118, 213, 103, 77, 209, 185, 218, 184, 233, 73, 27, 123, 237, 187, 25, 40, 28, 138, 254, 108, 205, 180, 137, 44, 149, 126, 197, 25, 180, 12, 61, 149, 193, 158, 103, 13, 175, 96, 1, 83, 94, 98, 118, 222, 240, 210, 43, 4, 34, 70, 93, 162, 134, 162, 146, 30, 52, 127, 35, 220, 121, 199, 131},
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
