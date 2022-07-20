package eth

import (
	"fmt"
	"math/big"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/common"
)

const RedPacketABI = `[{"inputs":[{"internalType":"address","name":"_admin","type":"address"},{"internalType":"address","name":"_beneficiary","type":"address"},{"internalType":"uint256","name":"_base_fee","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"_old","type":"address"},{"indexed":false,"internalType":"address","name":"_new","type":"address"}],"name":"AdminChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"_old","type":"address"},{"indexed":false,"internalType":"address","name":"_new","type":"address"}],"name":"BeneficiaryChanged","type":"event"},{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"address","name":"maybe_creator","type":"address"}],"name":"close","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contract IERC20","name":"token","type":"address"},{"internalType":"uint256","name":"count","type":"uint256"},{"internalType":"uint256","name":"total_balance","type":"uint256"}],"name":"create","outputs":[],"stateMutability":"payable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"_fee","type":"uint256"}],"name":"NewBasePrepaidFee","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"_id","type":"uint256"},{"indexed":false,"internalType":"contract IERC20","name":"_token","type":"address"},{"indexed":false,"internalType":"uint256","name":"_count","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_balance","type":"uint256"}],"name":"NewRedEnvelop","type":"event"},{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"},{"internalType":"address[]","name":"luck_accounts","type":"address[]"},{"internalType":"uint256[]","name":"balances","type":"uint256[]"}],"name":"open","outputs":[],"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"new_admin","type":"address"}],"name":"set_admin","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"new_beneficiary","type":"address"}],"name":"set_beneficiary","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"new_fee","type":"uint256"}],"name":"set_prepaid_fee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"_id","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_remain_count","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"_remain_balance","type":"uint256"}],"name":"UpdateRedEnvelop","type":"event"},{"inputs":[],"name":"admin","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"base_prepaid_fee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"beneficiary","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"count","type":"uint256"}],"name":"calc_prepaid_fee","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"id","type":"uint256"}],"name":"is_valid","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"max_count","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"next_id","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"red_envelop_infos","outputs":[{"internalType":"contract IERC20","name":"token","type":"address"},{"internalType":"uint256","name":"remain_count","type":"uint256"},{"internalType":"uint256","name":"remain_balance","type":"uint256"}],"stateMutability":"view","type":"function"}]`

const (
	RPAMethodCreate = "create"
	RPAMethodOpen   = "open"
	RPAMethodClose  = "close"
)

type RedPacketAction struct {
	Method string
	Params []interface{}
}

// 用户发红包 的操作
func NewRedPacketActionCreate(erc20TokenAddress string, count int, amount string) (*RedPacketAction, error) {
	addr := common.HexToAddress(erc20TokenAddress)
	c := big.NewInt(int64(count))
	a, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid red packet amount %v", amount)
	}
	return &RedPacketAction{
		Method: RPAMethodCreate,
		Params: []interface{}{addr, c, a},
	}, nil
}

// 批量打开红包 的操作
func NewRedPacketActionOpen(packetId int64, addresses []string, amounts []string) (*RedPacketAction, error) {
	id := big.NewInt(packetId)
	if len(addresses) != len(amounts) {
		return nil, fmt.Errorf("The number of opened addresses is not the same as the amount")
	}
	addrs := make([]common.Address, len(addresses))
	for index, address := range addresses {
		addrs[index] = common.HexToAddress(address)
	}
	amountInts := make([]*big.Int, len(amounts))
	for index, amount := range amounts {
		aInt, ok := big.NewInt(0).SetString(amount, 10)
		if !ok {
			return nil, fmt.Errorf("Invalid red packet amount %v", amount)
		}
		amountInts[index] = aInt
	}
	return &RedPacketAction{
		Method: RPAMethodOpen,
		Params: []interface{}{id, addrs, amountInts},
	}, nil
}

// 结束红包领取 的操作
func NewRedPacketActionClose(packetId int64, creator string) (*RedPacketAction, error) {
	id := big.NewInt(packetId)
	addr := common.HexToAddress(creator)
	return &RedPacketAction{
		Method: RPAMethodClose,
		Params: []interface{}{id, addr},
	}, nil
}

func (rpa *RedPacketAction) EstimateAmount() string {
	switch rpa.Method {
	case RPAMethodCreate:
		count := rpa.Params[1].(*big.Int).Int64()
		rate := 200.0
		switch {
		case count <= 10:
			rate = 4
		case count <= 100:
			rate = 16
		case count <= 1000:
			rate = 200
		}
		feeFloat := big.NewFloat(0.025 * rate)
		feeFloat.Mul(feeFloat, big.NewFloat(1e18))
		feeInt, _ := feeFloat.Int(big.NewInt(0))
		return feeInt.String()
	default:
		return "0"
	}
}

// 保证用户发 erc20 的红包时，红包合约可以有权限操作用户的资产
// @param account 要发红包的用户的账号，也许需要用到私钥来发起授权交易
// @param chain evm 链
// @param erc20Contract 要用作发红包的币种
// @param coins 如果需要发起新授权，指定要授权的币个数 default 10^6
// @return 如果授权成功，不会返回错误，如果有新授权，会返回授权交易的 hash
func (rpa *RedPacketAction) EnsureApprovedTokens(account *Account, chain *Chain, spender string, coins int) (string, error) {
	// only red packet **create** need approve
	if rpa.Method != RPAMethodCreate {
		return "", nil
	}

	erc20Token := rpa.Params[0].(common.Address).String()
	amount := rpa.Params[2].(*big.Int)

	token := NewErc20Token(chain, erc20Token)
	approved, err := token.Allowance(account.Address(), spender)
	if err != nil {
		return "", err
	}
	if approved.Cmp(amount) >= 0 {
		return "", nil
	}

	decimal, err := token.Decimal()
	if err != nil {
		decimal = 18
		err = nil
	}
	if coins <= 0 {
		coins = 1e6 // default 1 million coins
	}
	oneCoin := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	approveValue := big.NewInt(0).Mul(oneCoin, big.NewInt(int64(coins)))
	approveValue = base.MaxBigInt(approveValue, amount)

	return token.Approve(account, spender, approveValue)
}

// @param fromAddress 要调用红包业务的操作者
// @param contractAddress 红包合约地址
// @param chain 要发红包的链
func (rpa *RedPacketAction) TransactionFrom(fromAddress, contractAddress string, chain *Chain) (*Transaction, error) {
	data, err := EncodeAbiData(RedPacketABI, rpa.Method, rpa.Params...)
	if err != nil {
		return nil, err
	}
	gasPrice, err := chain.SuggestGasPrice()
	if err != nil {
		return nil, err
	}

	value := rpa.EstimateAmount()
	msg := NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(contractAddress)
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(data)
	msg.SetValue(value)

	gasLimit, err := chain.EstimateGasLimit(msg)
	if err != nil {
		gasLimit = &base.OptionalString{Value: "200000"}
		err = nil
	}
	msg.SetGasLimit(gasLimit.Value)

	return msg.TransferToTransaction(), nil
}
