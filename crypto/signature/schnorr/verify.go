package schnorr

import (
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// VerifySignature verifies the provided signature against the given hash and public key
func VerifySignature(sig []byte, sigHash []byte, pubKeyBytes []byte) bool {
	signature, err := schnorr.ParseSignature(sig)
	if err != nil {
		return false
	}

	pubKey, err := schnorr.ParsePubKey(pubKeyBytes)
	if err != nil {
		return false
	}

	return signature.Verify(sigHash, pubKey)
}
