package btc

import (
	"encoding/json"
	"errors"
	"fmt"
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

	// only has value if the content type starts with 'text/'
	ContentText string `json:"contentText"`
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

func (bi *Brc20Inscription) AsNFT() *base.NFT {
	return &base.NFT{
		Timestamp:  bi.Timestamp,
		HashString: bi.GenesisTransaction,

		Id:              bi.InscriptionId,
		Name:            fmt.Sprintf("Inscription %v", bi.InscriptionNumber),
		Image:           bi.Content,
		Standard:        "brc-20",
		Collection:      "",
		Descr:           "",
		ContractAddress: "",

		RelatedUrl: bi.Preview,
	}
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

type NFTPage struct {
	*inter.SdkPageable[*base.NFT]
}

func (bp *Brc20InscriptionPage) AsNFTPage() *NFTPage {
	nftArr := make([]*base.NFT, len(bp.Items))
	for i, bi := range bp.Items {
		nftArr[i] = bi.AsNFT()
	}
	return &NFTPage{SdkPageable: &inter.SdkPageable[*base.NFT]{
		TotalCount_:    bp.TotalCount_,
		CurrentCount_:  bp.CurrentCount(),
		CurrentCursor_: bp.CurrentCursor_,
		HasNextPage_:   bp.HasNextPage_,

		Items: nftArr,
	}}
}

type Brc20TransferableInscription struct {
	InscriptionId     string `json:"inscriptionId"`
	InscriptionNumber int64  `json:"inscriptionNumber"`
	Amount            string `json:"amount"`
	Ticker            string `json:"ticker"`
	Unconfirmed       bool   `json:"unconfirmed,omitempty"`
}

func (j *Brc20TransferableInscription) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
}
func NewBrc20TransferableInscriptionWithJsonString(str string) (*Brc20TransferableInscription, error) {
	var o Brc20TransferableInscription
	err := base.FromJsonString(str, &o)
	return &o, err
}

func (a *Brc20TransferableInscription) AsAny() *base.Any {
	return &base.Any{Value: a}
}
func AsBrc20TransferableInscription(a *base.Any) *Brc20TransferableInscription {
	if r, ok := a.Value.(*Brc20TransferableInscription); ok {
		return r
	}
	if r, ok := a.Value.(Brc20TransferableInscription); ok {
		return &r
	}
	return nil
}

type Brc20TransferableInscriptionPage struct {
	*inter.SdkPageable[*Brc20TransferableInscription]
}

func NewBrc20TransferableInscriptionPageWithJsonString(str string) (*Brc20TransferableInscriptionPage, error) {
	var o Brc20TransferableInscriptionPage
	err := base.FromJsonString(str, &o)
	return &o, err
}

type unisatTokenSummary struct {
	TransferableList []*Brc20TransferableInscription `json:"transferableList"`
	TokenBalance     *Brc20TokenBalance              `json:"tokenBalance"`
	// historyList, tokenInfo
}

type Brc20UTXO struct {
	Txid  string `json:"txid"`
	Index int64  `json:"index"`
}

func NewBrc20UTXO(txid string, index int64) *Brc20UTXO {
	return &Brc20UTXO{Txid: txid, Index: index}
}

type Brc20UTXOArray struct {
	inter.AnyArray[*Brc20UTXO]
}

func NewBrc20UTXOArray() *Brc20UTXOArray {
	return &Brc20UTXOArray{}
}

func (a *Brc20UTXOArray) UnmarshalJSON(data []byte) error {
	var out []*Brc20UTXO
	err := json.Unmarshal(data, &out)
	if err == nil {
		*a = Brc20UTXOArray{
			inter.AnyArray[*Brc20UTXO](out),
		}
	}
	return err
}

type Brc20CommitCustom struct {
	BaseTx string          `json:"baseTx"`
	Utxos  *Brc20UTXOArray `json:"utxos"`
}

func (a *Brc20CommitCustom) UnmarshalJSON(data []byte) error {
	var temp struct {
		BaseTx string          `json:"baseTx"`
		Utxos  *Brc20UTXOArray `json:"utxos"`
	}
	err := json.Unmarshal(data, &temp)
	if err == nil {
		(*a).BaseTx = temp.BaseTx
		(*a).Utxos = temp.Utxos
		return nil
	}

	var strs []string
	err = json.Unmarshal(data, &strs)
	if err == nil {
		if len(strs) < 3 || len(strs)%2 == 0 {
			return errors.New("invalid commit data")
		}
		utxos := NewBrc20UTXOArray()
		for i := 0; i < len(strs)/2; i++ {
			txid := strs[i*2+1]
			index, err := strconv.ParseInt(strs[i*2+2], 10, 64)
			if err != nil {
				return err
			}
			utxos.Append(NewBrc20UTXO(txid, index))
		}
		(*a).BaseTx = strs[0]
		(*a).Utxos = utxos
		return nil
	}
	return err
}

func (j *Brc20CommitCustom) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
}

func NewBrc20CommitCustomJson(json string) *Brc20CommitCustom {
	var o Brc20CommitCustom
	_ = base.FromJsonString(json, &o)
	return &o
}
