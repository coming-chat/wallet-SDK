package solana

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/program/metaplex/token_metadata"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/coming-chat/wallet-SDK/core/base"
	"github.com/coming-chat/wallet-SDK/pkg/httpUtil"
)

type SPLToken struct {
	chain       *Chain
	MintAddress string
}

func NewSPLToken(chain *Chain, mintAddress string) (*SPLToken, error) {
	if !IsValidAddress(mintAddress) {
		return nil, base.ErrInvalidAddress
	}
	return &SPLToken{
		chain:       chain,
		MintAddress: mintAddress,
	}, nil
}

// MARK - Implement the protocol Token

func (t *SPLToken) Chain() base.Chain {
	return t.chain
}

func (t *SPLToken) TokenInfo() (*base.TokenInfo, error) {
	return t.TokenInfoChainid(101) // mainnet chainid
}

// TokenInfoChainid
// @param chainId mainnet 101, testnet 102, devnet 103
func (t *SPLToken) TokenInfoChainid(chainId int) (info *base.TokenInfo, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	fetchApiInfo := func(chainId int) (*base.TokenInfo, error) {
		header := map[string]string{
			"content-type": "application/json",
		}
		bodyString := fmt.Sprintf(`{"addresses":["%v"]}`, t.MintAddress)
		bodyBytes := []byte(bodyString)

		url := fmt.Sprintf("https://token-list-api.solana.cloud/v1/mints?chainId=%v", chainId)
		res, err := httpUtil.Request(http.MethodPost, url, header, bodyBytes)
		if err != nil {
			return nil, err
		}
		if res.Code != 200 && res.Code != 201 {
			return nil, fmt.Errorf("fetch token info error, code: %v", res.Code)
		}

		var tokeninfo struct {
			Content []struct {
				Name     string `json:"name"`
				Symbol   string `json:"symbol"`
				Decimals int16  `json:"decimals"`
				// Address  string `json:"address"`
				// LogoURI  string `json:"logoURI"`
			} `json:"content"`
		}
		err = json.Unmarshal(res.Body, &tokeninfo)
		if err != nil {
			return nil, err
		}
		if len(tokeninfo.Content) == 0 {
			return nil, fmt.Errorf("not found")
		}
		tmp := tokeninfo.Content[0]
		return &base.TokenInfo{
			Name:    tmp.Name,
			Symbol:  tmp.Symbol,
			Decimal: tmp.Decimals,
		}, nil
	}
	info, err = fetchApiInfo(chainId) // mainnet chainid
	if err != nil {
		// pass, continue other fetch way
	} else {
		return info, nil
	}

	info = &base.TokenInfo{
		Name:   t.MintAddress,
		Symbol: t.MintAddress,
	}
	cli := t.chain.Client()
	updateMetadata := func() {
		tokenPub := common.PublicKeyFromString(t.MintAddress)
		metaPubkey, err := token_metadata.GetTokenMetaPubkey(tokenPub)
		if err != nil {
			return
		}
		accoutInfo, err := cli.GetAccountInfo(context.Background(), metaPubkey.ToBase58())
		if err != nil {
			return
		}
		metadata, err := token_metadata.MetadataDeserialize(accoutInfo.Data)
		if err != nil {
			return
		}
		info.Name = metadata.Data.Name
		info.Symbol = metadata.Data.Symbol
	}
	updateMetadata()

	// fetch decimal
	amt, err := cli.GetTokenSupply(context.Background(), t.MintAddress)
	if err != nil {
		return nil, err
	}
	info.Decimal = int16(amt.Decimals)
	return info, nil
}

func (t *SPLToken) BalanceOfAddress(address string) (*base.Balance, error) {
	balances, _, err := t.TokenAccountOfAddress(address)
	if err != nil {
		return nil, err
	}
	if len(balances) == 0 {
		return nil, ErrNoTokenAccount
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

// BuildTransfer
// This method will automatically create an token account for the receiver if receiver does not own it.
func (t *SPLToken) BuildTransfer(sender string, receiver string, amount string) (txn base.Transaction, err error) {
	return t.BuildTransferAuto(sender, receiver, amount, false, true)
}

// CanTransferAll
// Available
func (t *SPLToken) CanTransferAll() bool {
	return true
}

// BuildTransferAll
// This method will automatically create an token account for the receiver if receiver does not own it.
func (t *SPLToken) BuildTransferAll(sender string, receiver string) (txn base.Transaction, err error) {
	return t.BuildTransferAuto(sender, receiver, "0", true, true)
}

// BuildTransferAuto
// @param transferAll if true will transfer all balance, else transfer the amount
// @param autoCreateAccount if true will auto create token account for receiver, else throw error if receiver no has token account
func (t *SPLToken) BuildTransferAuto(sender, receiver, amount string, transferAll bool, autoCreateAccount bool) (txn base.Transaction, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	if !IsValidAddress(sender) || !IsValidAddress(receiver) {
		return nil, base.ErrInvalidAccountAddress
	}
	if transferAll {
		amount = "1"
	}
	amountInt, ok := big.NewInt(0).SetString(amount, 10)
	if !ok || amountInt.Cmp(big.NewInt(0)) <= 0 {
		return nil, base.ErrInvalidAmount
	}

	// 确定好要转账的主账号信息
	senderAccounts, _, err := t.TokenAccountOfAddress(sender)
	if err != nil {
		return nil, err
	}
	total := big.NewInt(0)
	for _, acc := range senderAccounts {
		total.Add(total, big.NewInt(int64(acc.Amount)))
	}
	if total.Cmp(amountInt) <= 0 {
		return nil, base.ErrInsufficientBalance
	}
	realSenderPubkey := common.PublicKeyFromString(senderAccounts[0].Owner)

	instructions := make([]types.Instruction, 0)

	// 确定好要接受转账的真实账号地址
	var realReceiverPubkey common.PublicKey
	receiverAccounts, unmatchToken, err := t.TokenAccountOfAddress(receiver)
	if err != nil {
		return nil, err
	}
	if unmatchToken {
		return nil, errors.New("the receiver's token account does not match the token type to be transferred")
	}
	if len(receiverAccounts) == 0 {
		if !autoCreateAccount {
			return nil, ErrNoTokenAccount
		}
		// 如果需要创建账号，则创建固定不变的关联账号
		receiverPubkey := common.PublicKeyFromString(receiver)
		mintPubkey := common.PublicKeyFromString(t.MintAddress)
		realReceiverPubkey, _, err = common.FindAssociatedTokenAddress(receiverPubkey, mintPubkey)
		if err != nil {
			return nil, err
		}
		instructions = append(instructions,
			associated_token_account.CreateAssociatedTokenAccount(associated_token_account.CreateAssociatedTokenAccountParam{
				Funder:                 realSenderPubkey,
				Owner:                  receiverPubkey,
				Mint:                   mintPubkey,
				AssociatedTokenAccount: realReceiverPubkey,
			}),
		)
	} else {
		addr := receiverAccounts[0].Address
		for _, acc := range receiverAccounts {
			if acc.AccountType == SPLAccountTypeAssociated { // 优先转账到关联账号
				addr = acc.Address
				break
			}
		}
		realReceiverPubkey = common.PublicKeyFromString(addr)
	}

	// 构建转账命令
	appendTransferInstruction := func(from string, amount uint64) {
		if amount != 0 {
			ins := token.Transfer(token.TransferParam{
				From:    common.PublicKeyFromString(from),
				To:      realReceiverPubkey,
				Auth:    realSenderPubkey,
				Signers: []common.PublicKey{},
				Amount:  amount,
			})
			instructions = append(instructions, ins)
		}
	}
	if transferAll {
		for _, acc := range senderAccounts {
			appendTransferInstruction(acc.Address, acc.Amount)
		}
	} else {
		needAmt := amountInt
		for _, acc := range senderAccounts {
			amt := big.NewInt(int64(acc.Amount))
			if needAmt.Cmp(amt) > 0 { // needAmt > amt
				appendTransferInstruction(acc.Address, amt.Uint64())
				needAmt = needAmt.Sub(needAmt, amt)
			} else {
				appendTransferInstruction(acc.Address, needAmt.Uint64())
				break
			}
		}
	}

	cli := t.chain.Client()
	lastestBlock, err := cli.GetLatestBlockhash(context.Background())
	if err != nil {
		return nil, err
	}
	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        realSenderPubkey,
		RecentBlockhash: lastestBlock.Blockhash,
		Instructions:    instructions,
	})
	return &Transaction{
		Message: message,
	}, nil
}

// MARK - Help func

const (
	SPLAccountTypeRandom     = "Random"
	SPLAccountTypeAssociated = "Associated"
)

type TokenAccount struct {
	Amount  uint64
	Address string
	Owner   string

	AccountType string // "Random" or "Associated"
}

func (t *SPLToken) TokenAccountOfAddress(address string) (res []TokenAccount, unmatchToken bool, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	cli := t.chain.Client()

	// 先假设该地址是 token 账号地址
	tokenAcc, err := cli.GetTokenAccount(context.Background(), address)
	if err != nil {
		// pass
	} else {
		if t.MintAddress != tokenAcc.Mint.ToBase58() {
			return []TokenAccount{}, true, nil
		}
		return []TokenAccount{TransformTokenAccount(tokenAcc, address)}, false, nil
	}

	// 否则假设该地址是普通账号地址
	tokenAccs, err := cli.GetTokenAccountsByOwnerByMint(context.Background(), address, t.MintAddress)
	if err != nil {
		return nil, false, err
	}
	res = make([]TokenAccount, len(tokenAccs))
	for idx, acc := range tokenAccs {
		res[idx] = TransformTokenAccount(acc.TokenAccount, acc.PublicKey.ToBase58())
	}
	return res, false, nil
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
	accounts, _, err := t.TokenAccountOfAddress(ownerAddress)
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

	cli := t.chain.Client()
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
