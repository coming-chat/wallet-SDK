package polka

import "testing"

func TestChain_EstimateFeeForTransaction(t *testing.T) {
	type transactionCreator func(tx *Tx) (*Transaction, error)
	balanceCreator := func(address, amount string) transactionCreator {
		return func(tx *Tx) (*Transaction, error) {
			return tx.NewBalanceTransferTx(address, amount)
		}
	}
	tests := []struct {
		name    string
		rpcinfo rpcInfo
		creator transactionCreator
		wantErr bool
	}{
		{
			name:    "chainx prod PCX very much",
			rpcinfo: rpcs.chainxProd,
			creator: balanceCreator(accountCase.address44, "100000000000000000000000000000000"),
		},
		{
			name:    "chainx pre PCX 0",
			rpcinfo: rpcs.chainxTest,
			creator: balanceCreator(accountCase.address44, "0"),
		},
		{
			name:    "chainx prod XBTC transfer",
			rpcinfo: rpcs.chainxProd,
			creator: func(tx *Tx) (*Transaction, error) {
				return tx.NewXAssetsTransferTx(accountCase.address0, "1000000")
			},
		},
		{
			name:    "minix pre MINI very very much",
			rpcinfo: rpcs.minixTest,
			creator: balanceCreator(accountCase.address2, "999999999999999999999999999999999999999999999999"),
			wantErr: true,
		},
		{
			name:    "minix pre MINI 0",
			rpcinfo: rpcs.minixTest,
			creator: balanceCreator(accountCase.address2, "0"),
		},
		{
			name:    "minix prod CID transfer",
			rpcinfo: rpcs.minixProd,
			creator: func(tx *Tx) (*Transaction, error) {
				return tx.NewComingNftTransferTx(accountCase.address0, 123456)
			},
		},
		{
			name:    "sherpax prod KSX transfer",
			rpcinfo: rpcs.sherpaxProd,
			creator: balanceCreator(accountCase.address2, "1000"),
		},
		{
			name:    "polkadot prod DOT very much",
			rpcinfo: rpcs.polkadot,
			creator: balanceCreator(accountCase.address2, "99999999999999999999999999999999"),
		},
		{
			name:    "polkadot prod DOT 0",
			rpcinfo: rpcs.polkadot,
			creator: balanceCreator(accountCase.address2, "0"),
		},
		{
			name:    "kusama prod KSM transfer",
			rpcinfo: rpcs.kusama,
			creator: balanceCreator(accountCase.address2, "1000"),
		},
		{
			name:    "error amount",
			rpcinfo: rpcs.minixProd,
			creator: balanceCreator(accountCase.address2, "abcdddddd"),
			wantErr: true,
		},
		{
			name:    "error address",
			rpcinfo: rpcs.sherpaxTest,
			creator: balanceCreator(accountCase.address2+"xx", "1000"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := tt.rpcinfo.Chain()
			tx, err := c.GetTx()
			if err != nil {
				t.Errorf("EstimateFeeForTransaction() get tx error = %v", err)
				return
			}
			transaction, err := tt.creator(tx)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("EstimateFeeForTransaction() get new transaction error = %v", err)
				}
				return
			}
			gotS, err := c.EstimateFeeForTransaction(transaction)
			if (err != nil) != tt.wantErr {
				t.Errorf("EstimateFeeForTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				t.Logf("Got fee: %v, But we cannot assert the fee.", gotS)
			}
		})
	}
}
