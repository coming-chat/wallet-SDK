package solanaswap

import (
	"math"
	"math/big"
)

func tick_index_to_sqrt_price_positive(tick_index int32) *big.Int {
	tick_index_shifted := tick_index

	ratio := big.NewInt(0)
	if tick_index_shifted&1 == 0 {
		ratio = MustInt("79228162514264337593543950336") // 0
	} else {
		ratio = MustInt("79232123823359799118286999567") // 1
	}
	precalculated_factor := []string{
		"79236085330515764027303304731", // 2
		"79244008939048815603706035061", // 4
		"79259858533276714757314932305", // 8
		"79291567232598584799939703904", // ...
		"79355022692464371645785046466",
		"79482085999252804386437311141",
		"79736823300114093921829183326",
		"80248749790819932309965073892",
		"81282483887344747381513967011",
		"83390072131320151908154831281",
		"87770609709833776024991924138",
		"97234110755111693312479820773",
		"119332217159966728226237229890",
		"179736315981702064433883588727",
		"407748233172238350107850275304",
		"2098478828474011932436660412517",
		"55581415166113811149459800483533",
		"38992368544603139932233054999993551", // 262144
	}
	for _, factor := range precalculated_factor {
		tick_index_shifted = tick_index_shifted >> 1
		if tick_index_shifted&1 != 0 {
			ratio = mul_shift(ratio, MustInt(factor), 96)
		}
	}
	return mul_shift(ratio, big.NewInt(1), 32)
}

func tick_index_to_sqrt_price_negative(tick_index int32) *big.Int {
	tick_index_shifted := int32(math.Abs(float64(tick_index)))

	ratio := big.NewInt(0)
	if tick_index_shifted&1 == 0 {
		ratio = MustInt("18446744073709551616") // 0
	} else {
		ratio = MustInt("18445821805675392311") // 1
	}
	precalculated_factor := []string{
		"18444899583751176498", // 2
		"18443055278223354162", // 4
		"18439367220385604838", // 8
		"18431993317065449817", // ...
		"18417254355718160513",
		"18387811781193591352",
		"18329067761203520168",
		"18212142134806087854",
		"17980523815641551639",
		"17526086738831147013",
		"16651378430235024244",
		"15030750278693429944",
		"12247334978882834399",
		"8131365268884726200",
		"3584323654723342297",
		"696457651847595233",
		"26294789957452057",
		"37481735321082", // 262144
	}
	for _, factor := range precalculated_factor {
		tick_index_shifted = tick_index_shifted >> 1
		if tick_index_shifted&1 != 0 {
			ratio = mul_shift(ratio, MustInt(factor), 64)
		}
	}
	return ratio
}

func mul_shift(a, b *big.Int, shift uint) *big.Int {
	mulInt := big.NewInt(0).Mul(a, b)
	return mulInt.Rsh(mulInt, shift)
}

func MustInt(num string) *big.Int {
	b, _ := big.NewInt(0).SetString(num, 10)
	return b
}

func divRoundUp(x, y *big.Int) *big.Int {
	var quotient, remainder big.Int
	quotient.DivMod(x, y, &remainder)
	if remainder.Cmp(big.NewInt(0)) > 0 {
		return quotient.Add(&quotient, big.NewInt(1))
	} else {
		return &quotient
	}
}

func orderASC(x, y *big.Int) (samll, big *big.Int) {
	if x.Cmp(y) <= 0 {
		return x, y
	} else {
		return y, x
	}
}

func shiftRightRoundUp(n *big.Int) *big.Int {
	result := big.NewInt(0).Rsh(n, 64)
	U64_MAX := big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(64), nil), big.NewInt(1))

	if big.NewInt(0).Mod(n, U64_MAX).Cmp(big.NewInt(0)) > 0 {
		result.Add(result, big.NewInt(1))
	}
	return result
}
