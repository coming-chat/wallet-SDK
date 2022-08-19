package base

import (
	"encoding/json"
	"fmt"
	"math/big"
)

const (
	RPAMethodCreate = "create"
	RPAMethodOpen   = "open"
	RPAMethodClose  = "close"
)

type RedPacketContract interface {
	PackTransaction(Account, *RedPacketAction) (*OptionalString, error)
	FetchRedPacketCreationDetail(hash string) (*RedPacketDetail, error)
	EstimateFee(*RedPacketAction) (string, error) // create red packet fee
}

type RedPacketAction struct {
	Method string

	CreateParams *RedPacketCreateParams
	OpenParams   *RedPacketOpenParams
	CloseParams  *RedPacketCloseParams
}

type RedPacketCreateParams struct {
	TokenAddress string // erc20 tokenAddress, aptos coin type
	Count        int
	Amount       string
}
type RedPacketOpenParams struct {
	PacketId  int64
	Addresses []string
	Amounts   []string
}
type RedPacketCloseParams struct {
	PacketId int64
	Creator  string
}

type RedPacketDetail struct {
	*TransactionDetail

	AmountName    string
	AmountDecimal int16
}

// 用户发红包 的操作
func NewRedPacketActionCreate(tokenAddress string, count int, amount string) (*RedPacketAction, error) {
	_, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid red packet amount %v", amount)
	}
	return &RedPacketAction{
		Method: RPAMethodCreate,
		CreateParams: &RedPacketCreateParams{
			TokenAddress: tokenAddress,
			Count:        count,
			Amount:       amount,
		},
	}, nil
}

// 批量打开红包 的操作
func NewRedPacketActionOpen(packetId int64, addresses []string, amounts []string) (*RedPacketAction, error) {
	if len(addresses) != len(amounts) {
		return nil, fmt.Errorf("the number of opened addresses is not the same as the amount")
	}
	for _, amount := range amounts {
		_, ok := big.NewInt(0).SetString(amount, 10)
		if !ok {
			return nil, fmt.Errorf("invalid red packet amount %v", amount)
		}
	}
	return &RedPacketAction{
		Method: RPAMethodOpen,
		OpenParams: &RedPacketOpenParams{
			PacketId:  packetId,
			Addresses: addresses,
			Amounts:   amounts,
		},
	}, nil
}

// 结束红包领取 的操作
func NewRedPacketActionClose(packetId int64, creator string) (*RedPacketAction, error) {
	return &RedPacketAction{
		Method: RPAMethodClose,
		CloseParams: &RedPacketCloseParams{
			PacketId: packetId,
			Creator:  creator,
		},
	}, nil
}

func (rpa *RedPacketAction) Do(chain Chain, account Account, contract RedPacketContract, password string) (string, error) {
	signedTx, err := contract.PackTransaction(account, rpa)
	if err != nil {
		return "", err
	}
	return chain.SendRawTransaction(signedTx.Value)
}

func (d *RedPacketDetail) JsonString() string {
	bytes, err := json.Marshal(d)
	if err != nil {
		return ""
	}
	return string(bytes)
}
