package polka

import "encoding/json"

func (a *Account) CheckPassword(password string) error {
	if a.keystore == nil {
		return ErrNilKeystore
	}
	return a.keystore.CheckPassword(password)
}

func IsValidKeystore(keysotreJson, password string) bool {
	var keyStore keystore
	err := json.Unmarshal([]byte(keysotreJson), &keyStore)
	if err != nil {
		return false
	}
	if err = keyStore.CheckPassword(password); err != nil {
		return false
	}
	return true
}
