package eth

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (e *EthChain) CallContractConstant(out interface{}, contractAddress, abiStr, methodName string, opts *bind.CallOpts, params ...interface{}) error {
	parsedAbi, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return err
	}
	inputParams, err := parsedAbi.Pack(methodName, params...)
	if err != nil {
		return err
	}

	method, ok := parsedAbi.Methods[methodName]
	if !ok {
		return errors.New("method not found")
	}
	return e.CallContractConstantWithPayload(out, contractAddress, hex.EncodeToString(inputParams), method.Outputs, opts)
}

func (e *EthChain) CallContractConstantWithPayload(out interface{}, contractAddress, payload string, outputTypes abi.Arguments, opts *bind.CallOpts) error {
	if opts == nil {
		opts = new(bind.CallOpts)
	}

	contractAddressObj := common.HexToAddress(contractAddress)

	payload = strings.TrimPrefix(payload, "0x")
	payloadBuf, err := hex.DecodeString(payload)
	if err != nil {
		return err
	}
	var (
		msg    = ethereum.CallMsg{From: opts.From, To: &contractAddressObj, Data: payloadBuf}
		ctx    = opts.Context
		code   []byte
		output []byte
	)
	if ctx == nil {
		ctxTemp, cancel := context.WithTimeout(context.Background(), e.timeout)
		defer cancel()
		ctx = ctxTemp
	}
	if opts.Pending {
		pb := bind.PendingContractCaller(e.RemoteRpcClient)
		output, err = pb.PendingCallContract(ctx, msg)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = pb.PendingCodeAt(ctx, contractAddressObj); err != nil {
				return err
			} else if len(code) == 0 {
				return errors.New(bind.ErrNoCode.Error())
			}
		}
	} else {
		output, err = bind.ContractCaller(e.RemoteRpcClient).CallContract(ctx, msg, opts.BlockNumber)
		if err != nil {
			return err
		}
		if len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = bind.ContractCaller(e.RemoteRpcClient).CodeAt(ctx, contractAddressObj, opts.BlockNumber); err != nil {
				return err
			} else if len(code) == 0 {
				return errors.New(bind.ErrNoCode.Error())
			}
		}
	}
	err = e.UnpackParams(out, outputTypes, hex.EncodeToString(output))
	if err != nil {
		return err
	}
	return nil
}
