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

func FetchBrc20Inscription(owner string, cursor string, pageSize int) (page *Brc20InscriptionPage, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)
	if cursor == "" {
		cursor = "0"
	}
	offset, err := strconv.ParseInt(cursor, 10, 64)
	if err != nil {
		return nil, errors.New("invalid cursor")
	}

	header := unisatRequestHeader()
	url := fmt.Sprintf("https://unisat.io/wallet-api-v4/address/inscriptions?address=%v&cursor=%v&size=%v", owner, offset, pageSize)
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
