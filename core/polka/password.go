package polka

import "encoding/json"

func (a *Account) CheckPassword(password string) error {
	if a.keystore == nil {
		return ErrNilKeystore
	}
	if a.keystore.CheckPassword(password) != nil {
		return ErrPassword
	}
	return nil
}

func CheckKeystorePassword(keystoreJson, password string) error {
	var keystore keystore
	err := json.Unmarshal([]byte(keystoreJson), &keystore)
	if err != nil {
		return err
	}
	if keystore.CheckPassword(password) != nil {
		return ErrPassword
	}
	return nil
}
