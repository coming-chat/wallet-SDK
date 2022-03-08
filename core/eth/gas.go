package eth

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/itering/scale.go/pkg/go-ethereum/crypto/sha3"
)

// 获取标准gas价格
func (e *EthChain) SuggestGasPrice() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasPrice, err := e.RemoteRpcClient.SuggestGasPrice(ctx)

	if err != nil {
		return "0", nil
	}
	return gasPrice.String(), err
}

// erc20 代币 Transfer，Approve GasLimit 估计
func (e *EthChain) EstimateErc20GasLimit(toAddress string, amount string) (string, error) {

	toAddressHex := common.HexToAddress(toAddress)
	amountBigInt, _ := new(big.Int).SetString(amount, 10)

	// 我们需要找出我们将要调用的智能合约函数名，以及函数将接收的输入。
	// 然后我们使用函数名的keccak-256哈希来检索 方法ID，它是前8个字符（4个字节）。
	// 然后，我们附加我们发送的地址，并附加我们打算转账的代币数量。 这些输入需要256位长（32字节）并填充左侧。
	// 方法ID不需填充
	transferFnSignature := []byte("transfer(address,uint256)")
	hash := sha3.NewKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddressHex.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amountBigInt.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasLimit, err := e.RemoteRpcClient.EstimateGas(ctx, ethereum.CallMsg{
		To:   &toAddressHex,
		Data: data,
	})
	if err != nil {
		return "0", err
	}
	gasLimitStr := strconv.FormatUint(gasLimit, 10)
	return gasLimitStr, nil
}

func (e *EthChain) EstimateGasLimit(msg ethereum.CallMsg) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()
	gasCount, err := e.RemoteRpcClient.EstimateGas(ctx, msg)
	if err != nil {
		return "0", err
	}
	gasLimitStr := strconv.FormatUint(gasCount, 10)
	return gasLimitStr, nil
}
