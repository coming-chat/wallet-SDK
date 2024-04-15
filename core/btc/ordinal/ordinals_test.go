package ordinal

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeOrdFromWitness(t *testing.T) {
	witness, err := hex.DecodeString("20f4df9e108244fab3e423ab3fab0bdacb9047e4443a5d68195d256ea007644ba4ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d38004c4d7b200d0a20202270223a20226272632d3230222c0d0a2020226f70223a20227472616e73666572222c0d0a2020227469636b223a2022696e7363222c0d0a202022616d74223a202232220d0a7d68")
	require.NoError(t, err)
	ord, err := DecodeOrdFromWitness(witness)
	require.NoError(t, err)
	t.Log(ord)
}
