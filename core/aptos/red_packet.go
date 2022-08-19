package aptos

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/coming-chat/go-aptos/aptostypes"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type aptosRedPacketContract struct {
	chain   *Chain
	address string
}

func NewRedPacketContract(chain *Chain, contractAddress string) base.RedPacketContract {
	return &aptosRedPacketContract{chain: chain, address: contractAddress}
}

func (contract *aptosRedPacketContract) PackTransaction(account base.Account, rpa *base.RedPacketAction) (*base.OptionalString, error) {
	fromAddress := account.Address()

	client, err := contract.chain.client()
	if err != nil {
		return nil, err
	}

	accountData, err := client.GetAccount(fromAddress)
	if err != nil {
		return nil, err
	}
	ledgerInfo, err := client.LedgerInfo()
	if err != nil {
		return nil, err
	}

	payload, err := contract.createPayload(rpa)
	if err != nil {
		return nil, err
	}

	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            contract.estimateGas(rpa),
		GasUnitPrice:            GasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // timeout 10 mins
	}

	signingMessage, err := client.CreateTransactionSigningMessage(transaction)
	if err != nil {
		return nil, err
	}
	signatureData, _ := account.Sign(signingMessage, "")

	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: account.PublicKeyHex(),
		Signature: types.HexEncodeToString(signatureData),
	}

	signedTransactionData, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	return &base.OptionalString{Value: types.HexEncodeToString(signedTransactionData)}, nil
}

func (contract *aptosRedPacketContract) estimateGas(rpa *base.RedPacketAction) uint64 {
	if rpa.Method != base.RPAMethodOpen {
		return MaxGasAmount
	}
	// devnet 测试结果
	// OpenPacketCount  |  GasUsed
	// 10 | 19
	// 100 | 249
	// 200 | 483
	// 500 | 1158
	// 1000 | 2355
	if len(rpa.OpenParams.Addresses) > 200 {
		return 2000
	} else if len(rpa.OpenParams.Addresses) > 500 {
		return 3000
	} else {
		return MaxGasAmount
	}
}

func (contract *aptosRedPacketContract) createPayload(rpa *base.RedPacketAction) (*aptostypes.Payload, error) {
	switch rpa.Method {
	case base.RPAMethodCreate:
		if nil == rpa.CreateParams {
			return nil, fmt.Errorf("create params is nil")
		}
		return &aptostypes.Payload{
			Type:          "script_function_payload",
			Function:      contract.address + "::red_packet::create",
			TypeArguments: []string{},
			Arguments: []interface{}{
				strconv.FormatInt(int64(rpa.CreateParams.Count), 10),
				rpa.CreateParams.Amount,
			},
		}, nil
	case base.RPAMethodOpen:
		if nil == rpa.OpenParams {
			return nil, fmt.Errorf("open params is nil")
		}
		return &aptostypes.Payload{
			Type:          "script_function_payload",
			Function:      contract.address + "::red_packet::open",
			TypeArguments: []string{},
			Arguments: []interface{}{
				strconv.FormatInt(int64(rpa.OpenParams.PacketId), 10),
				rpa.OpenParams.Addresses,
				rpa.OpenParams.Amounts,
			},
		}, nil
	case base.RPAMethodClose:
		if nil == rpa.CloseParams {
			return nil, fmt.Errorf("close params is nil")
		}
		return &aptostypes.Payload{
			Type:          "script_function_payload",
			Function:      contract.address + "::red_packet::close",
			TypeArguments: []string{},
			Arguments: []interface{}{
				strconv.FormatInt(int64(rpa.CloseParams.PacketId), 10),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsopported red packet method %s", rpa.Method)
	}
}
