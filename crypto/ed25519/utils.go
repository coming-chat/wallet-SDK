package ed25519

import "crypto/ed25519"

func IsValidSignature(publicKey, msg, signature []byte) bool {
	return ed25519.Verify(publicKey, msg, signature)
}
