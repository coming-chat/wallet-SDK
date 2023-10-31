package solanaswap

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
)

type WhirlpoolRewardInfoData struct {
	Mint      common.PublicKey
	Vault     common.PublicKey
	Authority common.PublicKey

	EmissionsPerSecondX64 *big.Int // u128
	GrowthGlobalX64       *big.Int // u128
}

type WhirlpoolData struct {
	WhirlpoolsConfig common.PublicKey
	WhirlpoolBump    []byte
	TickSpacing      uint16
	TickSpacingSeed  []byte
	FeeRate          uint16
	ProtocolFeeRate  uint16
	Liquidity        *big.Int
	SqrtPrice        *big.Int
	TickCurrentIndex int32
	ProtocolFeeOwedA uint64
	ProtocolFeeOwedB uint64
	TokenMintA       common.PublicKey
	TokenVaultA      common.PublicKey
	FeeGrowthGlobalA *big.Int
	TokenMintB       common.PublicKey
	TokenVaultB      common.PublicKey
	FeeGrowthGlobalB *big.Int

	RewardLastUpdatedTimestamp uint64
	RewardInfos                []WhirlpoolRewardInfoData
}

func (d *WhirlpoolData) Deserializer(data []byte) error {
	if len(data) < 8+261+128 {
		return errors.New("data length not enough")
	}
	ds := NewSolanaDeserializer(data)
	_ = ds.TakeBytes(8) // ignore

	d.WhirlpoolsConfig = ds.TakePublicKey()
	d.WhirlpoolBump = ds.TakeBytes(1)
	d.TickSpacing = ds.TakeU16()
	d.TickSpacingSeed = ds.TakeBytes(2)
	d.FeeRate = ds.TakeU16()
	d.ProtocolFeeRate = ds.TakeU16()
	d.Liquidity = ds.TakeU128()
	d.SqrtPrice = ds.TakeU128()
	d.TickCurrentIndex = ds.TakeI32()
	d.ProtocolFeeOwedA = ds.TakeU64()
	d.ProtocolFeeOwedB = ds.TakeU64()
	d.TokenMintA = ds.TakePublicKey()
	d.TokenVaultA = ds.TakePublicKey()
	d.FeeGrowthGlobalA = ds.TakeU128()
	d.TokenMintB = ds.TakePublicKey()
	d.TokenVaultB = ds.TakePublicKey()
	d.FeeGrowthGlobalB = ds.TakeU128()
	d.RewardLastUpdatedTimestamp = ds.TakeU64()

	infos := make([]WhirlpoolRewardInfoData, 0, ds.EndLength()/128)
	for ds.EndLength() >= 128 {
		info := WhirlpoolRewardInfoData{
			Mint:                  ds.TakePublicKey(),
			Vault:                 ds.TakePublicKey(),
			Authority:             ds.TakePublicKey(),
			EmissionsPerSecondX64: ds.TakeU128(),
			GrowthGlobalX64:       ds.TakeU128(),
		}
		infos = append(infos, info)
	}
	d.RewardInfos = infos

	return nil
}

// PDA: Program Derive Address
func getWhirlpoolPDA(programId string, whirlpoolConfigKey, tokenMintAKey, tokenMintBKey string, tickSpacing uint16) (common.PublicKey, error) {
	const PDA_WHIRLPOOL_SEED = "whirlpool"
	whirlpoolPub := common.PublicKeyFromString(whirlpoolConfigKey)
	aPub := common.PublicKeyFromString(tokenMintAKey)
	bPub := common.PublicKeyFromString(tokenMintBKey)
	if bytes.Compare(aPub[:], bPub[:]) == 1 {
		temp := aPub
		aPub = bPub
		bPub = temp
	}
	return findProgramAddressSeed(programId, [][]byte{
		[]byte(PDA_WHIRLPOOL_SEED),
		whirlpoolPub[:],
		aPub[:],
		bPub[:],
		binary.LittleEndian.AppendUint16(make([]byte, 0), tickSpacing),
	})
}

func getPoolData(cli *client.Client, poolAddr string) (*WhirlpoolData, error) {
	datas, err := getPoolsData(cli, []string{poolAddr})
	if err != nil {
		return nil, err
	}
	data := datas[poolAddr]
	if data != nil {
		return data, nil
	} else {
		return nil, errors.New("whirlpool not found")
	}
}

func getPoolsData(cli *client.Client, addresses []string) (map[string]*WhirlpoolData, error) {
	infos, err := cli.GetMultipleAccounts(context.Background(), addresses)
	if err != nil {
		return nil, err
	}
	res := make(map[string]*WhirlpoolData, 0)
	for idx, info := range infos {
		if info.Data == nil {
			continue
		}
		key := addresses[idx]
		var value WhirlpoolData
		err = value.Deserializer(info.Data)
		if err == nil {
			res[key] = &value
		}
	}
	return res, nil
}

type SearchedWhirlpoolData struct {
	Tick1   *WhirlpoolData
	Tick8   *WhirlpoolData
	Tick64  *WhirlpoolData
	Tick128 *WhirlpoolData
}

func SearchWhirlpool(cli *client.Client, programId string, whirlpoolConfigKey, mintA, mintB string) (*SearchedWhirlpoolData, error) {
	address1, err := getWhirlpoolPDA(programId, whirlpoolConfigKey, mintA, mintB, 1)
	if err != nil {
		return nil, err
	}
	address8, err := getWhirlpoolPDA(programId, whirlpoolConfigKey, mintA, mintB, 8)
	if err != nil {
		return nil, err
	}
	address64, err := getWhirlpoolPDA(programId, whirlpoolConfigKey, mintA, mintB, 64)
	if err != nil {
		return nil, err
	}
	address128, err := getWhirlpoolPDA(programId, whirlpoolConfigKey, mintA, mintB, 128)
	if err != nil {
		return nil, err
	}
	dataMap, err := getPoolsData(cli, []string{
		address1.ToBase58(),
		address8.ToBase58(),
		address64.ToBase58(),
		address128.ToBase58(),
	})
	if err != nil {
		return nil, err
	}
	return &SearchedWhirlpoolData{
		Tick1:   dataMap[address1.ToBase58()],
		Tick8:   dataMap[address8.ToBase58()],
		Tick64:  dataMap[address64.ToBase58()],
		Tick128: dataMap[address128.ToBase58()],
	}, nil
}
