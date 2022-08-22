package eth

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/ethereum/go-ethereum/common"
)

// ethRedPacketContract implement base.RedPacketContract interface
type ethRedPacketContract struct {
	chain   *Chain
	address string
}

func NewRedPacketContract(chain *Chain, contractAddress string) base.RedPacketContract {
	return &ethRedPacketContract{chain: chain, address: contractAddress}
}

func (contract *ethRedPacketContract) EstimateFee(rpa *base.RedPacketAction) (string, error) {
	switch rpa.Method {
	case RPAMethodCreate:
		count := rpa.CreateParams.Count
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
		return feeInt.String(), nil
	default:
		return "0", nil
	}
}

func (contract *ethRedPacketContract) FetchRedPacketCreationDetail(hash string) (*base.RedPacketDetail, error) {
	detail, err := contract.chain.FetchRedPacketCreationDetail(hash)
	if err != nil {
		return nil, err
	}
	return &base.RedPacketDetail{
		TransactionDetail: detail.TransactionDetail,
		AmountName:        detail.AmountName,
		AmountDecimal:     detail.AmountDecimal,
	}, nil
}

func (contract *ethRedPacketContract) SendTransaction(account base.Account, rpa *base.RedPacketAction) (string, error) {
	params, err := contract.packParams(rpa)
	if err != nil {
		return "", err
	}
	data, err := EncodeContractData(RedPacketABI, rpa.Method, params...)
	if err != nil {
		return "", err
	}
	gasPrice, err := contract.chain.SuggestGasPrice()
	if err != nil {
		return "", err
	}

	value, err := contract.EstimateFee(rpa)
	if err != nil {
		return "", err
	}
	msg := NewCallMsg()
	msg.SetFrom(account.Address())
	msg.SetTo(contract.address)
	msg.SetGasPrice(gasPrice.Value)
	msg.SetData(data)
	msg.SetValue(value)

	gasLimit, err := contract.chain.EstimateGasLimit(msg)
	if err != nil {
		gasLimit = &base.OptionalString{Value: "200000"}
		err = nil
	}
	msg.SetGasLimit(gasLimit.Value)
	tx := msg.TransferToTransaction()
	privateKeyHex, err := account.PrivateKeyHex()
	if err != nil {
		return "", err
	}
	signedTx, err := contract.chain.SignTransaction(privateKeyHex, tx)
	if err != nil {
		return "", err
	}
	return contract.chain.SendRawTransaction(signedTx.Value)
}

func (contract *ethRedPacketContract) packParams(rpa *base.RedPacketAction) ([]interface{}, error) {
	switch rpa.Method {
	case base.RPAMethodCreate:
		if rpa.CreateParams == nil {
			return nil, errors.New("invalid create params")
		}
		addr := common.HexToAddress(rpa.CreateParams.TokenAddress)
		c := big.NewInt(int64(rpa.CreateParams.Count))
		a, ok := big.NewInt(0).SetString(rpa.CreateParams.Amount, 10)
		if !ok {
			return nil, fmt.Errorf("Invalid red packet amount %v", rpa.CreateParams.Amount)
		}
		return []interface{}{addr, c, a}, nil
	case base.RPAMethodOpen:
		if rpa.OpenParams == nil {
			return nil, errors.New("invalid open params")
		}
		id := big.NewInt(rpa.OpenParams.PacketId)
		if len(rpa.OpenParams.Addresses) != len(rpa.OpenParams.Amounts) {
			return nil, fmt.Errorf("The number of opened addresses is not the same as the amount")
		}
		addrs := make([]common.Address, len(rpa.OpenParams.Addresses))
		for index, address := range rpa.OpenParams.Addresses {
			addrs[index] = common.HexToAddress(address)
		}
		amountInts := make([]*big.Int, len(rpa.OpenParams.Amounts))
		for index, amount := range rpa.OpenParams.Amounts {
			aInt, ok := big.NewInt(0).SetString(amount, 10)
			if !ok {
				return nil, fmt.Errorf("Invalid red packet amount %v", amount)
			}
			amountInts[index] = aInt
		}
		return []interface{}{id, addrs, amountInts}, nil
	case base.RPAMethodClose:
		if rpa.CloseParams == nil {
			return nil, errors.New("invalid close params")
		}
		id := big.NewInt(rpa.CloseParams.PacketId)
		addr := common.HexToAddress(rpa.CloseParams.Creator)
		return []interface{}{id, addr}, nil
	default:
		return nil, errors.New("invalid method")
	}
}
