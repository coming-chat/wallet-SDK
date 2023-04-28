package sui

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"net/url"
	"strconv"

	"github.com/coming-chat/go-sui/sui_types"
	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) FetchSuiCatGlobalData() (data *SuiCatGlobalData, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.Client()
	if err != nil {
		return
	}
	globalId, err := types.NewAddressFromHex(SuiCatTestnetConfig.GlobalId)
	if err != nil {
		return
	}
	res, err := cli.GetObject(context.Background(), *globalId, &types.SuiObjectDataOptions{
		ShowContent: true,
	})
	if err != nil {
		return
	}
	if res.Data.Content == nil || res.Data.Content.Data.MoveObject == nil {
		return nil, errors.New("sui cat data not found")
	}

	contentData, err := json.Marshal(res.Data.Content.Data.MoveObject)
	if err != nil {
		return
	}
	var rawData struct {
		Fields rawSuiCatGlobalData `json:"fields"`
	}
	err = json.Unmarshal(contentData, &rawData)
	if err != nil {
		return
	}

	// save whitelist id
	SuiCatTestnetConfig.whitelistId = rawData.Fields.Whitelist.Fields.Id.Id

	data = rawData.Fields.mapToPubData()
	return
}

func (c *Chain) QueryIsInSuiCatWhiteList(address string) (b *base.OptionalBool, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.Client()
	if err != nil {
		return
	}

	if SuiCatTestnetConfig.whitelistId == "" {
		_, err = c.FetchSuiCatGlobalData()
		if err != nil {
			return
		}
	}

	parentId, err := types.NewAddressFromHex(SuiCatTestnetConfig.whitelistId)
	if err != nil {
		return
	}
	fieldName := sui_types.DynamicFieldName{
		Type:  "address",
		Value: address,
	}
	res, err := cli.GetDynamicFieldObject(context.Background(), *parentId, fieldName)
	if err != nil {
		if _, ok := err.(*url.Error); ok {
			return nil, err
		}
		return &base.OptionalBool{Value: false}, nil
	}

	if res.Error != nil || res.Data == nil || res.Data.Digest == "" {
		return &base.OptionalBool{Value: false}, nil
	} else {
		return &base.OptionalBool{Value: true}, nil
	}
}

func (c *Chain) MintSuiCatNFT(signer string, amount string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signerAddr, err := types.NewAddressFromHex(signer)
	if err != nil {
		return
	}
	amountInt, err := strconv.ParseUint(amount, 10, 64)
	if err != nil {
		return
	}
	add2, err := types.NewAddressFromHex("0x6")
	if err != nil {
		return
	}

	cli, err := c.Client()
	if err != nil {
		return
	}
	coins, err := cli.GetCoins(context.Background(), *signerAddr, nil, nil, MAX_INPUT_COUNT_MERGE)
	if err != nil {
		return
	}
	pickedCoins, err := types.PickupCoins(coins, *big.NewInt(0).SetUint64(amountInt), MAX_INPUT_COUNT_MERGE, true)
	if err != nil {
		return
	}
	return c.BaseMoveCall(signer,
		SuiCatTestnetConfig.PackageId,
		SuiCatTestnetConfig.ModuleName,
		SuiCatFuncMint,
		[]string{},
		[]any{
			SuiCatTestnetConfig.GlobalId,
			add2,
			pickedCoins.CoinIds(),
		})
}
