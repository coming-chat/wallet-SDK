package sui

import (
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
)

type SuiCatConfig struct {
	PackageId  string
	GlobalId   string
	ModuleName string

	whitelistId string
}

var SuiCatTestnetConfig = &SuiCatConfig{
	PackageId:  "0xc5b18811206c9ef35b516cd90f1736e7504f17fec147179298cc6851f2aa10a9",
	GlobalId:   "0x9876b64fad60ef76235f56c3221a4ee1aa891eaa3b86b10ed16195169c7c3e19",
	ModuleName: "suicat",
}

const (
	SuiCatFuncMint = "mint"
)

type SuiCatGlobalData struct {
	TotalMinted    int64  `json:"total_minted"`
	Supply         int64  `json:"supply"`
	PricePublic    string `json:"price_public"`
	PriceWhitelist string `json:"price_whitelist"`
	StartTimeMs    int64  `json:"start_time"`
	DurationMs     int64  `json:"duration"`
}

func (j *SuiCatGlobalData) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
}
func NewSuiCatGlobalDataWithJsonString(str string) (*SuiCatGlobalData, error) {
	var o SuiCatGlobalData
	err := base.FromJsonString(str, &o)
	return &o, err
}

type rawSuiCatGlobalData struct {
	Creator         string `json:"creator"`
	Duration        string `json:"duration"`
	MintedPublic    uint64 `json:"minted_public"`
	MintedWhitelist uint64 `json:"minted_whitelist"`
	Supply          uint64 `json:"supply"`
	Balance         string `json:"balance"`
	Beneficiary     string `json:"beneficiary"`
	PricePublic     string `json:"price_public"`
	PriceWhitelist  string `json:"price_whitelist"`
	StartTime       string `json:"start_time"`
	TeamReserve     uint64 `json:"team_reserve"`

	Whitelist struct {
		Fields struct {
			Id struct {
				Id string `json:"id"`
			} `json:"id"`
		} `json:"fields"`
	} `json:"whitelist"`
	Id struct {
		Id string `json:"id"`
	} `json:"id"`
	MintVault struct {
		Fields struct {
			Indexes []int `json:"indexes"`
		} `json:"fields"`
	} `json:"mint_vault"`
	// team_vault any
}

func (raw *rawSuiCatGlobalData) mapToPubData() *SuiCatGlobalData {
	if raw == nil {
		return nil
	}
	startTime, err := strconv.ParseInt(raw.StartTime, 10, 64)
	if err != nil {
		startTime = 0
		err = nil
	}
	duration, err := strconv.ParseInt(raw.Duration, 10, 64)
	if err != nil {
		duration = 0
		err = nil
	}
	return &SuiCatGlobalData{
		TotalMinted:    int64(len(raw.MintVault.Fields.Indexes)),
		Supply:         int64(raw.Supply),
		PricePublic:    raw.PricePublic,
		PriceWhitelist: raw.PriceWhitelist,
		StartTimeMs:    startTime,
		DurationMs:     duration,
	}
}
