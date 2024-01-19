package protocol

import (
	"encoding/hex"
	"fmt"

	"github.com/tidwall/gjson"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

// ValidateSymbol validates if the given symbol satisfies the rules
func ValidateSymbol(symbol string) error {
	if len(symbol) == 0 {
		return fmt.Errorf("symbol can not be empty")
	}

	if len(symbol) < MIN_SYMBOL_LEN || len(symbol) > MAX_SYMBOL_LEN {
		return fmt.Errorf("invalid symbol, the length must be between [%d,%d]: %s", MIN_SYMBOL_LEN, MAX_SYMBOL_LEN, symbol)
	}

	return nil
}

// ValidateAddress validates if the given address is valid btc address
func ValidateAddress(addr string, netParams *chaincfg.Params) error {
	if len(addr) == 0 {
		return fmt.Errorf("address can not be empty")
	}

	if _, err := btcutil.DecodeAddress(addr, netParams); err != nil {
		return fmt.Errorf("invalid address %s: %v", addr, err)
	}

	return nil
}

// ValidateMetadata validates if the given metadata is valid
func ValidateMetadata(metadata string) error {
	if gjson.Valid(metadata) {
		return fmt.Errorf("metadata is not valid JSON")
	}

	return nil
}

// ValidatePubKey validates if the given public key is valid schnorr public key
func ValidatePubKey(pubKey string) error {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return fmt.Errorf("invalid public key: %v", err)
	}

	_, err = schnorr.ParsePubKey(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("invalid schnorr public key: %v", err)
	}

	return nil
}

// ValidateSignature validates if the given signature is valid schnorr signature
func ValidateSignature(signature string) error {
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	_, err = schnorr.ParseSignature(sig)
	if err != nil {
		return fmt.Errorf("invalid schnorr signature: %v", err)
	}

	return nil
}
