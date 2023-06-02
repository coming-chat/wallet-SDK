package btc

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeOrdFromWitness(t *testing.T) {
	witness, err := hex.DecodeString("208651980c1703275f50edf7e024e8ec5acffe21003250704dfe8cfa07333ea545ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800307b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2232323232222c22616d74223a317d68")
	require.NoError(t, err)
	ord, err := DecodeOrdFromWitness(witness)
	require.NoError(t, err)
	t.Log(ord)
}
