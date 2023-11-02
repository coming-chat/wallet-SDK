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
	ds := NewSolanaDeserializer(data)
	if !ds.IsWhirlpoolDataType("Whirlpool") {
		return errors.New("invalid whirlpool data")
	}
	ds.StartWhirlpoolDataParse()

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
	Address common.PublicKey
	Data    WhirlpoolData
}

type SearchedWhirlpoolList struct {
	Tick1   *SearchedWhirlpoolData
	Tick8   *SearchedWhirlpoolData
	Tick64  *SearchedWhirlpoolData
	Tick128 *SearchedWhirlpoolData
}

func SearchWhirlpool(cli *client.Client, programId string, whirlpoolConfigKey, mintA, mintB string) (*SearchedWhirlpoolList, error) {
	addresses := []string{}
	keys := map[string]*common.PublicKey{}
	putSearchKey := func(key string, tickSpacing uint16) {
		if addr, err := getWhirlpoolPDA(programId, whirlpoolConfigKey, mintA, mintB, tickSpacing); err == nil {
			addresses = append(addresses, addr.ToBase58())
			keys[key] = &addr
		}
	}
	putSearchKey("tick1", 1)
	putSearchKey("tick8", 8)
	putSearchKey("tick64", 64)
	putSearchKey("tick128", 128)
	dataMap, err := getPoolsData(cli, addresses)
	if err != nil {
		return nil, err
	}
	takeWhirlpoolData := func(key string) *SearchedWhirlpoolData {
		if addr := keys[key]; addr != nil {
			if data := dataMap[addr.ToBase58()]; data != nil {
				return &SearchedWhirlpoolData{
					Address: *addr,
					Data:    *data,
				}
			}
		}
		return nil
	}
	return &SearchedWhirlpoolList{
		Tick1:   takeWhirlpoolData("tick1"),
		Tick8:   takeWhirlpoolData("tick8"),
		Tick64:  takeWhirlpoolData("tick64"),
		Tick128: takeWhirlpoolData("tick128"),
	}, nil
}
