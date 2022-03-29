package httpUtil

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

var (
	httpClient = &http.Client{
		Timeout: 50 * time.Second,
	}
)

type Res struct {
	Body []byte
	Code int
}

func Request(method, url string, header map[string]string, body []byte) (*Res, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}

	resp, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	res, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return nil, errRead
	}
	return &Res{
		Body: res,
		Code: resp.StatusCode,
	}, nil
}

func RequestSync(method, url string, header map[string]string, body []byte, httpWg *sync.WaitGroup, resp *Res, err error) {
	response, err := Request(method, url, header, body)
	if err == nil {
		*resp = *response
	}
	httpWg.Done()
}
