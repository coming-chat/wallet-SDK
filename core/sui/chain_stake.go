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

	totalRewards := decimal.Zero
	var validators = &base.AnyArray{}
	for _, v := range state.ActiveValidators {
		validator := mapRawValidator(&v, state.Epoch)
		if validator.isSpecified() {
			validators.Values = append([]any{validator}, validators.Values...)
		} else {
			validators.Values = append(validators.Values, validator)
		}
		totalRewards.Add(v.RewardsPool)
	}
	totalStake := state.TotalStake
	res := &ValidatorState{
		Epoch:      int64(state.Epoch),
		Validators: validators,

		TotalStaked:           strconv.FormatUint(totalStake, 10),
		TotalRewards:          totalRewards.String(),
		EpochDurationMs:       int64(state.EpochDurationMs),
		EpochStartTimestampMs: int64(state.EpochStartTimestampMs),
	}

	return res, nil
}

func (c *Chain) GetValidator(address string, useCache bool) (v *Validator, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	var state *types.SuiSystemStateSummary = nil
	if useCache && cachedSuiSystemState != nil {
		state = cachedSuiSystemState
	}
	if cachedSuiSystemState == nil {
		cli, err := c.Client()
		if err != nil {
			return nil, err
		}
		state, err = cli.GetLatestSuiSystemState(context.Background())
		if err != nil {
			return nil, err
		}
	}

	for _, val := range state.ActiveValidators {
		if types.IsSameStringAddress(address, val.SuiAddress.String()) {
			validator := mapRawValidator(&val, state.Epoch)
			return validator, nil
		}
	}
	return nil, errors.New("not found")
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

	var stakes = &base.AnyArray{}
	for _, s := range list {
		stakeArray := mapRawStake(&s)
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
	allCoins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	if err != nil {
		return
	}
	needCoins, gasCoin, err := allCoins.PickSUICoinsWithGas(amountInt, maxGasBudgetForStake, types.PickBigger)
	if err != nil {
		return
	}
	coinIds := []types.ObjectId{}
	for _, coin := range needCoins {
		coinIds = append(coinIds, coin.CoinObjectId)
	}
	gasId := gasCoin.CoinObjectId
	txBytes, err := cli.RequestAddStake(context.Background(), *signer, coinIds, amountInt.Uint64(), *validator, &gasId, maxGasBudgetForStake)
	if err != nil {
		return
	}
	return &Transaction{
		Txn:          *txBytes,
		MaxGasBudget: maxGasBudgetForStake,
	}, nil
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
	allCoins, err := cli.GetSuiCoinsOwnedByAddress(context.Background(), *signer)
	if err != nil {
		return
	}
	gasCoin, err := allCoins.PickCoinNoLess(maxGasBudgetForStake)
	if err != nil {
		return
	}
	gasId := gasCoin.CoinObjectId
	txnBytes, err := cli.RequestWithdrawStake(context.Background(), *signer, *stakeSui, &gasId, maxGasBudgetForStake)
	if err != nil {
		return
	}

	return &Transaction{
		Txn:          *txnBytes,
		MaxGasBudget: maxGasBudgetForStake,
	}, nil
}

func mapRawValidator(v *types.SuiValidatorSummary, epoch uint64) *Validator {
	if v == nil {
		return nil
	}
	// selfStaked := v.StakeAmount
	// delegatedStaked := v.DelegationStakingPool.SuiBalance
	// totalStaked := selfStaked + delegatedStaked
	rewardsPoolBalance := v.RewardsPool

	validator := Validator{
		Address:    v.SuiAddress.String(),
		Name:       v.Name,
		Desc:       v.Description,
		ImageUrl:   v.ImageUrl,
		ProjectUrl: v.ProjectUrl,
		APY:        v.CalculateAPY(epoch),

		Commission:      v.CommissionRate,
		SelfStaked:      "--",
		DelegatedStaked: "--",
		TotalStaked:     "--",
		TotalRewards:    rewardsPoolBalance.String(),
		GasPrice:        int64(v.GasPrice),
	}
	return &validator
}

func mapRawStake(s *types.DelegatedStake) []*DelegatedStake {
	if s == nil {
		return nil
	}

	var sameValidator *Validator = nil
	if cachedSuiSystemState != nil {
		for _, v := range cachedSuiSystemState.ActiveValidators {
			if v.SuiAddress.ShortString() == s.ValidatorAddress.ShortString() {
				sameValidator = mapRawValidator(&v, cachedSuiSystemState.Epoch)
			}
		}
	}

	stakeArray := []*DelegatedStake{}
	for _, stake := range s.Stakes {
		requestEpoch, _ := strconv.ParseInt(stake.StakeRequestEpoch, 10, 64)
		stakeRes := &DelegatedStake{
			StakeId:          stake.StakedSuiId.String(),
			ValidatorAddress: s.ValidatorAddress.String(),
			Principal:        strconv.FormatUint(stake.Principal, 10),
			RequestEpoch:     requestEpoch,
			Validator:        sameValidator,
		}
		if stake.EstimatedReward != nil {
			stakeRes.EarnedAmount = strconv.FormatUint(*stake.EstimatedReward, 10)
		} else {
			stakeRes.EarnedAmount = "0"
		}
		if stake.Status == types.StakeStatusActive {
			stakeRes.Status = DelegationStatusActived
		} else {
			stakeRes.Status = DelegationStatusPending
		}
		stakeArray = append(stakeArray, stakeRes)
	}

	return stakeArray
}
