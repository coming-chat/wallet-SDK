package btc

import (
	"strconv"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
)

type unisatRawPage[T any] struct {
	List  []T `json:"list"`
	Total int `json:"total"`
}

// If your output type is not the T, you should implement your custom `MapToSdkPage`
func (p *unisatRawPage[T]) MapToSdkPage(offset, size int) *inter.SdkPageable[T] {
	fetchedCount := offset + size
	return &inter.SdkPageable[T]{
		TotalCount_:    p.Total,
		CurrentCount_:  len(p.List),
		CurrentCursor_: strconv.FormatInt(int64(fetchedCount), 10),
		HasNextPage_:   p.Total > fetchedCount,
		Items:          p.List,
	}
}

// - MARK -

type Brc20TokenBalance struct {
	Ticker              string `json:"ticker"`              //: "zbit",
	OverallBalance      string `json:"overallBalance"`      //: "0",
	TransferableBalance string `json:"transferableBalance"` //: "0",
	AvailableBalance    string `json:"availableBalance"`    //: "0"
}

func (j *Brc20TokenBalance) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
}
func NewBrc20TokenBalanceWithJsonString(str string) (*Brc20TokenBalance, error) {
	var o Brc20TokenBalance
	err := base.FromJsonString(str, &o)
	return &o, err
}

func (a *Brc20TokenBalance) AsAny() *base.Any {
	return &base.Any{Value: a}
}
func AsBrc20TokenBalance(a *base.Any) *Brc20TokenBalance {
	if r, ok := a.Value.(*Brc20TokenBalance); ok {
		return r
	}
	if r, ok := a.Value.(Brc20TokenBalance); ok {
		return &r
	}
	return nil
}

type Brc20TokenBalancePage struct {
	*inter.SdkPageable[*Brc20TokenBalance]
}

func NewBrc20TokenBalancePageWithJsonString(str string) (*Brc20TokenBalancePage, error) {
	var o Brc20TokenBalancePage
	err := base.FromJsonString(str, &o)
	return &o, err
}

type rawBrc20TokenBalancePage struct {
	unisatRawPage[*Brc20TokenBalance]
}

// - MARK -

type Brc20Inscription struct {
	InscriptionId      string `json:"inscriptionId"`
	InscriptionNumber  int64  `json:"inscriptionNumber"`
	Address            string `json:"address"`
	OutputValue        int64  `json:"outputValue"`
	Preview            string `json:"preview"`
	Content            string `json:"content"`
	ContentLength      int64  `json:"contentLength"`
	ContentType        string `json:"contentType"`
	ContentBody        string `json:"contentBody"`
	Timestamp          int64  `json:"timestamp"`
	GenesisTransaction string `json:"genesisTransaction"`
	Location           string `json:"location"`
	Output             string `json:"output"`
	// Offset             int64  `json:"offset"`
}

func (j *Brc20Inscription) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
}
func NewBrc20InscriptionWithJsonString(str string) (*Brc20Inscription, error) {
	var o Brc20Inscription
	err := base.FromJsonString(str, &o)
	return &o, err
}

func (a *Brc20Inscription) AsAny() *base.Any {
	return &base.Any{Value: a}
}
func AsBrc20Inscription(a *base.Any) *Brc20Inscription {
	if r, ok := a.Value.(*Brc20Inscription); ok {
		return r
	}
	if r, ok := a.Value.(Brc20Inscription); ok {
		return &r
	}
	return nil
}

type Brc20InscriptionPage struct {
	*inter.SdkPageable[*Brc20Inscription]
}

func NewBrc20InscriptionPageWithJsonString(str string) (*Brc20InscriptionPage, error) {
	var o Brc20InscriptionPage
	err := base.FromJsonString(str, &o)
	return &o, err
}

type rawBrc20InscriptionPage struct {
	unisatRawPage[*Brc20Inscription]
}
