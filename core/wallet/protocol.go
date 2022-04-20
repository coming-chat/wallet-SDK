package wallet

type Util interface {
	// @param publicKey can start with 0x or not.
	EncodePublicKeyToAddress(publicKey string) (string, error)
	// @return publicKey that will start with 0x.
	DecodeAddressToPublicKey(address string) (string, error)

	IsValidAddress(address string) bool
}

type Account interface {
	// If the account generated using keystore, it will return empty
	// @return privateKey that will start with 0x.
	PrivateKey() string
	// @return publicKey that will start with 0x.
	PublicKey() string
	Address() string

	SignData(data []byte, password string) (string, error)
	SignHexData(hex string, password string) (string, error) // for sdk

	// Only available to accounts generated with keystore.
	// @return If the password is correct, will return nil
	CheckPassword(password string) error
}
