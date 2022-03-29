package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestValidAddress(t *testing.T) {
	// 检测地址需要用 common.IsHexAddress
	// 结论：地址检测只检测了字符串长度、是否 16 进制
	addresses := []string{
		"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c2", // 正确的
		"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c1", // 这个地址修改了末尾一个数，大概率应该是错误的，但是暂时无法检测

		// 错误的
		"0x7161ada3EA6e53E5652A45988DdfF1cE595E09c",
		"0x7161ada3EA6e53E5652A",
	}

	for _, a := range addresses {
		address := common.HexToAddress(a)
		t.Log(common.IsHexAddress(a), address)
	}
}
