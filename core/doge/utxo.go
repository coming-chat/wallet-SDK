package doge

import (
	"fmt"
	"math/big"
)

type UTXO struct {
	// input from like
	// https://api.blockcypher.com/v1/doge/main/addrs/D8aDCsK4TA9NYhmwiqw1BjZ4CP8LQ814Ea?limit=5&unspentOnly=true
	Txid  string   `json:"tx_hash"`
	Index int      `json:"tx_output_n"`
	Value *big.Int `json:"value"`
}

func (u *UTXO) MarshalJSON() ([]byte, error) {
	// output to like
	// https://yapi.coming.chat/project/13/interface/api/578
	jsonString := fmt.Sprintf("{\"txid\":\"%v\",\"index\":%v,\"value\":\"%v\",\"prevTx\":\"\"}", u.Txid, u.Index, u.Value.String())
	return []byte(jsonString), nil
}

type UTXOList struct {
	Utxos []*UTXO `json:"txrefs"`
}

type SDKUTXOList struct {
	Txids      []*UTXO `json:"txids"`
	FastestFee int     `json:"fastestFee"`
}
