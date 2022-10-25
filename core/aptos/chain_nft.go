package aptos

import (
	"encoding/hex"
	"errors"

	aptosnft "github.com/coming-chat/go-aptos/nft"
	txbuilder "github.com/coming-chat/go-aptos/transaction_builder"
	"github.com/coming-chat/wallet-SDK/core/base"
)

func (c *Chain) OfferTokenTransactionNFT(sender *Account, receiver string, nft *base.NFT) (*base.OptionalString, error) {
	return c.tokenTxnBuild(sender, receiver, nft.ContractAddress, nft.Collection, nft.Name, "offer")
}

/** build transaction that offer token
 * @param sender the transferring token owner
 * @param receiver the token receiver
 * @param creator the token creator
 * @param collection the token's collection name
 * @param name the token's name
 * @return the offer token raw transaction that signed by sender.
 */
func (c *Chain) OfferTokenTransactionParams(sender *Account, receiver, creator, collection, name string) (*base.OptionalString, error) {
	return c.tokenTxnBuild(sender, receiver, creator, collection, name, "offer")
}

func (c *Chain) CancelTokenOffer(sender *Account, receiver, creator, collection, name string) (*base.OptionalString, error) {
	return c.tokenTxnBuild(sender, receiver, creator, collection, name, "cancelOffer")
}

/** build transaction that claim token, the nft info will be obtaining through offer hash.
 * @param receiver the token receiver
 * @param offerHash the submitted hash of the transaction that offer the token
 * @return the claim token raw transaction that signed by receiver.
 */
func (c *Chain) ClaimTokenFromHash(receiver *Account, offerHash string) (signedTxn *base.OptionalString, err error) {
	client, err := c.client()
	if err != nil {
		return
	}
	offeredTxn, err := client.GetTransactionByHash(offerHash)
	if err != nil {
		return
	}
	if !offeredTxn.Success {
		if offeredTxn.VmStatus != "" {
			return nil, errors.New("Claim failed, the offer transaction failed: " + offeredTxn.VmStatus)
		} else {
			return nil, errors.New("Claim failed, the offer transaction may still be pending.")
		}
		return //lint:ignore
	}
	if offeredTxn.Payload.Function != "0x3::token_transfers::offer_script" {
		return nil, errors.New("Claim failed, the given hash is not an offer token transaction")
	}
	arguments := offeredTxn.Payload.Arguments
	if len(arguments) < 4 {
		return nil, errors.New("Claim failed, offer params invalid.")
	}
	receiverString := arguments[0].(string)
	if receiver.Address() != receiverString {
		return nil, errors.New("Claim failed, this token is not offer to the receiver.")
	}
	creator := arguments[1].(string)
	collectionName := arguments[2].(string)
	tokenName := arguments[3].(string)
	nftSender := offeredTxn.Sender

	return c.tokenTxnBuild(receiver, nftSender, creator, collectionName, tokenName, "claim")
}

/** build transaction that claim token
 * @param receiver the token receiver
 * @param sender the transferred token owner
 * @param creator the token creator
 * @param collection the token's collection name
 * @param name the token's name
 * @return the claim token raw transaction that signed by receiver.
 */
func (c *Chain) ClaimTokenTransactionParams(receiver *Account, sender, creator, collection, name string) (signedTxn *base.OptionalString, err error) {
	return c.tokenTxnBuild(receiver, sender, creator, collection, name, "claim")
}

// @param action enum of `offer`, `claim`, `cancelOffer`
func (c *Chain) tokenTxnBuild(signer *Account, senderOrReceiver, creator, collection, name string, action string) (signedTxn *base.OptionalString, err error) {
	defer base.CatchPanicAndMapToBasicError(&err)

	senderOrReceiverAddress, err := txbuilder.NewAccountAddressFromHex(senderOrReceiver)
	if err != nil {
		return
	}
	creatorAddress, err := txbuilder.NewAccountAddressFromHex(creator)
	if err != nil {
		return
	}
	if collection == "" || name == "" {
		return nil, errors.New("The `collection` and `name` cannot be empty.")
	}
	builder, err := aptosnft.NewNFTPayloadBuilder()
	if err != nil {
		return
	}

	var payload txbuilder.TransactionPayload
	switch action {
	case "offer":
		payload, err = builder.OfferToken(*senderOrReceiverAddress, *creatorAddress, collection, name, 1, 0)
	case "claim":
		payload, err = builder.ClaimToken(*senderOrReceiverAddress, *creatorAddress, collection, name, 0)
	case "cancelOffer":
		payload, err = builder.CancelTokenOffer(*senderOrReceiverAddress, *creatorAddress, collection, name, 0)
	default:
		return nil, errors.New("Invalid token action: " + action)
	}
	if err != nil {
		return
	}

	txn, err := c.createTransactionFromPayloadBCS(signer, payload)
	if err != nil {
		return
	}
	signedBytes, err := txbuilder.GenerateBCSTransaction(signer.account, txn)
	if err != nil {
		return
	}
	hexString := "0x" + hex.EncodeToString(signedBytes)
	return &base.OptionalString{Value: hexString}, nil
}
