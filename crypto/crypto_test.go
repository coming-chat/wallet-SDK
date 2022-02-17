package crypto

import (
	"encoding/hex"
	"fmt"
	"github.com/coming-chat/wallet-SDK/u8util"
	"reflect"
	"testing"
)

func TestNewSHA256(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte("abc"), "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{[]byte("hello"), "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewSHA256(tt.in)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewBlake2b256(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8"},
		{[]byte("abc"), "bddd813c634239723171ef3fee98579b94964e3bb1cb3e427262c8c068d52319"},
		{[]byte("hello"), "324dcf027dd4a30a932c441f365a25e86b173defa4b8e58948253471b81b72cf"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewBlake2b256(tt.in)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewBlake2b512(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "786a02f742015903c6c6fd852552d272912f4740e15847618a86e217f71f5419d25e1031afee585313896444934eb04b903a685b1448b755d56f701afe9be2ce"},
		{[]byte("abc"), "ba80a53f981c4d0d6a2797b69f12f6e94c212f14685ac4b74b12bb6fdbffa2d17d87c5392aab792dc252d5de4533cc9518d38aa8dbf1925ab92386edd4009923"},
		{[]byte("hello"), "e4cfa39a3d37be31c59609e807970799caa68a19bfaa15135f165085e01d41a65ba1e1b146aeb6bd0092b49eac214c103ccfa3a365954bbbe52f74a2b3620c94"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewBlake2b512(tt.in)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewXXHash(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "99e9d85137db46ef"},
		{[]byte("abc"), "990977adf52cbc44"},
		{[]byte("hello"), "a36d9f887d82c726"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewXXHash(tt.in, 64)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewXXHash64(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "99e9d85137db46ef"},
		{[]byte("abc"), "990977adf52cbc44"},
		{[]byte("hello"), "a36d9f887d82c726"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewXXHash64(tt.in)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewBlake2b256Sig(t *testing.T) {
	for i, tt := range []struct {
		key  []byte
		data []byte
		out  string
	}{
		{nil, []byte("abc"), "bddd813c…68d52319"},
		{[]byte{4, 5, 6}, []byte{1, 2, 3}, "af0e60f4…8a714a8f"},
		{[]byte("abc"), nil, "7a78f945…73b8f07b"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result, err := NewBlake2b256Sig(tt.key, tt.data)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(u8util.ToHex(result[:], 512/8, false), tt.out) {
				t.Fatal(u8util.ToHex(result[:], 512/8, false), tt.out)
			}
		})
	}
}

func TestNewXXHash128(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "99e9d85137db46ef4bbea33613baafd5"},
		{[]byte("abc"), "990977adf52cbc440889329981caa9be"},
		{[]byte("hello"), "a36d9f887d82c726b2a1d004cb71dd23"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewXXHash128(tt.in)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewXXHash256(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "99e9d85137db46ef4bbea33613baafd56f963c64b1f3685a4eb4abd67ff6203a"},
		{[]byte("abc"), "990977adf52cbc440889329981caa9bef7da5770b2b8a05303b75d95360dd62b"},
		{[]byte("hello"), "a36d9f887d82c726b2a1d004cb71dd231fe2fb3bf584fc533914a80e276583e0"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := NewXXHash256(tt.in)
			if hex.EncodeToString(result[:]) != tt.out {
				t.Fatalf("want %v; got %v", tt.out, hex.EncodeToString(result[:]))
			}
		})
	}
}

func TestNewBlake2b512Sig(t *testing.T) {
	for i, tt := range []struct {
		key  []byte
		data []byte
		out  string
	}{
		{nil, []byte("abc"), "ba80a53f…d4009923"},
		{[]byte{4, 5, 6}, []byte{1, 2, 3}, "7ed2dd42…2ec9362e"},
		{[]byte("abc"), nil, "91cc35fc…f47b00e5"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result, err := NewBlake2b512Sig(tt.key, tt.data)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(u8util.ToHex(result[:], 512/8, false), tt.out) {
				t.Fatal(u8util.ToHex(result[:], 512/8, false), tt.out)
			}
		})
	}
}

func TestNaclVerify(t *testing.T) {
	signature := []byte{28, 58, 206, 239, 249, 70, 59, 191, 166, 40, 219, 218, 235, 170, 25, 79, 10, 94, 9, 197, 34, 126, 1, 150, 246, 68, 28, 238, 36, 26, 172, 163, 168, 90, 202, 211, 126, 246, 57, 212, 43, 24, 88, 197, 240, 113, 118, 76, 37, 81, 91, 110, 236, 50, 144, 134, 100, 223, 220, 238, 34, 185, 211, 7}

	publicKey, _, err := NewNaclKeyPairFromSeed([]uint8("12345678901234567890123456789012"))
	if err != nil {
		t.Fatal(err)
	}

	t.Run("validates a correctly signed message", func(t *testing.T) {
		digest := []byte{0x61, 0x62, 0x63, 0x64}
		if !NaclVerify(digest, signature, publicKey) {
			t.Fail()
		}
	})

	t.Run("fails a correctly signed message (message changed)", func(t *testing.T) {
		digest := []byte{0x61, 0x62, 0x63, 0x64, 0x65}
		if NaclVerify(digest, signature, publicKey) {
			t.Fail()
		}
	})

	t.Run("fails a correctly signed message (signature changed)", func(t *testing.T) {
		signature[0] = 0xff
		digest := []byte{0x61, 0x62, 0x63, 0x64}
		if NaclVerify(digest, signature, publicKey) {
			t.Fail()
		}
	})
}
