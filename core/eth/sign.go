package eth

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (e *EthChain) SignMsg(privateKey string, data string) (string, error) {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKeyHex, err := hex.DecodeString(privateKey)
	if err != nil {
		return "", err
	}
	privateKeyObj, err := crypto.ToECDSA(privateKeyHex)
	if err != nil {
		return "", err
	}
	hash, err := e.SignHashForMsg(data)
	if err != nil {
		return "", err
	}
	hashBuf, _ := hex.DecodeString(hash)
	signature, err := crypto.Sign(hashBuf, privateKeyObj)
	if err != nil {
		return "", err
	}
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return hex.EncodeToString(signature), nil
}

// 以太坊的 hash 专门在数据前面加上了一段话
func (e *EthChain) SignHashForMsg(data string) (string, error) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return hex.EncodeToString(crypto.Keccak256([]byte(msg))), nil
}

func (e *EthChain) RecoverSignerAddressFromMsgHash(msgHash, sig string) (*common.Address, error) {

	sig = strings.TrimPrefix(sig, "0x")
	sigHex, err := hex.DecodeString(sig)
	if err != nil {
		return nil, err
	}

	msgHash = strings.TrimPrefix(msgHash, "0x")
	msgHashHex, err := hex.DecodeString(msgHash)
	if err != nil {
		return nil, err
	}

	if len(sigHex) != 65 {
		return nil, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigHex[64] != 27 && sigHex[64] != 28 {
		return nil, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigHex[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.Ecrecover(msgHashHex, sigHex)
	if err != nil {
		return nil, err
	}
	pubKey, err := crypto.UnmarshalPubkey(rpk)
	if err != nil {
		return nil, err
	}
	//pubKey := crypto.ToECDSAPub(rpk)
	//crypto.FromECDSAPub()
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return &recoveredAddr, nil
}

func (e *EthChain) RecoverSignerAddress(msg, sig string) (*common.Address, error) {
	hash, err := e.SignHashForMsg(msg)
	if err != nil {
		return nil, err
	}
	return e.RecoverSignerAddressFromMsgHash(hash, sig)
}
