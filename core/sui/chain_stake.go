package sui

import (
	"context"
	"errors"
	"math/big"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/coming-chat/go-sui/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/shopspring/decimal"
)

var cachedSuiSystemState *types.SuiSystemStateSummary
var cachedDelegatedStakesMap sync.Map

const maxGasBudgetForStake = 12000000

type ValidatorState struct {
	// The current epoch in Sui. An epoch takes approximately 24 hours and runs in checkpoints.
	Epoch int64 `json:"epoch"`
	// Array of `Validator` elements
	Validators *base.AnyArray `json:"validators"`

	// The amount of all tokens staked in the Sui Network.
	TotalStaked string `json:"totalStaked"`
	// The amount of rewards won by all Sui validators in the last epoch.
	TotalRewards string `json:"lastEpochReward"`

	EpochStartTimestampMs int64 `json:"epochStartTimestampMs"`
	EpochDurationMs       int64 `json:"epochDurationMs"`
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

	PoolShare float64 `json:"poolShare"`
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

// @return if time > 0 indicates how long it will take to get the reward;
// if time < 0 indicates how much time has passed since the reward was earned;
func (s *DelegatedStake) EarningAmountTimeAfterNowMs(stateInfo *ValidatorState) int64 {
	return s.EarningAmountTimeAfterTimestampMs(time.Now().UnixMilli(), *stateInfo)
}

func (s *DelegatedStake) EarningAmountTimeAfterTimestampMs(timestamp int64, stateInfo ValidatorState) int64 {
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
	var validators = &base.AnyArray{}
	for _, v := range state.ActiveValidators {
		validator := mapRawValidator(&v, apys)
		if validator.isSpecified() {
			validators.Values = append([]any{validator}, validators.Values...)
		} else {
			validators.Values = append(validators.Values, validator)
		}
		totalRewards.Add(v.RewardsPool.Decimal())
	}
	totalStake := state.TotalStake
	res := &ValidatorState{
		Epoch:      state.Epoch.Int64(),
		Validators: validators,

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
func (c *Chain) GetDelegatedStakes(owner string) (arr *base.AnyArray, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	addr, err := types.NewAddressFromHex(owner)
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

	var stakes = &base.AnyArray{}
	for _, s := range list {
		stakeArray := mapRawStake(&s, apys)
		for _, stake := range stakeArray {
			stakes.Values = append(stakes.Values, stake)
		}
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
	if stakes == nil {
		return 0
	}
	totalApy := float64(0)
	added := map[string]bool{}
	for _, val := range stakes.Values {
		if stake, ok := val.(*DelegatedStake); ok {
			if added[stake.ValidatorAddress] != true && stake.Validator != nil {
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
	cli, err := c.Client()
	if err != nil {
		return
	}
	coinType := SUI_COIN_TYPE
	coins, err := cli.GetCoins(context.Background(), *signer, &coinType, nil, MAX_INPUT_COUNT_STAKE)
	if err != nil {
		return
	}
	pickedCoins, err := types.PickupCoins(coins, *amountInt, MAX_INPUT_COUNT_STAKE, true)
	if err != nil {
		return
	}
	maxGasBudget := uint64(maxGasBudgetForStake)
	if pickedCoins.RemainingMaxCoinValue > 0 {
		maxGasBudget = base.Min(maxGasBudget, pickedCoins.RemainingMaxCoinValue)
	}
	return c.EstimateTransactionFeeAndRebuildTransaction(maxGasBudget, func(gasBudget uint64) (*Transaction, error) {
		gasInt := big.NewInt(0).SetUint64(gasBudget)
		txBytes, err := cli.RequestAddStake(context.Background(), *signer,
			pickedCoins.CoinIds(),
			decimal.NewFromBigInt(amountInt, 0),
			*validator,
			nil, decimal.NewFromBigInt(gasInt, 0))
		if err != nil {
			return nil, err
		}
		return &Transaction{Txn: *txBytes}, nil
	})
}

func (c *Chain) WithdrawDelegation(owner, stakeId string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	signer, err := types.NewAddressFromHex(owner)
	if err != nil {
		return
	}
	stakeSui, err := types.NewHexData(stakeId)
	if err != nil {
		return
	}
	cli, err := c.Client()
	if err != nil {
		return
	}
	return c.EstimateTransactionFeeAndRebuildTransaction(MinGasBudget, func(gasBudget uint64) (*Transaction, error) {
		gasInt := big.NewInt(0).SetUint64(gasBudget)
		txnBytes, err := cli.RequestWithdrawStake(context.Background(), *signer, *stakeSui, nil, decimal.NewFromBigInt(gasInt, 0))
		if err != nil {
			return nil, err
		}
		return &Transaction{Txn: *txnBytes}, nil
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
