package eth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDonut(t *testing.T) {
	owner := "0xC91A369B080638a1e1D0cFC81f3c420414E66aEe"
	ins, err := FetchDonutInscriptions(owner, "")
	require.Nil(t, err)
	t.Log(ins.JsonString())
}
