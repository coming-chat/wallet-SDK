package btc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/core/base/inter"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

// FetchBrc20Inscription
// @param cursor start from 0
func (c *Chain) FetchBrc20Inscription(owner string, cursor string, pageSize int) (page *Brc20InscriptionPage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}
	host, err := unisatHost(c.Chainnet)
	if err != nil {
		return nil, err
	}

	header := unisatRequestHeader()
	url := fmt.Sprintf("%v/wallet-api-v4/address/inscriptions?address=%v&cursor=%v&size=%v", host, owner, offset, pageSize)
	resp, err := httpUtil.Request(http.MethodGet, url, header, nil)
	if err != nil {
		return
	}
	var rawPage rawBrc20InscriptionPage
	if err = decodeUnisatResponseV4(*resp, &rawPage); err != nil {
		return
	}

	res := rawPage.MapToSdkPage(int(offset), pageSize)
	batchFetchContentText(res.Items)
	return &Brc20InscriptionPage{res}, nil
}

func batchFetchContentText(inscriptions []*Brc20Inscription) {
	listAny := make([]any, len(inscriptions))
	for i, v := range inscriptions {
		listAny[i] = v
	}
	base.MapListConcurrent(listAny, 10, func(i interface{}) (interface{}, error) {
		inscription, ok := i.(*Brc20Inscription)
		if !ok {
			return nil, nil
		}
		if !strings.HasPrefix(inscription.ContentType, "text/plain") {
			return nil, nil
		}
		res, err := httpUtil.Request(http.MethodGet, inscription.Content, nil, nil)
		if err != nil {
			return nil, nil
		}
		inscription.ContentText = string(res.Body)
		return nil, nil
	})
}

func (c *Chain) FetchBrc20TransferableInscription(owner string, ticker string) (page *Brc20TransferableInscriptionPage, err error) {
	summary, err := c.fetchTokenSummary(owner, ticker)
	if err != nil {
		return
	}
	unconfirmedList, err := c.fetchBrc20UnconfirmedTransferableInscription(owner, ticker)
	if err != nil {
		unconfirmedList = []*Brc20TransferableInscription{}
		err = nil
	}

	confirmedList := summary.TransferableList
	confirmedList = append(confirmedList, unconfirmedList...)
	return &Brc20TransferableInscriptionPage{
		SdkPageable: &inter.SdkPageable[*Brc20TransferableInscription]{
			TotalCount_:    len(confirmedList),
			CurrentCount_:  len(confirmedList),
			CurrentCursor_: "end",
			HasNextPage_:   false,
			Items:          confirmedList,
		},
	}, nil
}

func (c *Chain) fetchBrc20UnconfirmedTransferableInscription(owner string, ticker string) (arr []*Brc20TransferableInscription, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	host, err := unisatHost(c.Chainnet)
	if err != nil {
		return nil, err
	}

	header := unisatRequestHeader()
	url := fmt.Sprintf("%v/wallet-api-v4/address/inscriptions?address=%v&cursor=%v&size=%v", host, owner, 0, 50)
	resp, err := httpUtil.Request(http.MethodGet, url, header, nil)
	if err != nil {
		return
	}
	var rawPage rawBrc20InscriptionPage
	if err = decodeUnisatResponseV4(*resp, &rawPage); err != nil {
		return
	}

	var unconfirmedInscription = []interface{}{}
	for _, inscription := range rawPage.List {
		if inscription.Timestamp <= 10 {
			unconfirmedInscription = append(unconfirmedInscription, inscription)
		}
	}
	unconfirmedTransferableInscription, err := base.MapListConcurrent(unconfirmedInscription, 10, func(i interface{}) (interface{}, error) {
		inscription := i.(*Brc20Inscription)
		body, err := httpUtil.Get(inscription.Content, nil)
		if err != nil {
			return 0, nil
		}
		var obj struct {
			// P    string `json:"p"`
			// Op   string `json:"op"`
			Tick string `json:"tick"`
			Amt  string `json:"amt"`
		}
		err = json.Unmarshal(body, &obj)
		if err != nil {
			return 0, nil
		}
		if obj.Tick != ticker {
			return 0, nil
		}
		return &Brc20TransferableInscription{
			InscriptionId:     inscription.InscriptionId,
			InscriptionNumber: inscription.InscriptionNumber,
			Amount:            obj.Amt,
			Ticker:            obj.Tick,
			Unconfirmed:       true,
		}, nil
	})
	if err != nil {
		return nil, err
	}

	arr = []*Brc20TransferableInscription{}
	for _, inscription := range unconfirmedTransferableInscription {
		if item, ok := inscription.(*Brc20TransferableInscription); ok {
			arr = append(arr, item)
		}
	}
	return arr, nil
}
