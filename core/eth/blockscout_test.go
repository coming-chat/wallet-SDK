package eth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockScout_NFT(t *testing.T) {
	api := NewBlockScout(BlockScoutURLEth)
	owner := "0x30b31174e5FEd4598aA0dF3191EbAA5AAb48d43E" // nft very much
	// owner := "0x9dE416AB881b322eE0b4c189C2dE624090280cF2" // nft few

	showPage := func(p *BKSNFTPage) {
		println("item count: ", p.Count())
		f := p.ValueAt(0).ToBaseNFT()
		println("first item: ", f.Id, f.Name, f.Image)
		println("has next page: ", p.HasNextPage())
		println("====================")
	}

	page1, err := api.Nft(owner, nil)
	require.Nil(t, err)
	showPage(page1)

	require.Equal(t, page1.HasNextPage(), true)
	page2, err := api.Nft(owner, page1.NextPageParams())
	require.Nil(t, err)
	showPage(page2)
}
