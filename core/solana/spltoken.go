package solana

import (
	"context"
	"errors"
	"math/big"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/coming-chat/wallet-SDK/core/base"
)

type SPLToken struct {
	chain       *Chain
	MintAddress string
}

func NewSPLToken(chain *Chain, mintAddress string) *SPLToken {
	return &SPLToken{
		chain:       chain,
		MintAddress: mintAddress,
	}
}

// MARK - Implement the protocol Token

func (t *SPLToken) Chain() base.Chain {
	return t.chain
}

func (t *SPLToken) TokenInfo() (*base.TokenInfo, error) {
	return nil, errors.New("TODO func")
}

func (t *SPLToken) BalanceOfAddress(address string) (*base.Balance, error) {
	balances, err := t.TokenAccountOfAddress(address)
	if err != nil {
		return nil, err
	}
	if len(balances) == 0 {
		return nil, errors.New("the owner has not created the token account")
	}
	total := big.NewInt(0)
	for _, bal := range balances {
		total.Add(total, big.NewInt(int64(bal.Amount)))
	}
	return base.NewBalance(total.String()), nil
}
func (t *SPLToken) BalanceOfPublicKey(publicKey string) (*base.Balance, error) {
	addr, err := EncodePublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}
	return t.BalanceOfAddress(addr)
}
func (t *SPLToken) BalanceOfAccount(account base.Account) (*base.Balance, error) {
	return t.BalanceOfAddress(account.Address())
}

// BuildTransfer implements base.Token.
func (t *SPLToken) BuildTransfer(sender string, receiver string, amount string) (txn base.Transaction, err error) {
	return nil, errors.New("TODO func")
}

func (t *SPLToken) CanTransferAll() bool {
	return false
}
func (t *SPLToken) BuildTransferAll(sender string, receiver string) (txn base.Transaction, err error) {
	return nil, base.ErrUnsupportedFunction
}

// MARK - Help func

const (
	SPLAccountTypeRandom     = "Random"
	SPLAccountTypeAssociated = "Associated"
)

type TokenAccount struct {
	Address string
	Owner   string
	Amount  uint64

	AccountType string // "Random" or "Associated"
}

func (t *SPLToken) TokenAccountOfAddress(address string) (res []TokenAccount, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli := t.chain.client()

	// 先假设该地址是 token 账号地址
	tokenAcc, err := cli.GetTokenAccount(context.Background(), address)
	if err == nil {
		if t.MintAddress != tokenAcc.Mint.ToBase58() {
			return []TokenAccount{}, nil
		}
		return []TokenAccount{TransformTokenAccount(tokenAcc, address)}, nil
	}

	// 否则假设该地址是普通账号地址
	tokenAccs, err := cli.GetTokenAccountsByOwnerByMint(context.Background(), address, t.MintAddress)
	if err != nil {
		return nil, err
	}
	res = make([]TokenAccount, len(tokenAccs))
	for idx, acc := range tokenAccs {
		res[idx] = TransformTokenAccount(acc.TokenAccount, acc.PublicKey.ToBase58())
	}
	return res, nil
}

// TransformTokenAccount
func TransformTokenAccount(account token.TokenAccount, tokenAddress string) TokenAccount {
	res := TokenAccount{
		Address: tokenAddress,
		Owner:   account.Owner.ToBase58(),
		Amount:  account.Amount,
	}
	ata, _, err := common.FindAssociatedTokenAddress(account.Owner, account.Mint)
	if err == nil && ata.ToBase58() == tokenAddress {
		res.AccountType = SPLAccountTypeAssociated
	} else {
		res.AccountType = SPLAccountTypeRandom
	}
	return res
}

func (t *SPLToken) HasCreated(ownerAddress string) (b *base.OptionalBool, err error) {
	accounts, err := t.TokenAccountOfAddress(ownerAddress)
	if err != nil {
		return nil, err
	}
	return base.NewOptionalBool(len(accounts) > 0), nil
}

func (t *SPLToken) CreateTokenAccount(ownerAddress string, signerAddress string) (txn *Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	ownerPubkey := common.PublicKeyFromString(ownerAddress)
	tokenPubkey := common.PublicKeyFromString(t.MintAddress)
	signerPubkey := common.PublicKeyFromString(signerAddress)
	associateAccount, _, err := common.FindAssociatedTokenAddress(ownerPubkey, tokenPubkey)
	if err != nil {
		return
	}

	cli := t.chain.client()
	latestBlock, err := cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return
	}
	msg := types.NewMessage(types.NewMessageParam{
		FeePayer:        signerPubkey,
		RecentBlockhash: latestBlock.Blockhash,
		Instructions: []types.Instruction{
			associated_token_account.CreateAssociatedTokenAccount(associated_token_account.CreateAssociatedTokenAccountParam{
				Funder:                 signerPubkey,
				Owner:                  ownerPubkey,
				Mint:                   tokenPubkey,
				AssociatedTokenAccount: associateAccount,
			}),
		},
	})
	return &Transaction{
		Message: msg,
	}, nil
}
