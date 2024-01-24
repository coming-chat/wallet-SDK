package eth

import (
	"errors"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) TransferNFT(sender, receiver string, nft *base.NFT) (*Transaction, error) {
	return c.TransferNFTParams(sender, receiver, nft.Id, nft.ContractAddress, nft.Standard)
}

// TransferNFTParams
// - param nftStandard: only support erc-721 now, else throw error unsupported nft type.
func (c *Chain) TransferNFTParams(sender, receiver, nftId, nftContractAddress, nftStandard string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if strings.ToLower(nftStandard) != "erc-721" {
		return nil, errors.New("unsupported nft type")
	}
	data, err := EncodeErc721TransferFrom(sender, receiver, nftId)
	if err != nil {
		return nil, err
	}
	gasPrice, err := c.SuggestGasPrice()
	if err != nil {
		return nil, err
	}

	msg := NewCallMsg()
	msg.SetFrom(sender)
	msg.SetTo(nftContractAddress)
	msg.SetValue("0")
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(data)

	gasLimit, err := c.EstimateGasLimit(msg)
	if err != nil {
		return nil, err
	}
	msg.SetGasLimit(gasLimit.Value)

	return msg.TransferToTransaction(), nil
}
