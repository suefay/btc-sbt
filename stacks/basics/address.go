package basics

import (
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// IsSegWitAddress checks if the given address is the segwit type
func IsSegWitAddress(address btcutil.Address) bool {
	switch address.(type) {
	case *btcutil.AddressWitnessPubKeyHash:
		return true

	case *btcutil.AddressWitnessScriptHash:
		return true

	case *btcutil.AddressTaproot:
		return true

	default:
		return false
	}
}

// IsTaprootAddress checks if the given address is the taproot type
func IsTaprootAddress(address btcutil.Address) bool {
	switch address.(type) {
	case *btcutil.AddressTaproot:
		return true

	default:
		return false
	}
}

// IsP2SHAddress checks if the given address is the P2SH type
func IsP2SHAddress(address btcutil.Address) bool {
	switch address.(type) {
	case *btcutil.AddressScriptHash:
		return true

	default:
		return false
	}
}

// GetTaprootAddress gets the taproot address from the given private key
func GetTaprootAddress(key *secp256k1.PrivateKey, netParams *chaincfg.Params) (*btcutil.AddressTaproot, error) {
	return btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(key.PubKey())), netParams)
}

// GetTaprootPkScript gets the taproot pk script from the given private key
func GetTaprootPkScript(key *secp256k1.PrivateKey, netParams *chaincfg.Params) ([]byte, error) {
	taprootAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootKeyNoScript(key.PubKey())), netParams)
	if err != nil {
		return nil, err
	}

	return txscript.PayToAddrScript(taprootAddress)
}
