package cmd

import (
	"os"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"btc-sbt/stacks/basics"
)

// GetPrivateKeyAndAddress gets the private key from the given file path and generates the corresponding address by the given address type
func GetPrivateKeyAndAddress(keyPath string, addrType basics.AddressType, netParam *chaincfg.Params) (*secp256k1.PrivateKey, btcutil.Address, error) {
	keyWIFBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, nil, err
	}

	keyWIF, err := btcutil.DecodeWIF(string(keyWIFBytes))
	if err != nil {
		return nil, nil, err
	}

	addr, err := basics.GetAddress(keyWIF.PrivKey, addrType, netParam)
	if err != nil {
		return nil, nil, err
	}

	return keyWIF.PrivKey, addr, nil
}
