package basics

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// AddressType represents the basic address type
type AddressType uint8

const (
	Taproot AddressType = iota
	WitnessPubKeyHash
	WitnessScriptHash
	PubKeyHash
	ScriptHash
)

// String implement fmt.Stringer
func (a AddressType) String() string {
	switch a {
	case Taproot:
		return "p2tr"

	case WitnessPubKeyHash:
		return "p2wpkh"

	case WitnessScriptHash:
		return "p2wsh"

	case PubKeyHash:
		return "p2pkh"

	case ScriptHash:
		return "p2sh"

	default:
		return ""
	}
}

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

// GetAddress gets the address from the given private key according to the specified address type
func GetAddress(key *secp256k1.PrivateKey, addrType AddressType, netParams *chaincfg.Params) (btcutil.Address, error) {
	switch addrType {
	case Taproot:
		return GetTaprootAddress(key, netParams)

	case WitnessPubKeyHash:
		return GetWitnessPubKeyHashAddress(key, netParams)

	case PubKeyHash:
		return GetPubKeyHashAddress(key, netParams)
	}

	return nil, fmt.Errorf("unsupported address type: %s", addrType)
}

// GetPubKeyHashAddress gets the public key hash address from the given private key
func GetPubKeyHashAddress(key *secp256k1.PrivateKey, netParams *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
	return btcutil.NewAddressPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), netParams)
}

// GetWitnessPubKeyHashAddress gets the witness public key hash address from the given private key
func GetWitnessPubKeyHashAddress(key *secp256k1.PrivateKey, netParams *chaincfg.Params) (*btcutil.AddressWitnessPubKeyHash, error) {
	return btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), netParams)
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
