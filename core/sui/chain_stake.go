package sui

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"sync"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

var cachedSuiSystemState *types.SuiSystemState
var cachedDelegatedStakesMap sync.Map

const maxGasBudgetForStake = 20000

type ValidatorState struct {
	// The current epoch in Sui. An epoch takes approximately 24 hours and runs in checkpoints.
	Epoch int64 `json:"epoch"`
	// Array of `Validator` elements
	Validators *base.AnyArray `json:"validators"`

	// The amount of all tokens staked in the Sui Network.
	TotalStaked string `json:"totalStaked"`
	// The amount of rewards won by all Sui validators in the last epoch.
	TotalRewards string `json:"lastEpochReward"`

	Time int64 `json:"time"`
}

func (s *ValidatorState) JsonString() (*base.OptionalString, error) {
	return base.JsonString(s)
}

func NewValidatorStateWithJsonString(str string) (*ValidatorState, error) {
	var o ValidatorState
	err := base.FromJsonString(str, &o)
	return &o, err
}

type Validator struct {
	Address    string  `json:"address"`
	Name       string  `json:"name"`
	Desc       string  `json:"desc"`
	ImageUrl   string  `json:"imageUrl"`
	ProjectUrl string  `json:"projectUrl"`
	APY        float64 `json:"apy"`

	Commission      int64  `json:"commission"`
	TotalStaked     string `json:"totalStaked"`
	DelegatedStaked string `json:"delegatedStaked"`
	SelfStaked      string `json:"selfStaked"`
	TotalRewards    string `json:"totalRewards"`
	GasPrice        int64  `json:"gasPrice"`
}

func (s *Validator) JsonString() (*base.OptionalString, error) {
	return base.JsonString(s)
}

func NewValidatorWithJsonString(str string) (*Validator, error) {
	var o Validator
	err := base.FromJsonString(str, &o)
	return &o, err
}

func (o *Validator) AsAny() *base.Any {
	return &base.Any{Value: o}
}
func AsValidator(a *base.Any) *Validator {
	if r, ok := a.Value.(*Validator); ok {
		return r
	}
	if r, ok := a.Value.(Validator); ok {
		return &r
	}
	return nil
}

type DelegationStatus = base.SDKEnumInt

const (
	DelegationStatusPending DelegationStatus = 0
	DelegationStatusActived DelegationStatus = 1
)

type DelegatedStake struct {
	StakeId          string `json:"stakeId"`
	ValidatorAddress string `json:"validatorAddress"`
	Principal        string `json:"principal"`
	RequestEpoch     int64  `json:"requestEpoch"`

	Status       DelegationStatus `json:"status"`
	DelegationId string           `json:"delegationId"`
	EarnedAmount string           `json:"earnedAmount"`
	Validator    *Validator       `json:"validator"`
}

func (s *DelegatedStake) JsonString() (*base.OptionalString, error) {
	return base.JsonString(s)
}

func NewDelegatedStakeWithJsonString(str string) (*DelegatedStake, error) {
	var o DelegatedStake
	err := base.FromJsonString(str, &o)
	return &o, err
}

func NewDelegatedStakeArrayWithJsonString(str string) (*base.AnyArray, error) {
	var o []*DelegatedStake
	err := base.FromJsonString(str, &o)
	arr := make([]any, len(o))
	for i, v := range o {
		arr[i] = v
	}
	return &base.AnyArray{Values: arr}, err
}

func (o *DelegatedStake) AsAny() *base.Any {
	return &base.Any{Value: o}
}
func AsDelegatedStake(a *base.Any) *DelegatedStake {
	if r, ok := a.Value.(*DelegatedStake); ok {
		return r
	}
	if r, ok := a.Value.(DelegatedStake); ok {
		return &r
	}
	return nil
}

func (c *Chain) GetValidatorState() (s *ValidatorState, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.client()
	if err != nil {
		return nil, err
	}
	state, err := cli.GetSuiSystemState(context.Background())
	if err != nil {
		return nil, err
	}
	cachedSuiSystemState = state // cache

	totalRewards := big.NewInt(0)
	var validators = &base.AnyArray{}
	for _, v := range state.Validators.ActiveValidators {
		validator := mapRawValidator(&v, state.Epoch)
		validators.Values = append(validators.Values, validator)
		reward := big.NewInt(int64(v.DelegationStakingPool.RewardsPool.Value))
		totalRewards.Add(totalRewards, reward)
	}

	delegationStake := big.NewInt(int64(state.Validators.DelegationStake))
	validatorStake := big.NewInt(int64(state.Validators.ValidatorStake))
	totalStake := big.NewInt(0).Add(delegationStake, validatorStake)

	res := &ValidatorState{
		Epoch:      int64(state.Epoch),
		Validators: validators,

		TotalStaked:  totalStake.String(),
		TotalRewards: totalRewards.String(),
	}

	return res, nil
}

// @return Array of `DelegatedStake` elements
func (c *Chain) GetDelegatedStakes(owner string) (arr *base.AnyArray, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	addr, err := types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}
	cli, err := c.client()
	if err != nil {
		return nil, err
	}
	list, err := cli.GetDelegatedStakes(context.Background(), *addr)
	if err != nil {
		return nil, err
	}
	if cachedSuiSystemState == nil {
		cachedSuiSystemState, _ = cli.GetSuiSystemState(context.Background())
	}

	var stakes = &base.AnyArray{}
	for _, s := range list {
		stake := mapRawStake(&s)
		stakes.Values = append(stakes.Values, stake)
	}

	cachedDelegatedStakesMap.Store(owner, stakes)
	return stakes, nil
}

// @useCache If true, when there is cached data, the result will be returned directly without requesting data on the chain.
func (c *Chain) TotalStakedSuiAtValidator(validator, owner string, useCache bool) (sui *base.OptionalString, err error) {
	var stakes *base.AnyArray
	if useCache {
		if cachedStakes, ok := cachedDelegatedStakesMap.Load(owner); ok {
			if v, ok := cachedStakes.(*base.AnyArray); ok {
				stakes = v
			}
		}
	}
	if stakes == nil {
		stakes, err = c.GetDelegatedStakes(owner)
		if err != nil {
			return nil, err
		}
	}

	total := big.NewInt(0)
	for _, val := range stakes.Values {
		if stake, ok := val.(*DelegatedStake); ok {
			if principalInt, ok := big.NewInt(0).SetString(stake.Principal, 10); ok {
				total = total.Add(total, principalInt)
			}
		}
	}
	return &base.OptionalString{Value: total.String()}, nil
}

func AverageApyOfDelegatedStakes(stakes *base.AnyArray) float64 {
	totalApy := float64(0)
	added := map[string]bool{}
	for _, val := range stakes.Values {
		if stake, ok := val.(*DelegatedStake); ok {
			if added[stake.ValidatorAddress] != true {
				totalApy = totalApy + stake.Validator.APY
				added[stake.ValidatorAddress] = true
				continue
			}
		}
	}

	if count := len(added); count > 0 {
		return totalApy / float64(count)
	} else {
		return 0
	}
}

func (c *Chain) AddDelegation(owner, amount string, validatorAddress string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	validator, err := types.NewAddressFromHex(validatorAddress)
	if err != nil {
		return
	}
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, errors.New("invalid stake amount")
	}
	cli, err := c.client()
	if err != nil {
		return
	}
	gasPrice, err := cli.GetReferenceGasPrice(context.Background())
	if err != nil {
		return
	}
	maxGasFee := gasPrice * maxGasBudgetForStake
	allCoins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	if err != nil {
		return
	}
	needCoins, gasCoin, err := allCoins.PickSUICoinsWithGas(amountInt, maxGasFee, types.PickBigger)
	if err != nil {
		return
	}
	coinIds := []types.ObjectId{}
	for _, coin := range needCoins {
		coinIds = append(coinIds, coin.Reference.ObjectId)
	}
	gasId := gasCoin.Reference.ObjectId
	txBytes, err := cli.RequestAddDelegation(context.Background(), *signer, coinIds, amountInt.Uint64(), *validator, gasId, maxGasBudgetForStake)
	if err != nil {
		return
	}
	return &Transaction{
		Txn:          *txBytes,
		MaxGasBudget: int64(maxGasFee),
	}, nil
}

func (c *Chain) WithdrawDelegation(owner, delegationId, stakeId string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	delegation, err := types.NewHexData(delegationId)
	if err != nil {
		return
	}
	stakeSui, err := types.NewHexData(stakeId)
	if err != nil {
		return
	}
	cli, err := c.client()
	if err != nil {
		return
	}
	gasPrice, err := cli.GetReferenceGasPrice(context.Background())
	if err != nil {
		return
	}
	maxGasFee := gasPrice * maxGasBudgetForStake
	allCoins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	if err != nil {
		return
	}
	gasCoin, err := allCoins.PickCoinNoLess(maxGasFee)
	if err != nil {
		return
	}
	gasId := gasCoin.Reference.ObjectId
	txnBytes, err := cli.RequestWithdrawDelegation(context.Background(), *signer, *delegation, *stakeSui, gasId, maxGasBudgetForStake)
	if err != nil {
		return
	}

	return &Transaction{
		Txn:          *txnBytes,
		MaxGasBudget: int64(maxGasFee),
	}, nil
}

func mapRawValidator(v *types.Validator, epoch uint64) *Validator {
	if v == nil {
		return nil
	}
	meta := v.Metadata
	selfStaked := v.StakeAmount
	delegatedStaked := v.DelegationStakingPool.SuiBalance
	totalStaked := selfStaked + delegatedStaked
	rewardsPoolBalance := v.DelegationStakingPool.RewardsPool.Value

	validator := Validator{
		Address:    meta.SuiAddress.String(),
		Name:       string(meta.Name),
		Desc:       string(meta.Description),
		ImageUrl:   string(meta.ImageUrl),
		ProjectUrl: string(meta.ProjectUrl),
		APY:        v.CalculateAPY(epoch),

		Commission:      int64(v.CommissionRate),
		SelfStaked:      strconv.FormatUint(selfStaked, 10),
		DelegatedStaked: strconv.FormatUint(delegatedStaked, 10),
		TotalStaked:     strconv.FormatUint(totalStaked, 10),
		TotalRewards:    strconv.FormatUint(rewardsPoolBalance, 10),
		GasPrice:        int64(v.GasPrice),
	}
	return &validator
}

func mapRawStake(s *types.DelegatedStake) *DelegatedStake {
	if s == nil {
		return nil
	}
	stake := &DelegatedStake{
		StakeId:          s.StakedSui.Id.Id.String(),
		ValidatorAddress: s.StakedSui.ValidatorAddress.String(),
		Principal:        strconv.FormatUint(s.StakedSui.Principal.Value, 10),
		RequestEpoch:     int64(s.StakedSui.DelegationRequestEpoch),
	}

	if cachedSuiSystemState != nil {
		earned, validator := s.CalculateEarnAmount(cachedSuiSystemState.Validators.ActiveValidators)
		stake.EarnedAmount = strconv.FormatUint(earned, 10)
		stake.Validator = mapRawValidator(validator, cachedSuiSystemState.Epoch)
	}

	statusData, err := json.Marshal(s.DelegationStatus)
	if err != nil {
		stake.Status = DelegationStatusPending
		return stake
	}
	var activeStatus types.ActiveDelegationStatus
	err = json.Unmarshal(statusData, &activeStatus)
	if err != nil {
		stake.Status = DelegationStatusPending
		return stake
	}

	stake.Status = DelegationStatusActived
	stake.DelegationId = activeStatus.Active.Id.Id.String()
	return stake
}
