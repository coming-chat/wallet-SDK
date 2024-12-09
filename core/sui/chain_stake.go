package sui

import (
	"context"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/coming-chat/go-sui/v2/client"
	"github.com/coming-chat/go-sui/v2/sui_types"
	"github.com/coming-chat/go-sui/v2/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
	"github.com/shopspring/decimal"
)

var cachedSuiSystemState *types.SuiSystemStateSummary
var cachedDelegatedStakesMap sync.Map

const maxGasBudgetForStake = 12000000

type ValidatorState struct {
	// The current epoch in Sui. An epoch takes approximately 24 hours and runs in checkpoints.
	Epoch                 int64 `json:"epoch"`
	EpochStartTimestampMs int64 `json:"epochStartTimestampMs"`
	EpochDurationMs       int64 `json:"epochDurationMs"`

	// Array of `Validator` elements
	Validators *ValidatorArray `json:"validators"`

	// The amount of all tokens staked in the Sui Network.
	TotalStaked string `json:"totalStaked"`
	// The amount of rewards won by all Sui validators in the last epoch.
	TotalRewards string `json:"lastEpochReward"`
}

func (s *ValidatorState) JsonString() (*base.OptionalString, error) {
	return base.JsonString(s)
}

// @return if time > 0 indicates how long it will take to get the reward;
// if time < 0 indicates how much time has passed since the reward was earned;
func (s *ValidatorState) EarningAmountTimeAfterNowMs() int64 {
	return s.EarningAmountTimeAfterTimestampMs(time.Now().UnixMilli())
}

func (s *ValidatorState) EarningAmountTimeAfterTimestampMs(timestamp int64) int64 {
	ranTime := timestamp - s.EpochStartTimestampMs
	leftTime := s.EpochDurationMs*2 - ranTime
	return leftTime
}

func NewValidatorState() *ValidatorState {
	return &ValidatorState{}
}

func NewValidatorStateWithJsonString(str string) (*ValidatorState, error) {
	var o ValidatorState
	err := base.FromJsonString(str, &o)
	return &o, err
}

type Validator struct {
	APY        float64 `json:"apy"`
	PoolShare  float64 `json:"poolShare"`
	Commission int64   `json:"commission"`
	GasPrice   int64   `json:"gasPrice"`

	Address    string `json:"address"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	ImageUrl   string `json:"imageUrl"`
	ProjectUrl string `json:"projectUrl"`

	TotalStaked     string `json:"totalStaked"`
	DelegatedStaked string `json:"delegatedStaked"`
	SelfStaked      string `json:"selfStaked"`
	TotalRewards    string `json:"totalRewards"`
}

type ValidatorArray struct {
	inter.AnyArray[*Validator]
}

func (v *Validator) isSpecified() bool {
	if v == nil {
		return false
	}
	reg := regexp.MustCompile(`(?i)^Coming[ ._-]*Chat$`)
	nameMatched := reg.MatchString(v.Name)
	return nameMatched || v.Address == "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
}

func (s *Validator) JsonString() (*base.OptionalString, error) {
	return base.JsonString(s)
}

func NewValidator() *Validator {
	return &Validator{}
}

func NewValidatorWithJsonString(str string) (*Validator, error) {
	var o Validator
	err := base.FromJsonString(str, &o)
	return &o, err
}

type DelegationStatus = base.SDKEnumInt

const (
	DelegationStatusPending DelegationStatus = 0
	DelegationStatusActived DelegationStatus = 1
)

type DelegatedStake struct {
	RequestEpoch     int64  `json:"requestEpoch"`
	StakeId          string `json:"stakeId"`
	ValidatorAddress string `json:"validatorAddress"`
	Principal        string `json:"principal"`

	Status       DelegationStatus `json:"status"`
	DelegationId string           `json:"delegationId"`
	EarnedAmount string           `json:"earnedAmount"`
	Validator    *Validator       `json:"validator"`
}

type DelegatedStakeArray struct {
	inter.AnyArray[*DelegatedStake]
}

// @return if time > 0 indicates how long it will take to get the reward;
// if time < 0 indicates how much time has passed since the reward was earned;
func (s *DelegatedStake) EarningAmountTimeAfterNowMs(stateInfo *ValidatorState) int64 {
	return s.EarningAmountTimeAfterTimestampMs(time.Now().UnixMilli(), stateInfo)
}

func (s *DelegatedStake) EarningAmountTimeAfterTimestampMs(timestamp int64, stateInfo *ValidatorState) int64 {
	rewardEpoch := s.RequestEpoch + 2
	leftEpoch := rewardEpoch - stateInfo.Epoch

	ranTime := timestamp - stateInfo.EpochStartTimestampMs
	leftTime := stateInfo.EpochDurationMs*leftEpoch - ranTime
	return leftTime
}

func (s *DelegatedStake) JsonString() (*base.OptionalString, error) {
	return base.JsonString(s)
}

func NewDelegatedStake() *DelegatedStake {
	return &DelegatedStake{}
}

func NewDelegatedStakeWithJsonString(str string) (*DelegatedStake, error) {
	var o DelegatedStake
	err := base.FromJsonString(str, &o)
	return &o, err
}

func NewDelegatedStakeArrayWithJsonString(str string) (*DelegatedStakeArray, error) {
	var o DelegatedStakeArray
	err := base.FromJsonString(str, &o)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (c *Chain) GetValidatorState() (s *ValidatorState, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli, err := c.Client()
	if err != nil {
		return nil, err
	}
	state, err := cli.GetLatestSuiSystemState(context.Background())
	if err != nil {
		return nil, err
	}
	cachedSuiSystemState = state // cache

	apys := c.getValidatorsApy(true)

	totalRewards := decimal.Zero
	var validators = []*Validator{}
	for _, v := range state.ActiveValidators {
		validator := mapRawValidator(&v, apys)
		if validator.isSpecified() {
			validators = append([]*Validator{validator}, validators...)
		} else {
			validators = append(validators, validator)
		}
		totalRewards.Add(v.RewardsPool.Decimal())
	}
	totalStake := state.TotalStake
	res := &ValidatorState{
		Epoch:      state.Epoch.Int64(),
		Validators: &ValidatorArray{validators},

		TotalStaked:           strconv.FormatUint(totalStake.Uint64(), 10),
		TotalRewards:          totalRewards.String(),
		EpochDurationMs:       state.EpochDurationMs.Int64(),
		EpochStartTimestampMs: state.EpochStartTimestampMs.Int64(),
	}

	return res, nil
}

func (c *Chain) GetValidator(address string, useCache bool) (v *Validator, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	var state *types.SuiSystemStateSummary = nil
	if useCache && cachedSuiSystemState != nil {
		state = cachedSuiSystemState
	}
	if state == nil {
		cli, err := c.Client()
		if err != nil {
			return nil, err
		}
		state, err = cli.GetLatestSuiSystemState(context.Background())
		if err != nil {
			return nil, err
		}
		cachedSuiSystemState = state
	}
	apys := c.getValidatorsApy(true)

	for _, val := range state.ActiveValidators {
		if types.IsSameStringAddress(address, val.SuiAddress.String()) {
			validator := mapRawValidator(&val, apys)
			return validator, nil
		}
	}
	return nil, errors.New("not found")
}

var apysCache map[string]float64

// @return not null
func (c *Chain) getValidatorsApy(useCache bool) map[string]float64 {
	if useCache && apysCache != nil {
		return apysCache
	}
	cli, err := c.Client()
	if err != nil {
		return map[string]float64{}
	}
	apys, err := cli.GetValidatorsApy(context.Background())
	if err != nil {
		return map[string]float64{}
	}
	apysCache = apys.ApyMap()
	return apysCache
}

// @return Array of `DelegatedStake` elements
func (c *Chain) GetDelegatedStakes(owner string) (arr *DelegatedStakeArray, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	addr, err := sui_types.NewAddressFromHex(owner)
	if err != nil {
		return nil, err
	}
	cli, err := c.Client()
	if err != nil {
		return nil, err
	}
	list, err := cli.GetStakes(context.Background(), *addr)
	if err != nil {
		return nil, err
	}
	if cachedSuiSystemState == nil {
		cachedSuiSystemState, _ = cli.GetLatestSuiSystemState(context.Background())
	}
	apys := c.getValidatorsApy(true)

	var stakes = []*DelegatedStake{}
	for _, s := range list {
		stakeArray := mapRawStake(&s, apys)
		for _, stake := range stakeArray {
			stakes = append(stakes, stake)
		}
	}

	cachedDelegatedStakesMap.Store(owner, stakes)
	return &DelegatedStakeArray{stakes}, nil
}

// @useCache If true, when there is cached data, the result will be returned directly without requesting data on the chain.
func (c *Chain) TotalStakedSuiAtValidator(validator, owner string, useCache bool) (sui *base.OptionalString, err error) {
	var stakes *DelegatedStakeArray
	if useCache {
		if cachedStakes, ok := cachedDelegatedStakesMap.Load(owner); ok {
			if v, ok := cachedStakes.(*DelegatedStakeArray); ok {
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
	for _, stake := range stakes.AnyArray {
		if principalInt, ok := big.NewInt(0).SetString(stake.Principal, 10); ok {
			total = total.Add(total, principalInt)
		}
	}
	return &base.OptionalString{Value: total.String()}, nil
}

func AverageApyOfDelegatedStakes(stakes *DelegatedStakeArray) float64 {
	if stakes == nil {
		return 0
	}
	totalApy := float64(0)
	added := map[string]bool{}
	for _, stake := range stakes.AnyArray {
		if added[stake.ValidatorAddress] != true && stake.Validator != nil {
			totalApy = totalApy + stake.Validator.APY
			added[stake.ValidatorAddress] = true
			continue
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

	signer, err := sui_types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	validator, err := sui_types.NewAddressFromHex(validatorAddress)
	if err != nil {
		return
	}
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok {
		return nil, errors.New("invalid stake amount")
	}
	cli, err := c.Client()
	if err != nil {
		return
	}
	coinType := SUI_COIN_TYPE
	coins, err := cli.GetCoins(context.Background(), *signer, &coinType, nil, MAX_INPUT_COUNT_MERGE)
	if err != nil {
		return
	}
	pickedCoins, err := types.PickupCoins(coins, *amountInt, maxGasBudgetForStake, MAX_INPUT_COUNT_MERGE, 0)
	if err != nil {
		return
	}
	gasPrice, _ := c.CachedGasPrice()
	maxGasBudget := maxGasBudget(pickedCoins, maxGasBudgetForStake)
	return c.EstimateTransactionFeeAndRebuildTransactionBCS(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		txBytes, err := client.BCS_RequestAddStake(*signer,
			pickedCoins.CoinRefs(),
			types.NewSafeSuiBigInt(amountInt.Uint64()),
			*validator,
			gasBudget,
			gasPrice,
		)
		if err != nil {
			return nil, err
		}
		return &Transaction{
			TxnBytes: txBytes,
		}, nil
	})
}

func (c *Chain) WithdrawDelegation(owner, stakeId string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := sui_types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	stakeSui, err := sui_types.NewObjectIdFromHex(stakeId)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}
	pickedGas, err := c.PickGasCoins(*signer, maxGasBudgetForStake)
	if err != nil {
		return
	}
	stakeSuiObject, err := cli.GetObject(context.Background(), *stakeSui, nil)
	if err != nil {
		return
	}

	maxGasBudget := pickedGas.SuggestMaxGasBudget()
	gasPrice, _ := c.CachedGasPrice()
	return c.EstimateTransactionFeeAndRebuildTransactionBCS(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		txnBytes, err := client.BCS_RequestWithdrawStake(*signer,
			stakeSuiObject.Data.Reference(), pickedGas.CoinRefs(),
			gasBudget, gasPrice)
		if err != nil {
			return nil, err
		}
		return &Transaction{
			TxnBytes: txnBytes,
		}, nil
	})
}

func mapRawValidator(v *types.SuiValidatorSummary, apys map[string]float64) *Validator {
	if v == nil {
		return nil
	}
	// selfStaked := v.StakeAmount
	// delegatedStaked := v.DelegationStakingPool.SuiBalance
	totalStaked := v.StakingPoolSuiBalance
	rewardsPoolBalance := v.RewardsPool

	validator := Validator{
		Address:    v.SuiAddress.String(),
		Name:       v.Name,
		Desc:       v.Description,
		ImageUrl:   v.ImageUrl,
		ProjectUrl: v.ProjectUrl,
		APY:        apys[v.SuiAddress.String()] * 100,

		Commission:      v.CommissionRate.Int64(),
		SelfStaked:      "--",
		DelegatedStaked: "--",
		TotalStaked:     strconv.FormatInt(totalStaked.Int64(), 10),
		TotalRewards:    strconv.FormatInt(rewardsPoolBalance.Int64(), 10),
		GasPrice:        v.GasPrice.Int64(),
	}
	if cachedSuiSystemState.TotalStake.Int64() != 0 {
		validator.PoolShare, _ = v.StakingPoolSuiBalance.Decimal().
			Div(cachedSuiSystemState.TotalStake.Decimal()).
			Mul(decimal.NewFromInt(100)).Float64()
	}
	return &validator
}

func mapRawStake(s *types.DelegatedStake, apys map[string]float64) []*DelegatedStake {
	if s == nil {
		return nil
	}

	var sameValidator *Validator = nil
	if cachedSuiSystemState != nil {
		for _, v := range cachedSuiSystemState.ActiveValidators {
			if v.SuiAddress.ShortString() == s.ValidatorAddress.ShortString() {
				sameValidator = mapRawValidator(&v, apys)
			}
		}
	}

	stakeArray := []*DelegatedStake{}
	for _, json := range s.Stakes {
		stake := json.Data
		requestEpoch := stake.StakeRequestEpoch.Int64()
		stakeRes := &DelegatedStake{
			StakeId:          stake.StakedSuiId.String(),
			ValidatorAddress: s.ValidatorAddress.String(),
			Principal:        strconv.FormatUint(stake.Principal.Uint64(), 10),
			RequestEpoch:     requestEpoch,
			Validator:        sameValidator,
		}
		status := stake.StakeStatus.Data
		if status.Active != nil {
			stakeRes.EarnedAmount = strconv.FormatUint(status.Active.EstimatedReward.Uint64(), 10)
			stakeRes.Status = DelegationStatusActived
		} else {
			stakeRes.EarnedAmount = "0"
			stakeRes.Status = DelegationStatusPending
		}
		stakeArray = append(stakeArray, stakeRes)
	}

	return stakeArray
}
