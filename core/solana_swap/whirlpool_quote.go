package solanaswap

import (
	"context"
	"errors"
	"math"
	"math/big"
	"strconv"

	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/common"
)

type Percentage struct {
	Numerator   uint64
	Denominator uint64
}

func NewPercentage(numerator, denominator uint64) *Percentage {
	return &Percentage{Numerator: numerator, Denominator: denominator}
}

var defaultSlippagePercentage = Percentage{1, 1000}
var ZERO_SLIPPAGE = Percentage{0, 1}

const WhirlpoolProgramId = "whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"
const MAX_TICK_ARRAY_CROSSINGS = 2
const TICK_ARRAY_SIZE = 88
const MAX_TICK_INDEX = 443636
const MIN_TICK_INDEX = -443636
const MAX_SQRT_PRICE = "79226673515401279992447579055"
const MIN_SQRT_PRICE = "4295048016"

type SwapQuoteParam struct {
	poolAddress       string
	tokenMint         string
	tokenAmount       uint64
	isInput           bool
	slippageTolerance Percentage
	refresh           bool
}

type SwapQuote struct {
	poolAddress          string
	otherAmountThreshold uint64
	sqrtPriceLimitX64    *big.Int
	amountIn             uint64
	amountOut            uint64
	aToB                 bool
	fixedInput           bool
}

const (
	SwapDirectionAtoB = true
	SwapDirectionBtoA = false
)

func GetSwapQuote(cli *client.Client, param SwapQuoteParam) (quote *SwapQuote, err error) {
	whirlpool, err := getPoolData(cli, param.poolAddress)
	if err != nil {
		return
	}
	swapDirection := whirlpool.TokenMintA.ToBase58() == param.tokenMint

	output, err := simulateSwap(cli,
		SwapSimulationBaseInput{
			poolAddress:   param.poolAddress,
			whirlpoolData: whirlpool,
			swapDirection: swapDirection,
		},
		SwapSimulationInput{
			amount:              big.NewInt(0).SetUint64(param.tokenAmount),
			currentSqrtPriceX64: whirlpool.SqrtPrice,
			currentTickIndex:    whirlpool.TickCurrentIndex,
			currentLiquidity:    whirlpool.Liquidity,
		},
	)
	if err != nil {
		return
	}

	otherAmountThreshold := adjustAmountForSlippage(
		output.amountIn,
		output.amountOut,
		param.slippageTolerance,
	)

	return &SwapQuote{
		poolAddress:          param.poolAddress,
		otherAmountThreshold: otherAmountThreshold.Uint64(),
		sqrtPriceLimitX64:    output.sqrtPriceLimitX64,
		amountIn:             quote.amountIn,
		amountOut:            quote.amountOut,
		aToB:                 swapDirection == SwapDirectionAtoB,
		fixedInput:           true,
	}, nil
}

type SwapSimulationBaseInput struct {
	poolAddress   string
	whirlpoolData *WhirlpoolData
	swapDirection bool // AtoB / BtoA
}

type SwapSimulationInput struct {
	amount              *big.Int
	currentSqrtPriceX64 *big.Int
	currentTickIndex    int32
	currentLiquidity    *big.Int
}

type SwapSimulationOutput struct {
	amountIn          *big.Int
	amountOut         *big.Int
	sqrtPriceLimitX64 *big.Int
}

type SwapStepSimulationInput struct {
	sqrtPriceX64      *big.Int
	tickIndex         int32
	liquidity         *big.Int
	amountRemaining   uint64
	tickArraysCrossed int
}

type SwapStepSimulationOutput struct {
	nextSqrtPriceX64   *big.Int
	nextTickIndex      int32
	input              *big.Int
	output             *big.Int
	tickArraysCrossed  int
	hasReachedNextTick bool
}

func simulateSwap(cli *client.Client,
	baseInput SwapSimulationBaseInput,
	input SwapSimulationInput) (out *SwapSimulationOutput, err error) {
	specifiedAmountLeft := input.amount
	currentLiquidity := input.currentLiquidity
	currentTickIndex := input.currentTickIndex
	currentSqrtPriceX64 := input.currentSqrtPriceX64

	swapDirection := baseInput.swapDirection

	if specifiedAmountLeft.Cmp(big.NewInt(0)) <= 0 {
		return nil, errors.New("amount must be nonzero")
	}
	otherAmountCalculated := big.NewInt(0)
	tickArraysCrossed := 0
	var sqrtPriceLimitX64 *big.Int

	for specifiedAmountLeft.Cmp(big.NewInt(0)) > 0 {
		if tickArraysCrossed > MAX_TICK_ARRAY_CROSSINGS {
			return nil, errors.New("Crossed the maximum number of tick arrays")
		}

		swapStepSimulationOutput, err := simulateSwapStep(cli, baseInput, SwapStepSimulationInput{
			sqrtPriceX64:      input.currentSqrtPriceX64,
			amountRemaining:   specifiedAmountLeft.Uint64(),
			tickIndex:         currentTickIndex,
			liquidity:         currentLiquidity,
			tickArraysCrossed: tickArraysCrossed,
		})
		if err != nil {
			return nil, err
		}

		nextTickIndex := swapStepSimulationOutput.nextTickIndex

		specifiedAmountUsed, otherAmount := swapStepSimulationOutput.input, swapStepSimulationOutput.output

		specifiedAmountLeft.Sub(specifiedAmountLeft, specifiedAmountUsed)
		otherAmountCalculated.Add(otherAmountCalculated, otherAmount)

		if swapStepSimulationOutput.hasReachedNextTick {
			nextTick, err := fetchTick(cli, baseInput, nextTickIndex)
			if err != nil {
				return nil, err
			}
			currentLiquidity = calculateNewLiquidity(currentLiquidity, nextTick.liquidityNet, swapDirection)
			if swapDirection == SwapDirectionAtoB {
				currentTickIndex = nextTickIndex - 1
			} else {
				currentTickIndex = nextTickIndex
			}
		}

		currentSqrtPriceX64 = swapStepSimulationOutput.nextSqrtPriceX64
		println(currentSqrtPriceX64.String())
		tickArraysCrossed := swapStepSimulationOutput.tickArraysCrossed

		if tickArraysCrossed > MAX_TICK_ARRAY_CROSSINGS {
			sqrtPriceLimitX64, _ = tickIndexToSqrtPriceX64(nextTickIndex)
		}
	}

	inputAmount, outputAmount := big.NewInt(0).Sub(input.amount, specifiedAmountLeft), otherAmountCalculated

	if sqrtPriceLimitX64 == nil {
		if swapDirection == SwapDirectionAtoB {
			sqrtPriceLimitX64 = MustInt(MIN_SQRT_PRICE)
		} else {
			sqrtPriceLimitX64 = MustInt(MAX_SQRT_PRICE)
		}
	}

	return &SwapSimulationOutput{
		amountIn:          inputAmount,
		amountOut:         outputAmount,
		sqrtPriceLimitX64: sqrtPriceLimitX64,
	}, nil
}

func simulateSwapStep(cli *client.Client,
	baseInput SwapSimulationBaseInput,
	input SwapStepSimulationInput) (out *SwapStepSimulationOutput, err error) {

	swapDirection := baseInput.swapDirection
	liquidity := input.liquidity
	sqrtPriceX64 := input.sqrtPriceX64

	feeRatePercentage := getFeeRate(baseInput.whirlpoolData.FeeRate)

	nextTickIndex, tickArraysCrossedUpdate, err := getNextInitializedTickIndex(cli, baseInput, input.tickIndex, input.tickArraysCrossed)
	if err != nil {
		return
	}
	targetSqrtPriceX64, err := tickIndexToSqrtPriceX64(nextTickIndex)
	if err != nil {
		return
	}
	fixedDelta := getAmountFixedDelta(
		sqrtPriceX64,
		targetSqrtPriceX64,
		liquidity,
		swapDirection)
	amountCalculated := calculateAmountAfterFees(input.amountRemaining, *feeRatePercentage)

	var nextSqrtPriceX64 *big.Int
	if amountCalculated.Cmp(fixedDelta) >= 0 {
		nextSqrtPriceX64 = targetSqrtPriceX64
	} else {
		nextSqrtPriceX64 = getNextSqrtPrice(sqrtPriceX64, liquidity, amountCalculated, swapDirection)
	}

	hasReachedNextTick := nextSqrtPriceX64.Cmp(targetSqrtPriceX64) == 0

	unfixedDelta := getAmountUnfixedDelta(
		sqrtPriceX64,
		nextSqrtPriceX64,
		liquidity,
		swapDirection,
	)
	if !hasReachedNextTick {
		fixedDelta = getAmountFixedDelta(
			sqrtPriceX64,
			nextSqrtPriceX64,
			liquidity,
			swapDirection,
		)
	}

	inputDelta, outputDelta := fixedDelta, unfixedDelta

	if !hasReachedNextTick {
		inputDelta = big.NewInt(0).SetUint64(input.amountRemaining)
	} else {
		temp := calculateFeesFromAmount(inputDelta.Uint64(), *feeRatePercentage)
		inputDelta = big.NewInt(0).Add(inputDelta, temp)
	}

	return &SwapStepSimulationOutput{
		nextTickIndex:      nextTickIndex,
		nextSqrtPriceX64:   nextSqrtPriceX64,
		input:              inputDelta,
		output:             outputDelta,
		tickArraysCrossed:  tickArraysCrossedUpdate,
		hasReachedNextTick: hasReachedNextTick,
	}, nil
}

func tickIndexToSqrtPriceX64(tickIndex int32) (*big.Int, error) {
	if tickIndex > MAX_TICK_INDEX || tickIndex < MIN_TICK_INDEX {
		return nil, errors.New("Provided tick index does not fit within supported tick index range.")
	}
	if tickIndex > 0 {
		return tick_index_to_sqrt_price_positive(tickIndex), nil
	} else {
		return tick_index_to_sqrt_price_negative(tickIndex), nil
	}
}

func getNextInitializedTickIndex(cli *client.Client,
	baseInput SwapSimulationBaseInput,
	currentTickIndex int32,
	tickArraysCrossed int) (tickIndex int32, crossed int, err error) {

	tickSpacing := baseInput.whirlpoolData.TickSpacing

	var nextInitializedTickIndex *int32
	for nextInitializedTickIndex == nil {
		currentTickArray, err := fetchTickArray(cli, baseInput, currentTickIndex)
		if err != nil {
			return 0, 0, err
		}
		temp := int32(0)
		if baseInput.swapDirection == SwapDirectionAtoB {
			temp, err = getPrevInitializedTickIndex(*currentTickArray, currentTickIndex, baseInput.whirlpoolData.TickSpacing)
		} else {
			temp, err = UgetNextInitializedTickIndex(*currentTickArray, currentTickIndex, baseInput.whirlpoolData.TickSpacing)
		}
		if err == nil {
			nextInitializedTickIndex = &temp
		} else if tickArraysCrossed == MAX_TICK_ARRAY_CROSSINGS {
			if baseInput.swapDirection == SwapDirectionAtoB {
				temp = currentTickArray.startTickIndex
			} else {
				temp = currentTickArray.startTickIndex + TICK_ARRAY_SIZE*int32(tickSpacing)
			}
			nextInitializedTickIndex = &temp
			tickArraysCrossed++
		} else {
			if baseInput.swapDirection == SwapDirectionAtoB {
				currentTickIndex = currentTickArray.startTickIndex - 1
			} else {
				currentTickIndex = currentTickArray.startTickIndex + TICK_ARRAY_SIZE*int32(tickSpacing)
			}
			tickArraysCrossed++
		}
	}

	return *nextInitializedTickIndex, tickArraysCrossed, nil
}

func fetchTickArray(cli *client.Client,
	baseInput SwapSimulationBaseInput,
	tickIndex int32) (data *TickArrayData, err error) {
	address, err := getPdaWithTickIndex(tickIndex, baseInput.whirlpoolData.TickSpacing, baseInput.poolAddress, WhirlpoolProgramId, 0)
	if err != nil {
		return nil, err
	}
	return getTickArray(cli, address.ToBase58())
}

func fetchTick(cli *client.Client, baseInput SwapSimulationBaseInput, tickIndex int32) (data *TickData, err error) {
	tickArray, err := fetchTickArray(cli, baseInput, tickIndex)
	if err != nil {
		return
	}
	tickSpacing := baseInput.whirlpoolData.TickSpacing
	realIndex := tickIndexToTickArrayIndex(tickArray.startTickIndex, tickIndex, tickSpacing)
	if realIndex < 0 || int(realIndex) >= len(tickArray.ticks) {
		return nil, errors.New("tick realIndex out of range")
	}
	return &tickArray.ticks[realIndex], nil
}

func getTickArray(cli *client.Client, address string) (data *TickArrayData, err error) {
	info, err := cli.GetAccountInfo(context.Background(), address)
	if err != nil {
		return
	}
	var value TickArrayData
	err = value.Deserializer(info.Data)
	if err != nil {
		return nil, err
	}
	return &value, nil
}

func getFeeRate(feeRate uint16) *Percentage {
	/**
	 * Smart Contract comment: https://github.com/orca-so/whirlpool/blob/main/programs/whirlpool/src/state/whirlpool.rs#L9-L11
	 * // Stored as hundredths of a basis point
	 * // u16::MAX corresponds to ~6.5%
	 * pub fee_rate: u16,
	 */
	return NewPercentage(uint64(feeRate), 1e6)
}

type TickArrayData struct {
	whirlpool      common.PublicKey
	startTickIndex int32
	ticks          []TickData
}

type TickData struct {
	initialized          bool
	liquidityNet         *big.Int
	liquidityGross       *big.Int
	feeGrowthOutsideA    *big.Int
	feeGrowthOutsideB    *big.Int
	rewardGrowthsOutside []*big.Int
}

func (d *TickArrayData) Deserializer(data []byte) error {
	// if len(data) < 8+261+128 {
	// 	return errors.New("data length not enough")
	// }
	ds := NewSolanaDeserializer(data)
	_ = ds.TakeBytes(8) // ignore

	d.startTickIndex = ds.TakeI32()
	ticks := make([]TickData, 0, 88)
	for i := 0; i < 88; i++ {
		tick := TickData{}
		tick.initialized = ds.TakeBool()
		tick.liquidityNet = ds.TakeI128()
		tick.liquidityGross = ds.TakeU128()
		tick.feeGrowthOutsideA = ds.TakeU128()
		tick.feeGrowthOutsideB = ds.TakeU128()
		tick.rewardGrowthsOutside = []*big.Int{
			ds.TakeU128(),
			ds.TakeU128(),
			ds.TakeU128(),
		}
	}
	d.ticks = ticks
	d.whirlpool = ds.TakePublicKey()
	return nil
}

func getPdaWithTickIndex(
	tickIndex int32,
	tickSpacing uint16,
	whirlpool string,
	programId string,
	tickArrayOffset int) (common.PublicKey, error) {
	startIndex, err := getStartTickIndex(tickIndex, tickSpacing, tickArrayOffset)
	if err != nil {
		return common.PublicKey{}, err
	}
	return getTickArrayPDA(programId, whirlpool, startIndex)
}

/**
 * Get the startIndex of the tick array containing tickIndex.
 *
 * @param tickIndex
 * @param tickSpacing
 * @param offset can be used to get neighboring tick array startIndex.
 * @returns
 */
func getStartTickIndex(tickIndex int32, tickSpacing uint16, offset int) (int, error) {
	realIndex := math.Floor(float64(tickIndex) / float64(tickSpacing) / float64(TICK_ARRAY_SIZE))
	startTickIndex := (int(realIndex) + offset) * int(tickSpacing) * TICK_ARRAY_SIZE

	ticksInArray := TICK_ARRAY_SIZE * int(tickSpacing)
	minTickIndex := MIN_TICK_INDEX - ((MIN_TICK_INDEX % ticksInArray) + ticksInArray)
	if startTickIndex < minTickIndex {
		return 0, errors.New("startTickIndex is too small")
	}
	if startTickIndex > MAX_TICK_INDEX {
		return 0, errors.New("startTickIndex is too large")
	}
	return startTickIndex, nil
}

func getTickArrayPDA(programId string, whirlpoolAddress string, startTick int) (common.PublicKey, error) {
	PDA_TICK_ARRAY_SEED := "tick_array"
	whirlpoolPub := common.PublicKeyFromString(whirlpoolAddress)
	return findProgramAddressSeed(programId, [][]byte{
		[]byte(PDA_TICK_ARRAY_SEED),
		whirlpoolPub[:],
		[]byte(strconv.FormatInt(int64(startTick), 10)),
	})
}

const (
	TickSearchDirectionLeft  = 0
	TickSearchDirectionRight = 1
)

func getPrevInitializedTickIndex(account TickArrayData,
	currentTickIndex int32,
	tickSpacing uint16) (int32, error) {
	return findInitializedTick(
		account,
		currentTickIndex,
		tickSpacing,
		TickSearchDirectionLeft)
}

func UgetNextInitializedTickIndex(account TickArrayData,
	currentTickIndex int32,
	tickSpacing uint16) (int32, error) {
	return findInitializedTick(
		account,
		currentTickIndex,
		tickSpacing,
		TickSearchDirectionRight)
}

func findInitializedTick(account TickArrayData,
	currentTickIndex int32,
	tickSpacing uint16,
	searchDirection int) (int32, error) {
	currentTickArrayIndex := tickIndexToTickArrayIndex(
		account.startTickIndex,
		currentTickIndex,
		tickSpacing)

	increment := int32(1)
	stepInitializedTickArrayIndex := int32(0)
	if searchDirection == TickSearchDirectionRight {
		increment = 1
		stepInitializedTickArrayIndex = currentTickArrayIndex + increment
	} else {
		increment = -1
		stepInitializedTickArrayIndex = currentTickIndex
	}

	for stepInitializedTickArrayIndex >= 0 &&
		stepInitializedTickArrayIndex < int32(len(account.ticks)) {
		if account.ticks[stepInitializedTickArrayIndex].initialized {
			return tickArrayIndexToTickIndex(account.startTickIndex,
				stepInitializedTickArrayIndex,
				tickSpacing), nil
		}
		stepInitializedTickArrayIndex++
	}

	return 0, errors.New("not found")
}

func tickIndexToTickArrayIndex(startTickIndex int32, tickIndex int32, tickSpacing uint16) int32 {
	return int32(math.Floor((float64(tickIndex - startTickIndex)) / float64(tickSpacing)))
}

func tickArrayIndexToTickIndex(startTickIndex int32, tickArrayIndex int32, tickSpacing uint16) int32 {
	return startTickIndex + tickArrayIndex*int32(tickSpacing)
}

func getAmountFixedDelta(currentSqrtPriceX64, targetSqrtPriceX64, liquidity *big.Int, swapDirection bool) *big.Int {
	if swapDirection == SwapDirectionAtoB {
		return getTokenAFromLiquidity(liquidity, currentSqrtPriceX64, targetSqrtPriceX64, true)
	} else {
		return getTokenBFromLiquidity(liquidity, currentSqrtPriceX64, targetSqrtPriceX64, true)
	}
}

func adjustAmountForSlippage(amountIn, amountOut *big.Int, precent Percentage) *big.Int {
	den := big.NewInt(0).SetUint64(precent.Denominator)
	num := big.NewInt(0).SetUint64(precent.Numerator)
	x := big.NewInt(0).Mul(amountOut, den)
	y := big.NewInt(0).Add(den, num)
	return big.NewInt(0).Div(x, y)
}

func getTokenAFromLiquidity(liquidity, sqrtPrice0X64, sqrtPrice1X64 *big.Int, roundUp bool) *big.Int {
	sqrtPriceLowerX64, sqrtPriceUpperX64 := orderASC(sqrtPrice0X64, sqrtPrice1X64)

	numerator_ := big.NewInt(0).Mul(liquidity, big.NewInt(0).Sub(sqrtPriceUpperX64, sqrtPriceLowerX64))
	numerator := numerator_.Lsh(numerator_, 64)
	denominator := big.NewInt(0).Mul(sqrtPriceUpperX64, sqrtPriceLowerX64)
	if roundUp {
		return divRoundUp(numerator, denominator)
	} else {
		return numerator.Div(numerator, denominator)
	}
}

func getTokenBFromLiquidity(liquidity, sqrtPrice0X64, sqrtPrice1X64 *big.Int, roundUp bool) *big.Int {
	sqrtPriceLowerX64, sqrtPriceUpperX64 := orderASC(sqrtPrice0X64, sqrtPrice1X64)

	result := big.NewInt(0).Mul(liquidity, big.NewInt(0).Sub(sqrtPriceUpperX64, sqrtPriceLowerX64))
	if roundUp {
		return shiftRightRoundUp(result)
	} else {
		return result.Rsh(result, 64)
	}
}

func calculateAmountAfterFees(amount uint64, feeRate Percentage) *big.Int {
	amt := big.NewInt(0).SetUint64(amount)
	den := big.NewInt(0).SetUint64(feeRate.Denominator)
	num := big.NewInt(0).SetUint64(feeRate.Numerator)

	res := amt.Mul(amt, big.NewInt(0).Sub(den, num))
	res = res.Div(res, den)
	return res
}

func calculateFeesFromAmount(amount uint64, feeRate Percentage) *big.Int {
	amt := big.NewInt(0).SetUint64(amount)
	den := big.NewInt(0).SetUint64(feeRate.Denominator)
	num := big.NewInt(0).SetUint64(feeRate.Numerator)

	return divRoundUp(big.NewInt(0).Mul(amt, num), big.NewInt(0).Sub(den, num))
}

func calculateNewLiquidity(liquidity, nextLiquidityNet *big.Int, swapDirection bool) *big.Int {
	if swapDirection == SwapDirectionAtoB {
		nextLiquidityNet = big.NewInt(0).Neg(nextLiquidityNet)
	}
	return big.NewInt(0).Add(liquidity, nextLiquidityNet)
}

func getAmountUnfixedDelta(currentSqrtPriceX64, targetSqrtPriceX64, liquidity *big.Int, swapDirection bool) *big.Int {
	if swapDirection == SwapDirectionAtoB {
		return getTokenBFromLiquidity(liquidity, currentSqrtPriceX64, targetSqrtPriceX64, false)
	} else {
		return getTokenAFromLiquidity(liquidity, currentSqrtPriceX64, targetSqrtPriceX64, false)
	}
}

func getNextSqrtPrice(sqrtPriceX64, liquidity, amount *big.Int, swapDirection bool) *big.Int {
	if swapDirection == SwapDirectionAtoB {
		return getLowerSqrtPriceFromTokenA(amount, liquidity, sqrtPriceX64)
	} else {
		return getUpperSqrtPriceFromTokenB(amount, liquidity, sqrtPriceX64)
	}
}

func getLowerSqrtPriceFromTokenA(amount, liquidity, sqrtPriceX64 *big.Int) *big.Int {
	numerator := big.NewInt(0).Lsh(big.NewInt(0).Mul(liquidity, sqrtPriceX64), 64)
	denominator := big.NewInt(0).Add(big.NewInt(0).Lsh(liquidity, 64), big.NewInt(0).Mul(amount, sqrtPriceX64))
	return divRoundUp(numerator, denominator)
}

func getUpperSqrtPriceFromTokenA(amount, liquidity, sqrtPriceX64 *big.Int) *big.Int {
	numerator := big.NewInt(0).Lsh(big.NewInt(0).Mul(liquidity, sqrtPriceX64), 64)
	denominator := big.NewInt(0).Sub(big.NewInt(0).Lsh(liquidity, 64), big.NewInt(0).Mul(amount, sqrtPriceX64))
	return divRoundUp(numerator, denominator)
}

func getLowerSqrtPriceFromTokenB(amount, liquidity, sqrtPriceX64 *big.Int) *big.Int {
	s := divRoundUp(big.NewInt(0).Lsh(amount, 64), liquidity)
	return big.NewInt(0).Sub(sqrtPriceX64, s)
}

func getUpperSqrtPriceFromTokenB(amount, liquidity, sqrtPriceX64 *big.Int) *big.Int {
	s := big.NewInt(0).Div(big.NewInt(0).Lsh(amount, 64), liquidity)
	return big.NewInt(0).Add(sqrtPriceX64, s)
}
