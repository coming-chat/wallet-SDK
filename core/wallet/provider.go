package wallet

type WalletInfoProvider interface {
	Mnemonic(walletId string) string
	Keystore(walletId string) string
	Password(walletId string) string
	PrivateKey(walletId string) string
}

var InfoProvider WalletInfoProvider = nil

func readTypeAndValue(walletId string) (WalletType, string) {
	if InfoProvider == nil {
		return WalletTypeError, ""
	}
	if m := InfoProvider.Mnemonic(walletId); len(m) > 24 {
		return WalletTypeMnemonic, m
	} else if k := InfoProvider.Keystore(walletId); len(k) > 0 {
		return WalletTypeKeystore, k
	} else if p := InfoProvider.PrivateKey(walletId); len(p) > 0 {
		return WalletTypePrivateKey, p
	} else {
		return WalletTypeError, ""
	}
}

func readValue(walletId string, typ WalletType) (WalletType, string) {
	if InfoProvider == nil {
		return WalletTypeError, ""
	}
	switch typ {
	case WalletTypeMnemonic:
		return typ, InfoProvider.Mnemonic(walletId)
	case WalletTypeKeystore:
		return typ, InfoProvider.Keystore(walletId)
	case WalletTypePrivateKey:
		return typ, InfoProvider.PrivateKey(walletId)
	default:
		return WalletTypeError, ""
	}
}
