package runes

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/wire"
	"math/big"
	"reflect"
	"testing"
)

func TestDecipher(t *testing.T) {
	type args struct {
		transaction string
	}
	tests := []struct {
		name    string
		args    args
		want    Artifact
		wantErr bool
	}{
		{
			name: "mint",
			args: args{
				transaction: "0200000000010193ca17b01a5256b8fdd28d387e1c32f9809312b3c7e8a3aeaab3f205fd53e8590000000000fdffffff021027000000000000225120bc3b91cde00e0200b157d637cd251dcd48ff90ca582fb06d0c67b0575d29dc4600000000000000000b6a5d0814a6e09d01148f010140178ebcc522b406817b6b6d76d1f5f66a75e90f15207adc1d54398c916fcf6b87c2b2e82bd7623502544ecfc23db39a836ff851e05d9f3f05ed8bc375fbd0f7ce00000000",
			},
		},
		{
			name: "deploy",
			args: args{
				transaction: "0200000000010152a68185bd38da39def6b2c10158da48928b496e1f855559d8d8c861a7fce556000000000006000000031027000000000000225120bc3b91cde00e0200b157d637cd251dcd48ff90ca582fb06d0c67b0575d29dc461027000000000000225120bc3b91cde00e0200b157d637cd251dcd48ff90ca582fb06d0c67b0575d29dc460000000000000000226a5d1f020304aff8b394dab7aec1c80f010306cc080ae80708e8071280c2d72f16010340a3bfbd5521c9021d7efbc03672eb537bc3698c39ac058eb8beae3d6a64c0f0fd5202267a332940ed56340fc05edf0172d6fef7c4aca4f6463b772757cb6c35dbfd4401201afefd27eb8a1ccd6b5e7684632c91a14f57113671d7e6e2d7981a302dadce1aac0063036f726401010d696d6167652f7376672b786d6c010200010d092ffc8ca2bdb982c807004cf93c3f786d6c2076657273696f6e3d22312e302220656e636f64696e673d225554462d38223f3e0a3c73766720786d6c6e733d22687474703a2f2f7777772e77332e6f72672f323030302f737667222076696577426f783d222d31303030202d313030302032303030203230303022207374796c653d226261636b67726f756e642d636f6c6f723a7768697465223e0a20203c636972636c6520723d22393630222066696c6c3d22776869746522207374726f6b653d22626c61636b22207374726f6b652d77696474683d223830222f3e0a20203c636972636c6520723d22363735222066696c6c3d22626c61636b222f3e0a3c2f7376673e0a6821c11afefd27eb8a1ccd6b5e7684632c91a14f57113671d7e6e2d7981a302dadce1a00000000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := wire.NewMsgTx(wire.TxVersion)
			raw, err := hex.DecodeString(tt.args.transaction)
			if err != nil {
				t.Error(err)
			}
			if err = tx.Deserialize(bytes.NewReader(raw)); err != nil {
				t.Error(err)
			}
			got, err := Decipher(tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decipher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got.(Runestone).Etching.Rune.String())
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decipher() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBigInt(t *testing.T) {
	a, _ := new(big.Int).SetString("340282366920938463463374607431768211455", 10)
	t.Log(a.BitLen())

	maxU128 := big.NewInt(1)

	// 左移 64 位，相当于乘以 2^64，再加上 64 位的最大值

	//maxU128.Add(maxU128, big.NewInt(0xFFFFFFFFFFFFFFFF))
	t.Log(new(big.Int).Sub(big.NewInt(1).Lsh(maxU128, 128), big.NewInt(1)).String())
}
