package basics

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// GetXOnlyPubKey returns the X-only public key
func GetXOnlyPubKey(pubKey string) ([]byte, error) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return nil, err
	}

	switch len(pubKeyBytes) {
	case schnorr.PubKeyBytesLen:
		return pubKeyBytes, nil

	case secp256k1.PubKeyBytesLenCompressed, secp256k1.PubKeyBytesLenUncompressed:
		return pubKeyBytes[1:33], nil

	default:
		return nil, fmt.Errorf("invalid public key length: %d", len(pubKeyBytes))
	}
}
