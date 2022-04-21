package wallet

type Util interface {
	// @param publicKey can start with 0x or not.
	EncodePublicKeyToAddress(publicKey string) (string, error)
	// @return publicKey that will start with 0x.
	DecodeAddressToPublicKey(address string) (string, error)

	IsValidAddress(address string) bool
}

type Account interface {
	// @return privateKey data
	PrivateKeyData() ([]byte, error)

	// @return privateKey string that will start with 0x.
	PrivateKey() (string, error)
	// @return publicKey string that will start with 0x.
	PublicKey() string
	// @return address string
	Address() string

	SignData(data []byte, password string) (string, error)
	SignHexData(hex string, password string) (string, error)
}
