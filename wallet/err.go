package wallet

import "errors"

var (
	ErrNilKey       = errors.New("no mnemonic or private key")
	ErrNilWallet    = errors.New("no mnemonic or private key or keystore")
	ErrNilKeystore  = errors.New("no keystore")
	ErrNilMetadata  = errors.New("no metadata")
	ErrNotSigned    = errors.New("transaction not signed")
	ErrNoPublicKey  = errors.New("transaction no public key")
	ErrNilExtrinsic = errors.New("nil extrinsic")
	ErrAddress      = errors.New("err address")
	ErrPublicKey    = errors.New("err publicKey")
	ErrSeedOrPhrase = errors.New("invalid seed length")

	ErrNoEncrypted     = errors.New("no encrypted data to decode")
	ErrEncryptedLength = errors.New("encrypted length is less than 24")
	ErrInvalidParams   = errors.New("invalid injected scrypt params found")
	ErrSecretLength    = errors.New("secret length is not 32")
	ErrEncoded         = errors.New("encoded is nil")
	ErrPkcs8Header     = errors.New("invalid Pkcs8 header found in body")
	ErrPkcs8Divider    = errors.New("invalid Pkcs8 divider found in body")

	ErrNonPkcs8        = errors.New("unable to decode non-pkcs8 type")
	ErrNilPassword     = errors.New("password required to decode encrypted data")
	ErrNoEncryptedData = errors.New("no encrypted data available to decode")
	ErrKeystore        = errors.New("decoded public keys are not equal")

	ErrPassword = errors.New("password err")

	ErrNumber = errors.New("illegal number")
	ErrSign   = errors.New("sign panic error")
)
