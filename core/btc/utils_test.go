package btc

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestValidAddress(t *testing.T) {
	addArray := [][2]string{
		// 有效的
		{"signet", "tb1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ystm5728"},
		{"signet", "tb1p4fwg0qlcsm94y90gnkwr0zkfsv9gxjlq43mpegf4cmn9xed02xcq3n0386"},
		{"mainnet", "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sg"},

		// 错误的
		{"signet", "tb1p4fwg0qlcsm94y90gnkwr0zkfsv9gxjlq43mpegf4cmn9xed02xcq3n0387"},
		{"signet", "tb1p4fwg0qlcsm94y90gnkwr0zkfsv9gxjlq43mpeg"},
		{"mainnet", "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7engl5z9ysunz3sh"},
		{"mainnet", "bc1p5uslzuqy8k40mc86jfdtdjh4624umtwjyjffrvvypc7e"},
	}

	for _, item := range addArray {
		net := item[0]
		addr := item[1]

		var cfg *chaincfg.Params
		switch net {
		case "signet":
			cfg = &chaincfg.SigNetParams
		case "mainnet":
			cfg = &chaincfg.MainNetParams
		}
		addobj, err := btcutil.DecodeAddress(addr, cfg)

		if err != nil {
			t.Log("false address, error = ", err)
		} else {
			t.Log("true address, address = ", addobj)
		}

		t.Log(IsValidAddress(addr, net))
	}

}
