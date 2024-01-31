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

// MARK - Src20Token

type Src20Token struct {
	chain *Chain
	Tick  string
	// current support bevm donut token, Default is "0xf414dF7d8260A8e1e007F72892Cf5F0A7955cf04"
	ContractAddress string
}

func NewSrc20Token(chain *Chain, tick string) *Src20Token {
	return &Src20Token{
		chain: chain,
		Tick:  tick,

		ContractAddress: "0xf414dF7d8260A8e1e007F72892Cf5F0A7955cf04",
	}
}

func (t *Src20Token) Chain() base.Chain {
	return t.chain
}

func (t *Src20Token) TokenInfo() (*base.TokenInfo, error) {
	return &base.TokenInfo{
		Name:    t.Tick,
		Symbol:  t.Tick,
		Decimal: 0,
	}, nil
}

func (t *Src20Token) BalanceOfAddress(address string) (*base.Balance, error) {
	tokens, err := fetchDonutInscriptions(address, t.Tick, "")
	if err != nil {
		return nil, err
	}
	if len(tokens) <= 0 {
		return base.NewBalance("0"), nil
	}
	return base.NewBalance(tokens[0].Amount), nil
}
func (t *Src20Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return t.BalanceOfAddress(publicKey)
}
func (t *Src20Token) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
}

func (t *Src20Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if !common.IsHexAddress(sender) || !common.IsHexAddress(receiver) ||
		!common.IsHexAddress(t.ContractAddress) {
		return nil, base.ErrInvalidAddress
	}
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, base.ErrInvalidAmount
	}
	data, err := encodeSrc20OpData(common.HexToAddress(receiver), "transfer", t.Tick, amountInt)
	if err != nil {
		return nil, err
	}
	gasPrice, err := t.chain.SuggestGasPrice()
	if err != nil {
		return nil, err
	}

	msg := NewCallMsg()
	msg.SetFrom(sender)
	msg.SetTo(t.ContractAddress)
	msg.SetValue("0")
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(data)

	gasLimit, err := t.chain.EstimateGasLimit(msg)
	if err != nil {
		return nil, err
	}
	msg.SetGasLimit(gasLimit.Value)

	return msg.TransferToTransaction(), nil
}

// Before invoking this method, it is best to check `CanTransferAll()`
func (t *Src20Token) CanTransferAll() bool {
	return false
}
func (t *Src20Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}

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
// - param graphURL: Default "https://bc.dnt.social/v1/common/search"
func FetchDonutInscriptions(owner string, graphURL string) (*DonutInscriptionArray, error) {
	arr, err := fetchDonutInscriptions(owner, "", graphURL)
	if err != nil {
		return nil, err
	}
	return &DonutInscriptionArray{AnyArray: arr}, nil
}

// fetchDonutInscriptions
// - param tick: if is empty, will query all of owner's ticks
// - param graphURL: Default "https://bc.dnt.social/v1/common/search"
func fetchDonutInscriptions(owner, tick, graphURL string) (arr []*DonutInscription, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if graphURL == "" {
		graphURL = "https://bc.dnt.social/v1/common/search"
	}
	funcStr := ""
	if tick == "" {
		funcStr = fmt.Sprintf(`src20Balances(holder: "%v", first: 100)`, owner)
	} else {
		funcStr = fmt.Sprintf(`src20Balances(holder: "%v", tick: "%v")`, owner, tick)
	}
	query := fmt.Sprintf(`{
		%v {
			edges{
				node{
					tick
					amount
				}
			}
		}
	}`, funcStr)

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
	inscriptions := make([]*DonutInscription, 0, len(out.Src20Balances.Edges))
	for _, node := range out.Src20Balances.Edges {
		if node.Node.Amount != "" && node.Node.Amount != "0" {
			inscriptions = append(inscriptions, node.Node)
		}
	}
	return inscriptions, nil
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
