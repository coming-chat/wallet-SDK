package eth

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
	"github.com/coming-chat/wallet-SDK/graphql"
	"github.com/ethereum/go-ethereum/common"
)

type donutGraphResp struct {
	Data struct {
		Type  string `json:"@type"`
		Value string `json:"value"`
	} `json:"data,omitempty"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func donutParser(resp donutGraphResp, out any) error {
	if resp.Code != 0 {
		return fmt.Errorf("error code %v: %v", resp.Code, resp.Msg)
	}
	return json.Unmarshal([]byte(resp.Data.Value), out)
}

type DonutInscription struct {
	Tick   string `json:"tick"`
	Amount string `json:"amount"`
}

type DonutInscriptionArray struct {
	inter.AnyArray[*DonutInscription]
}

// FetchDonutInscriptions
// - param graphURL Default "https://bc.dnt.social/v1/common/search"
func FetchDonutInscriptions(owner string, graphURL string) (arr *DonutInscriptionArray, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if graphURL == "" {
		graphURL = "https://bc.dnt.social/v1/common/search"
	}
	query := fmt.Sprintf(`{
		src20Balances(holder: "%v", first: 100) {
			edges{
				node{
					tick
					amount
				}
			}
		}
	}`, owner)

	var out struct {
		Src20Balances struct {
			Edges []struct {
				Node *DonutInscription `json:"node"`
			} `json:"edges"`
		} `json:"src20Balances"`
	}
	err = graphql.QueryString(query, graphURL, donutParser, &out)
	if err != nil {
		return
	}
	inscriptions := make([]*DonutInscription, len(out.Src20Balances.Edges))
	for idx, node := range out.Src20Balances.Edges {
		inscriptions[idx] = node.Node
	}
	return &DonutInscriptionArray{AnyArray: inscriptions}, nil
}

func (c *Chain) BuildDonutTransfer(sender, receiver, tick, amount string) (*Transaction, error) {
	return c.BuildSrc20Transfer(sender, receiver, "0xf414dF7d8260A8e1e007F72892Cf5F0A7955cf04", tick, amount)
}

func (c *Chain) BuildSrc20Transfer(sender, receiver string, src20ContractAddress string, tick string, amount string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if !common.IsHexAddress(sender) || !common.IsHexAddress(receiver) ||
		!common.IsHexAddress(src20ContractAddress) {
		return nil, base.ErrInvalidAddress
	}
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, base.ErrInvalidAmount
	}
	data, err := encodeSrc20OpData(common.HexToAddress(receiver), "transfer", tick, amountInt)
	if err != nil {
		return nil, err
	}
	gasPrice, err := c.SuggestGasPrice()
	if err != nil {
		return nil, err
	}

	msg := NewCallMsg()
	msg.SetFrom(sender)
	msg.SetTo(src20ContractAddress)
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

func encodeSrc20OpData(receiver common.Address, op, tick string, amount *big.Int) ([]byte, error) {
	return AbiCoderEncode([]string{"address", "address", "string", "string", "string", "uint256", "uint256", "uint256", "uint256", "uint16", "string"},
		common.HexToAddress("0x0"),
		receiver,
		"src-20",
		op,
		tick,
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		amount,
		uint16(0),
		"{}",
	)
}
