package solana

import "errors"

var (
	ErrNoTokenAccount = errors.New("the owner has not created the token account")
)

func IsNoTokenAccountError(err error) bool {
	return err.Error() == ErrNoTokenAccount.Error()
}
