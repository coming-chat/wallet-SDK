package btc

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
	return &Brc20InscriptionPage{res}, nil
}
