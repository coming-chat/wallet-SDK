package btc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type Brc20Token struct {
	Ticker string
}

func NewBrc20Token(ticker string) *Brc20Token {
	return &Brc20Token{
		Ticker: ticker,
	}
}

func (t *Brc20Token) Chain() base.Chain {
	return nil
}

func (t *Brc20Token) TokenInfo() (*base.TokenInfo, error) {
	info, err := t.FullTokenInfo()
	if err != nil {
		return nil, err
	}
	return &base.TokenInfo{
		Name:    info.Ticker,
		Symbol:  info.Ticker,
		Decimal: info.Decimal,
	}, nil
}

func (t *Brc20Token) BalanceOfAddress(address string) (*base.Balance, error) {
	return nil, base.ErrUnsupportedFunction
}
func (t *Brc20Token) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	return nil, base.ErrUnsupportedFunction
}
func (t *Brc20Token) BalanceOfAccount(account Account) (*base.Balance, error) {
	return nil, base.ErrUnsupportedFunction
}

func (t *Brc20Token) BuildTransfer(sender, receiver, amount string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}

// Before invoking this method, it is best to check `CanTransferAll()`
func (t *Brc20Token) CanTransferAll() bool {
	return false
}
func (t *Brc20Token) BuildTransferAll(sender, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}

var brc20InfoCache map[string]*Brc20TokenInfo

func (t *Brc20Token) FullTokenInfo() (info *Brc20TokenInfo, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	key := strings.ToLower(t.Ticker)
	if brc20InfoCache == nil {
		brc20InfoCache = make(map[string]*Brc20TokenInfo)
	} else if cache, exists := brc20InfoCache[key]; exists {
		return cache, nil
	}
	header := unisatRequestHeader()
	url := fmt.Sprintf("https://unisat.io/brc20-api-v2/brc20/%v/info", t.Ticker)
	resp, err := httpUtil.Request(http.MethodGet, url, header, nil)
	if err != nil {
		return
	}
	info = &Brc20TokenInfo{}
	if err = decodeUnisatResponseV2(*resp, info); err != nil {
		return
	}
	brc20InfoCache[key] = info
	return
}

type Brc20TokenInfo struct {
	Ticker                 string `json:"ticker"`
	HoldersCount           int64  `json:"holdersCount"`
	HistoryCount           int64  `json:"historyCount"`
	InscriptionNumber      int64  `json:"inscriptionNumber"`
	InscriptionId          string `json:"inscriptionId"`
	Max                    string `json:"max"`
	Limit                  string `json:"limit"`
	Minted                 string `json:"minted"`
	TotalMinted            string `json:"totalMinted"`
	ConfirmedMinted        string `json:"confirmedMinted"`
	ConfirmedMinted1h      string `json:"confirmedMinted1h"`
	ConfirmedMinted24h     string `json:"confirmedMinted24h"`
	MintTimes              int64  `json:"mintTimes"`
	Decimal                int16  `json:"decimal"`
	Creator                string `json:"creator"`
	Txid                   string `json:"txid"`
	DeployHeight           int64  `json:"deployHeight"`
	DeployBlocktime        int64  `json:"deployBlocktime"`
	CompleteHeight         int64  `json:"completeHeight"`
	CompleteBlocktime      int64  `json:"completeBlocktime"`
	InscriptionNumberStart int64  `json:"inscriptionNumberStart"`
	InscriptionNumberEnd   int64  `json:"inscriptionNumberEnd"`
}

func (j *Brc20TokenInfo) JsonString() (*base.OptionalString, error) {
	return base.JsonString(j)
}
func NewBrc20TokenInfoWithJsonString(str string) (*Brc20TokenInfo, error) {
	var o Brc20TokenInfo
	err := base.FromJsonString(str, &o)
	return &o, err
}

func FetchBrc20TokenBalance(owner string, cursor string, pageSize int) (page *Brc20TokenBalancePage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}

	header := unisatRequestHeader()
	url := fmt.Sprintf("https://unisat.io/wallet-api-v4/brc20/tokens?address=%v&cursor=%v&size=%v", owner, offset, pageSize)
	resp, err := httpUtil.Request(http.MethodGet, url, header, nil)
	if err != nil {
		return
	}
	var rawPage rawBrc20TokenBalancePage
	if err = decodeUnisatResponseV4(*resp, &rawPage); err != nil {
		return
	}

	res := rawPage.MapToSdkPage(int(offset), pageSize)
	return &Brc20TokenBalancePage{res}, nil
}
